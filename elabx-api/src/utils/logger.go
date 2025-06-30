// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/9 9:28
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : logger.go
// @Software: GoLand
package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var Apis []gin.RouteInfo

func GetAllRoutes(engine *gin.Engine) {
	for _, r := range engine.Routes() {
		Apis = append(Apis, r)
	}
	return
}
