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
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.BadRequestErr(c, errors.New("User does not exist or is unavailable.\n"))
		return
	}

	var user types.ElnUsers
	err := dao.OBCursor.Table("eln_users").Select("username", "roles", "permissions").
		Where("status = 1 and user_id = ?", username).Find(&user).Error
	if err != nil {
		utils.NotFoundError(c, fmt.Errorf("User does not exist or is unavailable. %s\n", err))
		return
	}

	utils.SuccessWithData(c, "success", gin.H{"permissions": strings.Split(user.Permissions, ","), "username": user.Username})
	return
}

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

	// Use a base query for count, and a separate query for paginated data to avoid GORM's Count+Limit/Offset bug.
	baseQuery := dao.OBCursor.Table("eln_users")

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// Get paginated data
	if err := baseQuery.
		Select("user_id, username, email, status, ip_addr, create_at, update_at").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{
		"items": users,
		"total": total,
	})
}

// GetUserListWithFilter handles /api/system/user/filter
func GetUserListWithFilter(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	query := req.Query
	if len(query) == 0 {
		utils.SuccessWithData(c, "", gin.H{
			"items": []types.ElnUsers{},
			"total": 0,
		})
		return
	}

	var users []types.ElnUsers
	var total int64

	db := dao.OBCursor.Table("eln_users").Where("username LIKE ?", "%"+query+"%").Where("status = 1")

	if err := db.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if err := db.
		Select("user_id, username, email, status, ip_addr, create_at, update_at").
		Limit(30).
		Find(&users).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{
		"items": users,
		"total": total,
	})
}
