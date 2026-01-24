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
	"eLabX/middleware"

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

	// Idempotency middleware for idempotent POST/PUT/PATCH requests
	// Note: This middleware is only effective for HTTP requests with an Idempotency-Key header.
	router.Use(middleware.GinIdempotencyMiddleware())

	router.Use(middleware.CasbinMiddleware())

	registerAuthRoutes(router)
	// 注册用户相关路由
	registerUserRoutes(router)

	// 注册其他路由
	registerSystemRoutes(router)
	// 注册实验笔记相关路由
	registerEnoteRoutes(router)
	// 注册其他路由
	registerOtherRoutes(router)
	// 注册ETL路由
	registerEtlRoutes(router)
	// 注册管理员路由
	registerAdminRoutes(router)
	return router
}

// 其他路由
func registerOtherRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api")
	{
		otherGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
