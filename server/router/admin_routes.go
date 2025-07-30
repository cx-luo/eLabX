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

	groupGroup := adminGroup.Group("/group")
	{
		groupGroup.POST("/list", admin.GetGroupList)
		groupGroup.POST("/detail", admin.GetGroupDetail)
		groupGroup.POST("/create", admin.CreateGroup)
		groupGroup.POST("/update", admin.UpdateGroup)
		groupGroup.POST("/delete", admin.DeleteGroup)
		groupGroup.POST("/assign-permissions", admin.AssignPermissionsToGroup)
		groupGroup.POST("/permissions", admin.GetGroupPermissions)
		groupGroup.POST("/add-user", admin.AddUserToGroup)
		groupGroup.POST("/remove-user", admin.RemoveUserFromGroup)
		groupGroup.POST("/users", admin.GetGroupUsers)
	}

	permissionGroup := adminGroup.Group("/permission")
	{
		permissionGroup.POST("/list", admin.GetPermissionList)
		permissionGroup.POST("/create", admin.CreatePermission)
		permissionGroup.POST("/update", admin.UpdatePermission)
		permissionGroup.POST("/delete", admin.DeletePermission)
	}
}
