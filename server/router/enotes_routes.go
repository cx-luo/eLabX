// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/9 11:51
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : enotes_routes.go
// @Software: GoLand
package router

import (
	"eLabX/src/api/enote"

	"github.com/gin-gonic/gin"
)

func registerEnoteRoutes(r *gin.Engine) {
	enoteGroup := r.Group("/api/enote")
	{
		enoteGroup.POST("/saveRxnToServer", enote.SaveRxnToServer)
	}
}
