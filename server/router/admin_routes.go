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
	projectGroup := adminGroup.Group("/project")
	{
		projectGroup.POST("/list", admin.GetProjectList)
		projectGroup.POST("/detail", admin.GetProjectDetail)
		projectGroup.POST("/create", admin.CreateProject)
		projectGroup.POST("/update", admin.UpdateProject)
		projectGroup.POST("/delete", admin.DeleteProject)
	}
}
