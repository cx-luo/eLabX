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
	middleware2 "eLabX/middleware"
	"eLabX/src/api"
	"eLabX/src/api/casbin"
	"eLabX/src/api/system"
	"eLabX/src/api/user"
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
	router.Use(middleware2.GinLogger(), middleware2.GinRecovery(true))

	// 注册其他中间件
	router.Use(middleware2.CORS())

	router.Use(middleware2.JwtAuth())

	router.Use(middleware2.CasbinMiddleware())

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
		userGroup.GET("/info", user.UserInfo)
		userGroup.POST("/name", user.FetchUserName)
		userGroup.GET("/list", user.GetUserList)
		userGroup.POST("/modify/pwd", user.ChangePwd)
		userGroup.POST("/modify/name", user.ModifyUserInfo)
		userGroup.POST("/forget/pwd", user.ForgetPwd)
	}
}

func registerCasbinRoutes(r *gin.Engine) {
	casbinGroup := r.Group("/api/casbin")
	{
		casbinGroup.POST("/add", casbin.AddPolicy)
	}
}

func registerSystemRoutes(r *gin.Engine) {
	systemGroup := r.Group("/api/system")
	menuGroup := systemGroup.Group("/menu")
	{
		menuGroup.POST("/tree", system.GetRouteTree)
		menuGroup.POST("/update", system.UpdateMenu)
		menuGroup.GET("/list", system.GetUserRouteList)
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
		userManagerGroup.POST("/reset/passwd", user.ResetPwd)
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

// 其他路由
func registerOtherRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api")
	{
		otherGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
