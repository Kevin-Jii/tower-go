package util

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

// Sign 生成SHA1签名
func Sign(signSource string) string {
	h := sha1.New()
	h.Write([]byte(signSource))
	result := fmt.Sprintf("%x", h.Sum(nil))
	return result
}

// GetMillisecond 获取当前毫秒时间戳
func GetMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// StrRepeat 重复字符串
func StrRepeat(str string, repeatTimes int) string {
	return strings.Repeat(str, repeatTimes)
}

// CalcGbkLenForPrint 计算GBK编码后的长度
func CalcGbkLenForPrint(data string) int {
	return len(data) // 简化版本，实际需要GBK编码
}

// CalcAsciiLenForPrint 计算ASCII长度
func CalcAsciiLenForPrint(data string) int {
	return len(data)
}