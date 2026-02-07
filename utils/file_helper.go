package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SaveImageFile 保存图片文件到指定目录，返回访问URL
// 用于 dingtalk 服务保存图片到本地
// filename: 文件名（不含路径）
// imageData: 图片数据
// 返回: 图片访问URL, 错误
func SaveImageFile(filename string, imageData []byte) (string, error) {
	// 创建按日期分类的子目录 (例如: 2024/01/15)
	today := time.Now().Format("2006/01/02")
	uploadDir := filepath.Join("uploads", "images", today)

	// 创建目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save image file: %w", err)
	}

	// 返回相对URL路径
	relativeURL := filepath.Join("/uploads/images", today, filename)
	// 转换为正斜杠（URL格式）
	relativeURL = filepath.ToSlash(relativeURL)

	return relativeURL, nil
}

// DeleteImageFile 删除图片文件
func DeleteImageFile(imageURL string) error {
	// 从URL提取相对路径
	// imageURL 格式: /uploads/images/2024/01/15/filename.jpg
	filePath := filepath.Join(".", imageURL)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，不报错
		}
		return err
	}

	// 删除文件
	return os.Remove(filePath)
}

// CleanOldImages 清理旧图片（可选功能，可用于定期清理）
// days: 保留天数
func CleanOldImages(days int) error {
	uploadDir := filepath.Join("uploads", "images")

	cutoffTime := time.Now().AddDate(0, 0, -days)

	return filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理文件，跳过目录
		if !info.IsDir() && info.ModTime().Before(cutoffTime) {
			return os.Remove(path)
		}

		return nil
	})
}
