// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/9 11:51
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : user_routes.go
// @Software: GoLand
package router

import (
	"eLabX/src/api/user"
	"github.com/gin-gonic/gin"
)

// 用户相关路由
func registerUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/info", user.UserInfo)
		userGroup.POST("/name", user.FetchUserName)
		userGroup.GET("/list", user.GetUserList)
		userGroup.POST("/modify/pwd", user.ChangePwd)
		userGroup.POST("/modify/name", user.ModifyUserInfo)
		userGroup.POST("/forget/pwd", user.ForgetPwd)
	}
}
