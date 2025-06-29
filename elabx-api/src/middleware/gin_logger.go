// Package middleware coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/29 11:27
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : gin_logger.go
// @Software: GoLand
package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// 使用自定义 ResponseWriter
		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw
		// 打印请求信息
		reqBody, _ := c.GetRawData()
		// 请求包体写回。
		if len(reqBody) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.String("method", c.Request.Method), // 请求方法类型 eg: GET
			zap.String("path", path),               // 请求路径 eg: /test
			zap.Int("status", c.Writer.Status()),   // 状态码 eg: 200
			zap.Duration("duration", cost),         // 返回花费时间
			zap.String("query", string(reqBody)),   // 请求参数 eg: name=1&password=2
			zap.String("ip", c.ClientIP()),         // 返回真实的客户端IP eg: ::1（这个就是本机IP，ipv6地址）
			//zap.String("user-agent", c.Request.UserAgent()),                      // 返回客户端的用户代理。 eg: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 返回Errors 切片中ErrorTypePrivate类型的错误

		)
	}
}
