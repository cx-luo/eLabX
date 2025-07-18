// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/9 11:52
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : system_routes.go
// @Software: GoLand
package router

import (
	"eLabX/src/api/system"
	"eLabX/src/api/user"
	"github.com/gin-gonic/gin"
)

func registerSystemRoutes(r *gin.Engine) {
	systemGroup := r.Group("/api/system")
	menuGroup := systemGroup.Group("/menu")
	{
		menuGroup.POST("/tree", system.GetRouteTree)
		menuGroup.POST("/update", system.UpdateMenu)
		menuGroup.GET("/list", system.GetUserRouteList)
		menuGroup.POST("/add", system.AddMenu)
	}

	apiGroup := systemGroup.Group("/api")
	{
		apiGroup.GET("/list", system.GetApiList)
		apiGroup.POST("/add", system.AddApi)
		apiGroup.POST("/delete", system.DeleteAPi)
		apiGroup.POST("/update", system.UpdateAPi)
		apiGroup.GET("/refresh", system.RefreshApis)
	}

	userManagerGroup := systemGroup.Group("/user")
	{
		userManagerGroup.POST("/reset-passwd", user.ResetPwd)
		userManagerGroup.POST("/list", system.GetSystemUserList)
	}

	roleGroup := systemGroup.Group("/role")
	{
		roleGroup.POST("/list", system.GetRoleList)
		roleGroup.POST("/assign", system.RoleAssign)
		roleGroup.POST("/add", system.RoleAdd)
		roleGroup.POST("/delete", system.DeleteRole)
		roleGroup.GET("/info/:id", system.RoleInfo)
		roleGroup.POST("/update", system.UpdateRole)
	}
}
