// Package system coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/18 16:49
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : users.go
// @Software: GoLand
package system

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

// GetSystemUserList handles /api/system/user/list
func GetSystemUserList(c *gin.Context) {
	var req struct {
		Page     int `json:"page,omitempty"`
		PageSize int `json:"pageSize,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 60
	}
	offset := (page - 1) * pageSize

	var users []types.ElnUsers
	var total int64

	db := dao.OBCursor.Table("eln_users").
		Select("user_id, username, email, status, ip_addr, created_at, update_at")

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// Get paginated data
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{
		"items":    users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
