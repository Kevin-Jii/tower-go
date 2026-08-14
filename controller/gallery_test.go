package controller

import (
	"mime/multipart"
	"testing"

	"github.com/Kevin-Jii/tower-go/pkg/apicode"
)

func TestValidateGalleryImage(t *testing.T) {
	tests := []struct {
		name   string
		header *multipart.FileHeader
		code   apicode.Code
	}{
		{
			name:   "missing image",
			header: nil,
			code:   apicode.ImageRequired,
		},
		{
			name:   "unsupported extension",
			header: &multipart.FileHeader{Filename: "image.svg", Size: 1024},
			code:   apicode.ImageFormatUnsupported,
		},
		{
			name:   "image exceeds limit",
			header: &multipart.FileHeader{Filename: "image.png", Size: maxGalleryImageSize + 1},
			code:   apicode.ImageTooLarge,
		},
		{
			name:   "image at limit",
			header: &multipart.FileHeader{Filename: "IMAGE.PNG", Size: maxGalleryImageSize},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGalleryImage(tt.header)
			if tt.code.Num == 0 {
				if err != nil {
					t.Fatalf("validateGalleryImage() error = %v", err)
				}
				return
			}
			if !apicode.Is(err, tt.code) {
				t.Fatalf("validateGalleryImage() error = %v, want code %d", err, tt.code.Num)
			}
		})
	}
}
