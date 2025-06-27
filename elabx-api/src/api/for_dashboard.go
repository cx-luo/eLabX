// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/2/5 19:59
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : for_dashboard.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRecentUsage(c *gin.Context) {
	var usageInfo []struct {
		OrderDay        string `json:"order_day" db:"order_day"`
		DailyOrderCount int    `json:"daily_order_count" db:"daily_order_count"`
	}
	err := dao.OBCursor.Select(&usageInfo, `SELECT DATE(access_time) AS order_day, COUNT(*) AS daily_order_count FROM eln_access_records GROUP BY order_day ORDER BY order_day desc limit 30`)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "success", Data: gin.H{"total": len(usageInfo), "output": usageInfo},
	})
}
