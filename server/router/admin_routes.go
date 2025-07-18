// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/17 17:04
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : admin_routes.go
// @Software: GoLand
package router

import (
	"eLabX/src/api/admin"
	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/api/admin")
	{
		adminGroup.POST("/project/list", admin.GetProjectList)
		adminGroup.POST("/project/detail", admin.GetProjectDetail)
		adminGroup.POST("/project/create", admin.CreateProject)
		adminGroup.POST("/project/update", admin.UpdateProject)
		adminGroup.POST("/project/delete", admin.DeleteProject)
	}
}
