// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/9 11:50
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : auth_routes.go
// @Software: GoLand
package router

import (
	"eLabX/src/api/system"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", system.UserLogin)
		authGroup.POST("/logout", system.UserLogout)
	}
}
