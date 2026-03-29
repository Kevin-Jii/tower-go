package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Kevin-Jii/tower-go/pkg/xpyun/model"
)

// httpClient HTTP客户端（带超时设置）
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// HttpPostJson 发送HTTP POST JSON请求
func HttpPostJson(url string, data interface{}) *model.XPYunResp {
	b, err := json.Marshal(&data)
	if err != nil {
		result := model.XPYunResp{
			HttpStatusCode: 500,
		}
		return &result
	}

	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("http post error: %v\n", err)
		result := model.XPYunResp{
			HttpStatusCode: 500,
		}
		return &result
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read body error: %v\n", err)
		result := model.XPYunResp{
			HttpStatusCode: resp.StatusCode,
		}
		return &result
	}

	result := model.XPYunResp{
		HttpStatusCode: resp.StatusCode,
	}

	var content model.XPYunRespContent
	err = json.Unmarshal(body, &content)
	if err == nil {
		result.Content = &content
	}

	return &result
}