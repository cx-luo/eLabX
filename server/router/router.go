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
	"eLabX/src/api/casbin"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 为了更好地管理和维护 API 路由，将不同模块的路由注册拆分到多个文件中。
// 例如：auth_routes.go, user_routes.go, casbin_routes.go, system_routes.go, other_routes.go
// 每个文件中定义对应的 registerXXXRoutes 函数，并在此处导入和调用。

// 例如：
// import "eLabX/router/auth_routes"
// import "eLabX/router/user_routes"
// import "eLabX/router/casbin_routes"
// import "eLabX/router/system_routes"
// import "eLabX/router/other_routes"

// 在 NewRouter 中调用：
// auth_routes.RegisterAuthRoutes(router)
// user_routes.RegisterUserRoutes(router)
// casbin_routes.RegisterCasbinRoutes(router)
// system_routes.RegisterSystemRoutes(router)
// other_routes.RegisterOtherRoutes(router)

// 这样可以实现 API 路由的模块化和解耦，便于后续维护和扩展。

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

	registerEtlRoutes(router)

	return router
}

func registerCasbinRoutes(r *gin.Engine) {
	casbinGroup := r.Group("/api/casbin")
	{
		casbinGroup.POST("/add", casbin.AddPolicy)
	}
}

// 其他路由
func registerOtherRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api")
	{
		otherGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
