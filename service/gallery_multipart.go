package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	minS3MultipartChunkSize int64 = 5 * 1024 * 1024
	maxS3MultipartParts           = 10000
	multipartRecoveryDelay        = 2 * time.Minute
)

var galleryImageExtensions = map[string]string{
	".gif":  "image/gif",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

func (s *GalleryService) InitMultipartUpload(
	ctx context.Context,
	userID, storeID uint,
	req *model.InitGalleryMultipartUploadReq,
) (*model.GalleryMultipartUploadStatus, error) {
	if s.rustfsService == nil {
		return nil, apicode.New(apicode.FileServiceUnavailable)
	}
	fileName, contentType, err := s.validateMultipartRequest(req)
	if err != nil {
		return nil, err
	}
	category, err := normalizeGalleryCategory(req.Category)
	if err != nil {
		return nil, err
	}

	existing, findErr := s.galleryModule.FindResumableUploadSession(userID, storeID, strings.ToLower(req.Fingerprint))
	if findErr == nil {
		if existing.FileName != fileName || existing.FileSize != req.FileSize || existing.Category != category {
			return nil, apicode.New(apicode.UploadSessionConflict)
		}
		if existing.Status == model.GalleryUploadStatusCompleted && time.Now().Before(existing.ExpiresAt) {
			gallery, err := s.galleryModule.GetByID(existing.GalleryID)
			if err != nil {
				return nil, apicode.Wrap(apicode.InternalError, err)
			}
			s.refreshGalleryURL(gallery)
			return completedMultipartStatus(existing, gallery), nil
		}
		if existing.Status == model.GalleryUploadStatusCompleted {
			// 幂等窗口结束后允许再次上传相同文件。
		} else if time.Now().After(existing.ExpiresAt) {
			_ = s.expireMultipartSession(ctx, existing)
		} else if existing.Status == model.GalleryUploadStatusCompleting {
			gallery, reset, err := s.recoverCompletingSession(ctx, existing)
			if err != nil {
				if apicode.Is(err, apicode.UploadAlreadyCompleting) {
					return completingMultipartStatus(existing), nil
				}
				return nil, err
			}
			if gallery != nil {
				return completedMultipartStatus(existing, gallery), nil
			}
			if reset {
				return s.multipartStatus(ctx, existing)
			}
		} else {
			return s.multipartStatus(ctx, existing)
		}
	} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, apicode.Wrap(apicode.InternalError, findErr)
	}

	chunkSize := calculateMultipartChunkSize(req.FileSize, s.multipartChunkSize)
	totalParts := int((req.FileSize + chunkSize - 1) / chunkSize)
	objectPath := newGalleryObjectPath(category, fileName)
	uploadID, err := s.rustfsService.NewMultipartUpload(ctx, objectPath, contentType)
	if err != nil {
		return nil, apicode.Wrap(apicode.MultipartInitFailed, err)
	}

	session := &model.GalleryUploadSession{
		ID:          uuid.NewString(),
		UploadID:    uploadID,
		ObjectPath:  objectPath,
		Fingerprint: strings.ToLower(req.Fingerprint),
		FileName:    fileName,
		FileSize:    req.FileSize,
		ContentType: contentType,
		Category:    category,
		Remark:      strings.TrimSpace(req.Remark),
		ChunkSize:   chunkSize,
		TotalParts:  totalParts,
		StoreID:     storeID,
		UserID:      userID,
		Status:      model.GalleryUploadStatusUploading,
		ExpiresAt:   time.Now().Add(s.uploadSessionTTL),
	}
	if err := s.galleryModule.CreateUploadSession(session); err != nil {
		_ = s.rustfsService.AbortMultipartUpload(ctx, objectPath, uploadID)
		return nil, apicode.Wrap(apicode.InternalError, err)
	}
	return s.multipartStatus(ctx, session)
}

