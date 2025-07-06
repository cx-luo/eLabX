// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2023/12/12 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : router.go
// @Software: GoLand
package router

import (
	_ "eLabX/docs"
	"eLabX/src/api"
	"eLabX/src/api/system"
	"eLabX/src/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter returns a new router.
func NewRouter(outputPath string, loglevel string) *gin.Engine {
	// 设置全局 Logger
	router := gin.New()

	// 为需要中间件的路由组注册中间件

	// 使用 Zap 中间件
	router.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 注册其他中间件
	router.Use(middleware.CORS())

	router.Use(middleware.JwtAuth())

	router.Use(middleware.CasbinMiddleware())

	registerAuthRoutes(router)
	// 注册用户相关路由
	registerUserRoutes(router)

	registerCasbinRoutes(router)

	// 注册其他路由
	registerSystemRoutes(router)

	registerOtherRoutes(router)

	return router
}

func registerAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", api.UserLogin)
		authGroup.POST("/logout", api.UserLogout)
		authGroup.POST("/setUserAuthorities", api.SetUserAuthorities)
	}
}

// 用户相关路由
func registerUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/info", api.UserInfo)
		userGroup.POST("/name", api.FetchUserName)
		userGroup.GET("/list", api.GetUserList)
		userGroup.POST("/modify/pwd", api.ChangePwd)
		userGroup.POST("/modify/name", api.ModifyUserInfo)
		userGroup.POST("/forget/pwd", api.ForgetPwd)
	}
}

func registerCasbinRoutes(r *gin.Engine) {
	casbinGroup := r.Group("/api/casbin")
	{
		casbinGroup.POST("/add", api.CasbinAddPolicy)
	}
}

func registerSystemRoutes(r *gin.Engine) {
	systemGroup := r.Group("/api/system")
	menuGroup := systemGroup.Group("/menu")
	{
		menuGroup.POST("/tree", api.GetRouteTree)
		menuGroup.POST("/update", api.UpdateMenu)
		menuGroup.GET("/list", api.GetUserRouteList)
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
		userManagerGroup.POST("/reset/passwd", api.ResetPwd)
	}
}

// 其他路由
func registerOtherRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api")
	{
		otherGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
