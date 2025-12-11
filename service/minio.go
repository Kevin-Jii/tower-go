package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// RustFSService RustFS文件服务（S3兼容）
type RustFSService struct {
	client       *minio.Client
	bucketName   string
	notifyBucket string // 通知图片专用bucket
	endpoint     string
	useSSL       bool
}

// RustFSUploadResult 上传结果
type RustFSUploadResult struct {
	Path string `json:"path"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	ETag string `json:"etag"`
}

// NewRustFSService 创建RustFS服务实例
func NewRustFSService(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*RustFSService, error) {
	return NewRustFSServiceWithNotify(endpoint, accessKey, secretKey, bucketName, "", useSSL)
}

// NewRustFSServiceWithNotify 创建RustFS服务实例（支持通知bucket）
func NewRustFSServiceWithNotify(endpoint, accessKey, secretKey, bucketName, notifyBucket string, useSSL bool) (*RustFSService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logging.LogError("RustFS客户端创建失败", zap.Error(err))
		return nil, fmt.Errorf("创建RustFS客户端失败: %v", err)
	}

	service := &RustFSService{
		client:       client,
		bucketName:   bucketName,
		notifyBucket: notifyBucket,
		endpoint:     endpoint,
		useSSL:       useSSL,
	}

	// 确保主bucket存在
	if err := service.ensureBucket(); err != nil {
		return nil, err
	}

	// 确保通知bucket存在
	if notifyBucket != "" && notifyBucket != bucketName {
		if err := service.ensureBucketByName(notifyBucket); err != nil {
			logging.LogWarn("通知bucket创建失败，将使用主bucket", zap.Error(err))
			service.notifyBucket = bucketName
		}
	}

	logging.LogInfo("RustFS服务初始化成功", zap.String("endpoint", endpoint), zap.String("bucket", bucketName), zap.String("notifyBucket", service.notifyBucket))
	return service, nil
}

// ensureBucket 确保bucket存在
func (s *RustFSService) ensureBucket() error {
	return s.ensureBucketByName(s.bucketName)
}

// ensureBucketByName 确保指定bucket存在
func (s *RustFSService) ensureBucketByName(bucketName string) error {
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("检查bucket失败: %v", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建bucket失败: %v", err)
		}
		logging.LogInfo("RustFS bucket创建成功", zap.String("bucket", bucketName))
	}

	return nil
}

// Upload 上传文件
func (s *RustFSService) Upload(folder, filename string, reader io.Reader, fileSize int64, contentType string) (*RustFSUploadResult, error) {
	ctx := context.Background()

	// 构建对象路径
	objectName := filepath.Join(folder, filename)
	objectName = strings.ReplaceAll(objectName, "\\", "/") // 统一使用/

	// 上传文件
	info, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		logging.LogError("RustFS上传失败", zap.Error(err), zap.String("object", objectName))
		return nil, fmt.Errorf("上传失败: %v", err)
	}

	logging.LogInfo("RustFS文件上传成功", zap.String("object", objectName), zap.Int64("size", info.Size))

	return &RustFSUploadResult{
		Path: objectName,
		URL:  s.GetPublicURL(objectName),
		Name: filename,
		Size: info.Size,
		ETag: info.ETag,
	}, nil
}

// UploadToNotify 上传文件到通知bucket（不加密）
func (s *RustFSService) UploadToNotify(folder, filename string, reader io.Reader, fileSize int64, contentType string) (*RustFSUploadResult, error) {
	ctx := context.Background()

	// 使用通知bucket，如果没有配置则使用主bucket
	bucket := s.notifyBucket
	if bucket == "" {
		bucket = s.bucketName
	}

	// 构建对象路径
	objectName := filepath.Join(folder, filename)
	objectName = strings.ReplaceAll(objectName, "\\", "/")

	// 上传文件
	info, err := s.client.PutObject(ctx, bucket, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		logging.LogError("RustFS通知图片上传失败", zap.Error(err), zap.String("object", objectName), zap.String("bucket", bucket))
		return nil, fmt.Errorf("上传失败: %v", err)
	}

	logging.LogInfo("RustFS通知图片上传成功", zap.String("object", objectName), zap.String("bucket", bucket), zap.Int64("size", info.Size))

	return &RustFSUploadResult{
		Path: objectName,
		URL:  s.GetPublicURLForBucket(bucket, objectName),
		Name: filename,
		Size: info.Size,
		ETag: info.ETag,
	}, nil
}

// UploadMultipart 上传multipart文件
func (s *RustFSService) UploadMultipart(folder string, file multipart.File, header *multipart.FileHeader) (*RustFSUploadResult, error) {
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return s.Upload(folder, header.Filename, file, header.Size, contentType)
}

// GetPublicURL 获取公开访问URL
func (s *RustFSService) GetPublicURL(objectName string) string {
	return s.GetPublicURLForBucket(s.bucketName, objectName)
}

// GetPublicURLForBucket 获取指定bucket的公开访问URL
func (s *RustFSService) GetPublicURLForBucket(bucket, objectName string) string {
	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, bucket, objectName)
}

// GetPresignedURL 获取预签名URL（临时访问）
func (s *RustFSService) GetPresignedURL(objectName string, expires time.Duration) (string, error) {
	return s.GetPresignedURLForBucket(s.bucketName, objectName, expires)
}

// GetPresignedURLForBucket 获取指定bucket的预签名URL
func (s *RustFSService) GetPresignedURLForBucket(bucket, objectName string, expires time.Duration) (string, error) {
	ctx := context.Background()
	url, err := s.client.PresignedGetObject(ctx, bucket, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}
	return url.String(), nil
}

// GetNotifyBucket 获取通知bucket名称
func (s *RustFSService) GetNotifyBucket() string {
	if s.notifyBucket != "" {
		return s.notifyBucket
	}
	return s.bucketName
}

// Delete 删除文件
func (s *RustFSService) Delete(objectName string) error {
	ctx := context.Background()
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		logging.LogError("RustFS删除失败", zap.Error(err), zap.String("object", objectName))
		return fmt.Errorf("删除失败: %v", err)
	}
	logging.LogInfo("RustFS文件删除成功", zap.String("object", objectName))
	return nil
}

// List 列出文件
func (s *RustFSService) List(prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	var objects []minio.ObjectInfo

	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("列出文件失败: %v", object.Err)
		}
		objects = append(objects, object)
	}

	return objects, nil
}

// GetObjectInfo 获取文件信息
func (s *RustFSService) GetObjectInfo(objectName string) (*minio.ObjectInfo, error) {
	ctx := context.Background()
	info, err := s.client.StatObject(ctx, s.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}
	return &info, nil
}