func (s *GalleryService) GetMultipartUploadStatus(
	ctx context.Context,
	sessionID string,
	userID, storeID uint,
) (*model.GalleryMultipartUploadStatus, error) {
	session, err := s.getMultipartSession(sessionID, userID, storeID)
	if err != nil {
		return nil, err
	}
	if session.Status == model.GalleryUploadStatusCompleted {
		gallery, err := s.galleryModule.GetByID(session.GalleryID)
		if err != nil {
			return nil, apicode.Wrap(apicode.InternalError, err)
		}
		s.refreshGalleryURL(gallery)
		return completedMultipartStatus(session, gallery), nil
	}
	if session.Status == model.GalleryUploadStatusExpired || time.Now().After(session.ExpiresAt) {
		_ = s.expireMultipartSession(ctx, session)
		return nil, apicode.New(apicode.UploadSessionExpired)
	}
	if session.Status == model.GalleryUploadStatusCompleting {
		gallery, reset, err := s.recoverCompletingSession(ctx, session)
		if err != nil {
			if apicode.Is(err, apicode.UploadAlreadyCompleting) {
				return completingMultipartStatus(session), nil
			}
			return nil, err
		}
		if gallery != nil {
			return completedMultipartStatus(session, gallery), nil
		}
		if reset {
			return s.multipartStatus(ctx, session)
		}
	}
	if session.Status != model.GalleryUploadStatusUploading {
		return nil, apicode.New(apicode.UploadSessionNotFound)
	}
	return s.multipartStatus(ctx, session)
}

func (s *GalleryService) UploadMultipartPart(
	ctx context.Context,
	sessionID string,
	userID, storeID uint,
	partNumber int,
	contentLength int64,
	reader io.Reader,
) (*model.GalleryUploadPartResult, error) {
	session, err := s.getActiveMultipartSession(ctx, sessionID, userID, storeID)
	if err != nil {
		return nil, err
	}
	expectedSize, ok := expectedMultipartPartSize(session, partNumber)
	if !ok || contentLength != expectedSize {
		return nil, apicode.Newf(
			apicode.InvalidUploadPart,
			"第%d个分片大小应为%d字节",
			partNumber,
			expectedSize,
		)
	}

	part, err := s.rustfsService.UploadPart(
		ctx,
		session.ObjectPath,
		session.UploadID,
		partNumber,
		reader,
		expectedSize,
	)
	if err != nil {
		return nil, apicode.Wrap(apicode.MultipartPartFailed, err)
	}
	return &model.GalleryUploadPartResult{PartNumber: part.PartNumber, Size: part.Size, ETag: part.ETag}, nil
}

func (s *GalleryService) CompleteMultipartUpload(
	ctx context.Context,
	sessionID string,
	userID, storeID uint,
) (*model.Gallery, error) {
	session, err := s.getMultipartSession(sessionID, userID, storeID)
	if err != nil {
		return nil, err
	}
	if session.Status == model.GalleryUploadStatusCompleted {
		gallery, err := s.galleryModule.GetByID(session.GalleryID)
		if err != nil {
			return nil, apicode.Wrap(apicode.InternalError, err)
		}
		s.refreshGalleryURL(gallery)
		return gallery, nil
	}
	if session.Status == model.GalleryUploadStatusCompleting {
		gallery, reset, err := s.recoverCompletingSession(ctx, session)
		if err != nil {
			return nil, err
		}
		if gallery != nil {
			return gallery, nil
		}
		if !reset {
			return nil, apicode.New(apicode.UploadAlreadyCompleting)
		}
	}
	if _, err := s.getActiveMultipartSession(ctx, sessionID, userID, storeID); err != nil {
		return nil, err
	}

	parts, err := s.rustfsService.ListUploadedParts(ctx, session.ObjectPath, session.UploadID)
	if err != nil {
		return nil, apicode.Wrap(apicode.MultipartCompleteFailed, err)
	}
	parts, err = validateMultipartParts(session, parts)
	if err != nil {
		return nil, err
	}
	marked, err := s.galleryModule.MarkUploadSessionCompleting(session.ID)
	if err != nil {
		return nil, apicode.Wrap(apicode.InternalError, err)
	}
	if !marked {
		return nil, apicode.New(apicode.UploadAlreadyCompleting)
	}
	session.Status = model.GalleryUploadStatusCompleting

	result, err := s.rustfsService.CompleteMultipartUpload(
		ctx,
		session.ObjectPath,
		session.UploadID,
		session.ContentType,
		parts,
		session.FileSize,
	)
	if err != nil {
		return nil, apicode.Wrap(apicode.MultipartCompleteFailed, err)
	}

	gallery := galleryFromMultipartResult(session, result)
	if err := s.galleryModule.CompleteUploadSession(session, gallery); err != nil {
		_ = s.rustfsService.Delete(result.Path)
		_ = s.galleryModule.UpdateUploadSessionStatus(session.ID, model.GalleryUploadStatusFailed)
		return nil, apicode.Wrap(apicode.GalleryRecordSaveFailed, err)
	}
	return gallery, nil
}

