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
	err := dao.OBCursor.Where("status = 1 and user_id = ?", username).Find(&user).Error
	if err != nil {
		utils.NotFoundError(c, fmt.Errorf("User does not exist or is unavailable. %s\n", err))
		return
	}

	// Get group IDs from user.GroupId (assuming it's a string of comma-separated IDs)
	var userGroups []types.ElnUserGroup
	err = dao.OBCursor.Model(&types.ElnUserGroup{}).Where("user_id = ? AND status = 1", user.UserId).Find(&userGroups).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	var groupIds []string
	for _, ug := range userGroups {
		groupIds = append(groupIds, fmt.Sprintf("%d", ug.GroupId))
	}

	// Query permissions for these group IDs
	var permissions []string
	var permissionIds []int64
	if len(groupIds) > 0 && groupIds[0] != "" {
		// First, get permissions string(s) for the user's groups
		var groupPermissions []string
		err := dao.OBCursor.Model(&types.ElnGroup{}).Select("permissions").Where("group_id IN (?)", groupIds).Pluck("permissions", &groupPermissions).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
		// Collect all permission IDs from all groups
		permissionIdSet := make(map[int64]struct{})
		for _, perms := range groupPermissions {
			for _, pidStr := range strings.Split(perms, ",") {
				pidStr = strings.TrimSpace(pidStr)
				if pidStr == "" {
					continue
				}
				var pid int64
				_, err := fmt.Sscan(pidStr, &pid)
				if err == nil {
					permissionIdSet[pid] = struct{}{}
				}
			}
		}
		permissionIds = make([]int64, 0, len(permissionIdSet))
		for pid := range permissionIdSet {
			permissionIds = append(permissionIds, pid)
		}
		// Now, get permission names
		if len(permissionIds) > 0 {
			err = dao.OBCursor.Model(&types.ElnPermission{}).Where("permission_id in (?)", permissionIds).Pluck("permission_name", &permissions).Error
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
		}
	}

	utils.SuccessWithData(c, "success", gin.H{
		"permissions": permissions,
		"username":    user.Username,
	})
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
		Query    string `json:"query,omitempty"`
		Status   int    `json:"status,omitempty"`
		Page     int    `json:"page,omitempty"`
		PageSize int    `json:"pageSize,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	query := req.Query
	var users []types.ElnUsers
	var total int64

	db := dao.OBCursor.Table("eln_users").Select("user_id, username, email, status, ip_addr, create_at, update_at").Where("status = 1")
	if err := db.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if len(query) != 0 {
		db.Where("username LIKE ?", "%"+query+"%")
	}
	var page, pageSize, offset int
	if req.Page > 0 {
		page = req.Page
		if page <= 0 {
			page = 1
		}
		pageSize = req.PageSize
		if pageSize <= 0 {
			pageSize = 60
		}
		offset = (page - 1) * pageSize
		db.Offset(offset).Limit(pageSize)
	}

	if err := db.Find(&users).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{
		"items": users,
		"total": total,
	})
}
