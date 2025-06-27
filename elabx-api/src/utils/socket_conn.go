// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/11 11:03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : socket_conn.go
// @Software: GoLand
package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

// SendDataByApi 发送数据到指定 API 地址，返回原始响应或错误
func SendDataByApi(data interface{}, apiUrl string) ([]byte, error) {
	client := resty.New()

	// 设置客户端通用配置（可选）
	client.SetTimeout(200 * time.Second) // 设置超时时间
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("accept", "*/*")

	resp, err := client.R().SetBody(data).Post(apiUrl)

	if err != nil {
		return nil, fmt.Errorf("api request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