func (s *GalleryService) AbortMultipartUpload(
	ctx context.Context,
	sessionID string,
	userID, storeID uint,
) error {
	session, err := s.getMultipartSession(sessionID, userID, storeID)
	if err != nil {
		return err
	}
	if session.Status == model.GalleryUploadStatusAborted || session.Status == model.GalleryUploadStatusExpired {
		return nil
	}
	if session.Status == model.GalleryUploadStatusCompleted {
		return apicode.New(apicode.UploadSessionConflict)
	}
	if err := s.rustfsService.AbortMultipartUpload(ctx, session.ObjectPath, session.UploadID); err != nil {
		return apicode.Wrap(apicode.MultipartAbortFailed, err)
	}
	if err := s.galleryModule.UpdateUploadSessionStatus(session.ID, model.GalleryUploadStatusAborted); err != nil {
		return apicode.Wrap(apicode.InternalError, err)
	}
	return nil
}

func (s *GalleryService) CleanupExpiredMultipartUploads(ctx context.Context) error {
	sessions, err := s.galleryModule.ListExpiredUploadSessions(time.Now(), 500)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if session.Status == model.GalleryUploadStatusCompleting {
			_ = s.rustfsService.Delete(session.ObjectPath)
		}
		if err := s.expireMultipartSession(ctx, session); err != nil {
			logging.LogWarn("清理过期图库上传会话失败", zap.String("sessionID", session.ID), zap.Error(err))
		}
	}
	return nil
}

func (s *GalleryService) validateMultipartRequest(req *model.InitGalleryMultipartUploadReq) (string, string, error) {
	fileName := sanitizeUploadFileName(req.FileName)
	ext := strings.ToLower(path.Ext(fileName))
	defaultContentType, ok := galleryImageExtensions[ext]
	if !ok {
		return "", "", apicode.New(apicode.ImageFormatUnsupported)
	}
	if req.FileSize > s.multipartMaxSize {
		return "", "", apicode.Newf(
			apicode.MultipartFileTooLarge,
			"文件大小不能超过%dMB",
			s.multipartMaxSize/(1024*1024),
		)
	}
	return fileName, defaultContentType, nil
}

func (s *GalleryService) getMultipartSession(sessionID string, userID, storeID uint) (*model.GalleryUploadSession, error) {
	session, err := s.galleryModule.GetUploadSession(sessionID, userID, storeID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apicode.New(apicode.UploadSessionNotFound)
	}
	if err != nil {
		return nil, apicode.Wrap(apicode.InternalError, err)
	}
	return session, nil
}

func (s *GalleryService) getActiveMultipartSession(
	ctx context.Context,
	sessionID string,
	userID, storeID uint,
) (*model.GalleryUploadSession, error) {
	session, err := s.getMultipartSession(sessionID, userID, storeID)
	if err != nil {
		return nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		_ = s.expireMultipartSession(ctx, session)
		return nil, apicode.New(apicode.UploadSessionExpired)
	}
	if session.Status != model.GalleryUploadStatusUploading {
		if session.Status == model.GalleryUploadStatusCompleting {
			return nil, apicode.New(apicode.UploadAlreadyCompleting)
		}
		return nil, apicode.New(apicode.UploadSessionNotFound)
	}
	return session, nil
}

