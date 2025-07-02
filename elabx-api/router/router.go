// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2023/12/12 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@pharmaron.com
// @File    : router.go
// @Software: GoLand
package router

import (
	_ "eLabX/docs"
	"eLabX/src/api"
	"eLabX/src/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter returns a new router.
func NewRouter(outputPath string, loglevel string) *gin.Engine {
	// 设置全局 Logger
	//logger := middleware.SetupLogger(outputPath, loglevel)
	//
	//// 延迟关闭 logger
	//defer func(logger *zap.Logger) {
	//	err := logger.Sync()
	//	if err != nil {
	//		middleware.Logger.Error(err.Error())
	//	}
	//}(logger)

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

	registerRouteRoutes(router)
	// 注册其他路由
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
		userGroup.POST("/fetchUserName", api.FetchUserName)
		userGroup.GET("/getUserList", api.GetUserList)
		userGroup.POST("/changePwd", api.ChangePwd)
		userGroup.POST("/changeUsername", api.ChangeUserName)
		userGroup.POST("/forgetPwd", api.ForgetPwd)
		userGroup.POST("/resetPwd", api.ResetPwd)
	}
}

func registerCasbinRoutes(r *gin.Engine) {
	casbinGroup := r.Group("/api/casbin")
	{
		casbinGroup.POST("/add", api.CasbinAddPolicy)
	}
}

func registerRouteRoutes(r *gin.Engine) {
	routeGroup := r.Group("/api/route")
	{
		routeGroup.POST("/tree", api.GetRouteTree)
		routeGroup.GET("/list", api.GetUserRouteList)
	}
}

// 其他路由
func registerOtherRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api")
	{
		otherGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
