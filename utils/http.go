package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Post 发送POST请求
func (c *HTTPClient) Post(url string, body interface{}, result interface{}) error {
	// 序列化请求体
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 如果需要解析响应
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// Get 发送GET请求
func (c *HTTPClient) Get(url string, result interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