func (s *GalleryService) multipartStatus(
	ctx context.Context,
	session *model.GalleryUploadSession,
) (*model.GalleryMultipartUploadStatus, error) {
	parts, err := s.rustfsService.ListUploadedParts(ctx, session.ObjectPath, session.UploadID)
	if err != nil {
		return nil, apicode.Wrap(apicode.ExternalServiceFailed, err)
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].PartNumber < parts[j].PartNumber })
	uploaded := make([]model.GalleryUploadedPart, 0, len(parts))
	var uploadedBytes int64
	for _, part := range parts {
		uploaded = append(uploaded, model.GalleryUploadedPart{PartNumber: part.PartNumber, Size: part.Size})
		uploadedBytes += part.Size
	}
	return &model.GalleryMultipartUploadStatus{
		SessionID:     session.ID,
		ChunkSize:     session.ChunkSize,
		TotalParts:    session.TotalParts,
		UploadedParts: uploaded,
		UploadedBytes: uploadedBytes,
		FileSize:      session.FileSize,
		Status:        session.Status,
		ExpiresAt:     session.ExpiresAt,
	}, nil
}

func (s *GalleryService) expireMultipartSession(ctx context.Context, session *model.GalleryUploadSession) error {
	if err := s.rustfsService.AbortMultipartUpload(ctx, session.ObjectPath, session.UploadID); err != nil {
		return err
	}
	return s.galleryModule.UpdateUploadSessionStatus(session.ID, model.GalleryUploadStatusExpired)
}

func (s *GalleryService) recoverCompletingSession(
	ctx context.Context,
	session *model.GalleryUploadSession,
) (*model.Gallery, bool, error) {
	// 正常合并期间不抢占状态；仅在进程异常退出后留下的陈旧 completing 会话上执行恢复。
	if time.Since(session.UpdatedAt) < multipartRecoveryDelay {
		return nil, false, apicode.New(apicode.UploadAlreadyCompleting)
	}
	if _, err := s.rustfsService.ListUploadedParts(ctx, session.ObjectPath, session.UploadID); err == nil {
		if err := s.galleryModule.UpdateUploadSessionStatus(session.ID, model.GalleryUploadStatusUploading); err != nil {
			return nil, false, apicode.Wrap(apicode.InternalError, err)
		}
		session.Status = model.GalleryUploadStatusUploading
		return nil, true, nil
	} else if !errors.Is(err, ErrRustFSMultipartUploadNotFound) {
		return nil, false, apicode.Wrap(apicode.ExternalServiceFailed, err)
	}

	result, err := s.rustfsService.StatObject(ctx, session.ObjectPath)
	if errors.Is(err, ErrRustFSObjectNotFound) {
		_ = s.galleryModule.UpdateUploadSessionStatus(session.ID, model.GalleryUploadStatusFailed)
		return nil, false, apicode.New(apicode.UploadSessionNotFound)
	}
	if err != nil {
		return nil, false, apicode.Wrap(apicode.ExternalServiceFailed, err)
	}
	if result.Size != session.FileSize {
		_ = s.rustfsService.Delete(session.ObjectPath)
		_ = s.galleryModule.UpdateUploadSessionStatus(session.ID, model.GalleryUploadStatusFailed)
		return nil, false, apicode.New(apicode.InvalidUploadPart)
	}
	gallery := galleryFromMultipartResult(session, result)
	if err := s.galleryModule.CompleteUploadSession(session, gallery); err != nil {
		return nil, false, apicode.Wrap(apicode.GalleryRecordSaveFailed, err)
	}
	session.Status = model.GalleryUploadStatusCompleted
	session.GalleryID = gallery.ID
	return gallery, false, nil
}

