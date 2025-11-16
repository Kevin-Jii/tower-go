package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
)

// SaveImageFile 保存图片文件到指定目录，返回访问URL
// filename: 文件名（不含路径）
// imageData: 图片数据
// 返回: 图片访问URL, 错误
func SaveImageFile(filename string, imageData []byte) (string, error) {
	cfg := config.GetConfig()
	uploadPath := cfg.App.ImageUploadPath
	baseURL := cfg.App.ImageBaseURL

	// 创建按日期分类的子目录 (例如: 2024/01/15)
	today := time.Now().Format("2006/01/02")
	targetDir := filepath.Join(uploadPath, today)

	// 确保目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 生成唯一文件名（时间戳 + MD5）
	timestamp := time.Now().Format("150405")
	hash := md5.Sum(imageData)
	hashStr := hex.EncodeToString(hash[:])[:8] // 取前8位
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png"
	}
	uniqueFilename := fmt.Sprintf("%s_%s%s", timestamp, hashStr, ext)

	// 完整文件路径
	filePath := filepath.Join(targetDir, uniqueFilename)

	// 写入文件
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// 生成访问URL（使用正斜杠，适配web路径）
	relPath := filepath.Join(today, uniqueFilename)
	// Windows路径转换为URL路径
	relPath = filepath.ToSlash(relPath)
	imageURL := fmt.Sprintf("%s/%s", baseURL, relPath)

	return imageURL, nil
}

// DeleteImageFile 删除图片文件
func DeleteImageFile(imageURL string) error {
	cfg := config.GetConfig()
	uploadPath := cfg.App.ImageUploadPath
	baseURL := cfg.App.ImageBaseURL

	// 从URL提取相对路径
	if len(imageURL) <= len(baseURL) {
		return fmt.Errorf("invalid image URL")
	}

	relPath := imageURL[len(baseURL)+1:] // 去掉baseURL和斜杠
	filePath := filepath.Join(uploadPath, relPath)

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// CleanOldImages 清理旧图片（可选功能，可用于定期清理）
// days: 保留天数
func CleanOldImages(days int) error {
	cfg := config.GetConfig()
	uploadPath := cfg.App.ImageUploadPath

	cutoffTime := time.Now().AddDate(0, 0, -days)

	return filepath.Walk(uploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件修改时间
		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				// 记录错误但继续处理其他文件
				fmt.Printf("Failed to delete old image %s: %v\n", path, err)
			}
		}

		return nil
	})
}
