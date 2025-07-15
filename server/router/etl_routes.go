// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/9 11:53
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : etl_routes.go
// @Software: GoLand
package router

import (
	"eLabX/src/api/etl"
	"github.com/gin-gonic/gin"
)

func registerEtlRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/etl")
	{
		userGroup.GET("/database/list", etl.GetDatabase)
		userGroup.GET("/table/list/:dbName", etl.GetTableList)
		userGroup.GET("/table/columns/:dbName/:tableName", etl.GetTableColumnsWithPK)
		userGroup.POST("/table/data/:dbName/:tableName", etl.GetTableData)
		userGroup.POST("/table/data/update", etl.UpdateTableDataApi)
	}
}