func galleryFromMultipartResult(session *model.GalleryUploadSession, result *RustFSUploadResult) *model.Gallery {
	return &model.Gallery{
		Name:     session.FileName,
		Path:     result.Path,
		URL:      result.URL,
		Size:     result.Size,
		MimeType: session.ContentType,
		Category: session.Category,
		StoreID:  session.StoreID,
		UploadBy: session.UserID,
		Remark:   session.Remark,
	}
}

func completedMultipartStatus(session *model.GalleryUploadSession, gallery *model.Gallery) *model.GalleryMultipartUploadStatus {
	return &model.GalleryMultipartUploadStatus{
		SessionID:     session.ID,
		ChunkSize:     session.ChunkSize,
		TotalParts:    session.TotalParts,
		UploadedParts: []model.GalleryUploadedPart{},
		UploadedBytes: session.FileSize,
		FileSize:      session.FileSize,
		Status:        session.Status,
		ExpiresAt:     session.ExpiresAt,
		Gallery:       gallery,
	}
}

func completingMultipartStatus(session *model.GalleryUploadSession) *model.GalleryMultipartUploadStatus {
	return &model.GalleryMultipartUploadStatus{
		SessionID:     session.ID,
		ChunkSize:     session.ChunkSize,
		TotalParts:    session.TotalParts,
		UploadedParts: []model.GalleryUploadedPart{},
		UploadedBytes: session.FileSize,
		FileSize:      session.FileSize,
		Status:        model.GalleryUploadStatusCompleting,
		ExpiresAt:     session.ExpiresAt,
	}
}

func calculateMultipartChunkSize(fileSize, configured int64) int64 {
	if configured < minS3MultipartChunkSize {
		configured = minS3MultipartChunkSize
	}
	minimumForPartLimit := (fileSize + maxS3MultipartParts - 1) / maxS3MultipartParts
	if minimumForPartLimit > configured {
		const mib = int64(1024 * 1024)
		configured = ((minimumForPartLimit + mib - 1) / mib) * mib
	}
	return configured
}

func expectedMultipartPartSize(session *model.GalleryUploadSession, partNumber int) (int64, bool) {
	if session == nil || partNumber < 1 || partNumber > session.TotalParts {
		return 0, false
	}
	if partNumber < session.TotalParts {
		return session.ChunkSize, true
	}
	return session.FileSize - int64(session.TotalParts-1)*session.ChunkSize, true
}

func validateMultipartParts(session *model.GalleryUploadSession, parts []RustFSUploadedPart) ([]RustFSUploadedPart, error) {
	if len(parts) != session.TotalParts {
		return nil, apicode.Newf(apicode.UploadIncomplete, "已上传%d/%d个分片", len(parts), session.TotalParts)
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].PartNumber < parts[j].PartNumber })
	for index, part := range parts {
		expectedPartNumber := index + 1
		expectedSize, _ := expectedMultipartPartSize(session, expectedPartNumber)
		if part.PartNumber != expectedPartNumber || part.Size != expectedSize || strings.TrimSpace(part.ETag) == "" {
			return nil, apicode.Newf(apicode.InvalidUploadPart, "第%d个分片校验失败", expectedPartNumber)
		}
	}
	return parts, nil
}

func sanitizeUploadFileName(fileName string) string {
	fileName = strings.ReplaceAll(strings.TrimSpace(fileName), "\\", "/")
	return path.Base(fileName)
}

func normalizeGalleryCategory(category string) (string, error) {
	category = strings.TrimSpace(category)
	if category == "" {
		return "other", nil
	}
	switch category {
	case "product", "supplier", "avatar", "purchase", "other":
		return category, nil
	default:
		return "", apicode.New(apicode.InvalidParameter.WithMessage("图库分类无效"))
	}
}

func newGalleryObjectPath(category, fileName string) string {
	now := time.Now()
	ext := strings.ToLower(path.Ext(fileName))
	uniqueName := fmt.Sprintf("%s_%s%s", now.Format("150405"), uuid.NewString()[:8], ext)
	return path.Join("gallery", category, now.Format("2006/01/02"), uniqueName)
}
