package service

import (
	"context"
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
)

func TestCalculateMultipartChunkSize(t *testing.T) {
	const mib = int64(1024 * 1024)
	tests := []struct {
		name       string
		fileSize   int64
		configured int64
		want       int64
	}{
		{name: "enforces S3 minimum", fileSize: mib, configured: mib, want: 5 * mib},
		{name: "keeps configured size", fileSize: 100 * mib, configured: 8 * mib, want: 8 * mib},
		{name: "stays within part limit", fileSize: 50001 * mib, configured: 5 * mib, want: 6 * mib},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMultipartChunkSize(tt.fileSize, tt.configured); got != tt.want {
				t.Fatalf("calculateMultipartChunkSize() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestExpectedMultipartPartSize(t *testing.T) {
	const mib = int64(1024 * 1024)
	session := &model.GalleryUploadSession{FileSize: 12 * mib, ChunkSize: 5 * mib, TotalParts: 3}
	for partNumber, want := range map[int]int64{1: 5 * mib, 2: 5 * mib, 3: 2 * mib} {
		got, ok := expectedMultipartPartSize(session, partNumber)
		if !ok || got != want {
			t.Fatalf("part %d size = %d, %v; want %d, true", partNumber, got, ok, want)
		}
	}
	if _, ok := expectedMultipartPartSize(session, 4); ok {
		t.Fatal("out-of-range part unexpectedly accepted")
	}
}

func TestValidateMultipartParts(t *testing.T) {
	session := &model.GalleryUploadSession{FileSize: 12, ChunkSize: 5, TotalParts: 3}
	valid := []RustFSUploadedPart{
		{PartNumber: 3, Size: 2, ETag: "three"},
		{PartNumber: 1, Size: 5, ETag: "one"},
		{PartNumber: 2, Size: 5, ETag: "two"},
	}
	parts, err := validateMultipartParts(session, valid)
	if err != nil {
		t.Fatalf("validateMultipartParts() error = %v", err)
	}
	for index, part := range parts {
		if part.PartNumber != index+1 {
			t.Fatalf("part index %d has number %d", index, part.PartNumber)
		}
	}

	if _, err := validateMultipartParts(session, valid[:2]); !apicode.Is(err, apicode.UploadIncomplete) {
		t.Fatalf("missing part error = %v", err)
	}
	invalid := append([]RustFSUploadedPart(nil), valid...)
	invalid[0].Size = 1
	if _, err := validateMultipartParts(session, invalid); !apicode.Is(err, apicode.InvalidUploadPart) {
		t.Fatalf("invalid part error = %v", err)
	}
}

func TestSanitizeUploadFileName(t *testing.T) {
	if got := sanitizeUploadFileName(`C:\\fakepath\\image.png`); got != "image.png" {
		t.Fatalf("sanitizeUploadFileName() = %q", got)
	}
}

func TestNormalizeGalleryCategory(t *testing.T) {
	if got, err := normalizeGalleryCategory(""); err != nil || got != "other" {
		t.Fatalf("empty category = %q, %v", got, err)
	}
	if got, err := normalizeGalleryCategory("product"); err != nil || got != "product" {
		t.Fatalf("product category = %q, %v", got, err)
	}
	if _, err := normalizeGalleryCategory("../../escape"); !apicode.Is(err, apicode.InvalidParameter) {
		t.Fatalf("invalid category error = %v", err)
	}
}

func TestValidateMultipartRequest(t *testing.T) {
	service := &GalleryService{multipartMaxSize: 10 * 1024 * 1024}
	fileName, contentType, err := service.validateMultipartRequest(&model.InitGalleryMultipartUploadReq{
		FileName:    "photo.PNG",
		FileSize:    1024,
		ContentType: "application/octet-stream",
	})
	if err != nil || fileName != "photo.PNG" || contentType != "image/png" {
		t.Fatalf("valid request = (%q, %q, %v)", fileName, contentType, err)
	}
	if _, _, err := service.validateMultipartRequest(&model.InitGalleryMultipartUploadReq{
		FileName: "photo.svg",
		FileSize: 1024,
	}); !apicode.Is(err, apicode.ImageFormatUnsupported) {
		t.Fatalf("unsupported format error = %v", err)
	}
	if _, _, err := service.validateMultipartRequest(&model.InitGalleryMultipartUploadReq{
		FileName: "photo.png",
		FileSize: 11 * 1024 * 1024,
	}); !apicode.Is(err, apicode.MultipartFileTooLarge) {
		t.Fatalf("oversized file error = %v", err)
	}
}

func TestFreshCompletingSessionIsNotRecovered(t *testing.T) {
	service := &GalleryService{}
	_, reset, err := service.recoverCompletingSession(context.Background(), &model.GalleryUploadSession{
		UpdatedAt: time.Now(),
	})
	if reset || !apicode.Is(err, apicode.UploadAlreadyCompleting) {
		t.Fatalf("recoverCompletingSession() = reset %v, error %v", reset, err)
	}
}
