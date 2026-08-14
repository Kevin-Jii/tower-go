package middleware

import "strings"

const galleryMultipartPathPrefix = "/api/v1/galleries/multipart/"

func isGalleryMultipartPath(path string) bool {
	return path == "/api/v1/galleries/multipart/init" || strings.HasPrefix(path, galleryMultipartPathPrefix)
}

func isGalleryMultipartPartPath(path string) bool {
	if !strings.HasPrefix(path, galleryMultipartPathPrefix) {
		return false
	}
	rest := strings.TrimPrefix(path, galleryMultipartPathPrefix)
	return strings.Contains(rest, "/parts/")
}
