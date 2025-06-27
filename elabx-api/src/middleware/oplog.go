// Package middleware coding=utf-8
// @Project : eLabX
// @Time    : 2024/2/5 11:05
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : oplog.go
// @Software: GoLand
package middleware

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
)

// APICallLog 结构表示API调用日志的数据模型
type APICallLog struct {
	StatusCode int    `json:"status_code" db:"status_code" gorm:"status_code"`
	Method     string `json:"method" db:"method" gorm:"method"`
	ApiPath    string `json:"api_path" db:"api_path" gorm:"api_path"`
	RemoteAddr string `json:"remote_addr" db:"remote_addr" gorm:"remote_addr"`
}

func Oplog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在API调用之前记录请求信息
		logEntry := APICallLog{
			ApiPath:    c.Request.URL.Path,
			Method:     c.Request.Method,
			RemoteAddr: c.RemoteIP(),
		}

		// 在API调用之后记录响应信息
		c.Next()

		// 记录响应状态码
		logEntry.StatusCode = c.Writer.Status()

		// 将日志存储到MySQL
		_, err := dao.OBCursor.Exec("INSERT INTO eln_access_records(api_path, method, status_code, remote_addr) VALUES (?, ?, ?, ?)",
			logEntry.ApiPath, logEntry.Method, logEntry.StatusCode, logEntry.RemoteAddr)
		if err != nil {
			utils.Logger.Error("Error inserting API call log:" + err.Error())
		}
	}
}
