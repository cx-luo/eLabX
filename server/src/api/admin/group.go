// Package admin coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/25 11:20
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : group.go
// @Software: GoLand
package admin

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GroupParam is used for creating and updating user groups
type GroupParam struct {
	GroupId     int64  `json:"groupId,omitempty"`
	GroupName   string `json:"groupName" binding:"required"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	CreatedBy   int64  `json:"createdBy"`
}

// PermissionParam is used for assigning permissions to a group
type PermissionParam struct {
	GroupId       int64   `json:"groupId" binding:"required"`
	PermissionIds []int64 `json:"permissionIds"`
}

// GetGroupList retrieves the list of user groups with pagination
func GetGroupList(c *gin.Context) {
	var req struct {
		Page      int    `json:"page,omitempty"`
		PageSize  int    `json:"pageSize,omitempty"`
		SortField string `json:"sortField,omitempty"`
		SortOrder string `json:"sortOrder,omitempty"`
		Status    int8   `json:"status,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize
	var groups []types.ElnUserGroup
	var total int64

	// Count total number of groups
	if err := dao.OBCursor.Model(&types.ElnUserGroup{}).Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	db := dao.OBCursor.Model(&types.ElnUserGroup{})
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if err := db.Limit(req.PageSize).Offset(offset).Find(&groups).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Success", gin.H{"items": groups, "total": total})
}

// GetGroupDetail retrieves the details of a single user group
func GetGroupDetail(c *gin.Context) {
	var param struct {
		GroupId int64 `json:"groupId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var group types.ElnUserGroup
	err := dao.OBCursor.Model(&types.ElnUserGroup{}).Where("group_id = ?", param.GroupId).First(&group).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", group)
}

// CreateGroup creates a new user group
func CreateGroup(c *gin.Context) {
	var param GroupParam
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	node, err := snowflake.NewNode(2)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	group := types.ElnUserGroup{
		GroupId:     node.Generate().Int64(),
		GroupName:   param.GroupName,
		Description: param.Description,
		Status:      1,
		CreatedBy:   param.CreatedBy,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}
	if err := dao.OBCursor.Model(&types.ElnUserGroup{}).Create(&group).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", group)
}

// UpdateGroup updates the information of a user group
func UpdateGroup(c *gin.Context) {
	var param types.ElnUserGroup
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	if err := dao.OBCursor.Where("group_id = ?", param.GroupId).Updates(&types.ElnUserGroup{
		GroupName:   param.GroupName,
		Description: param.Description,
		Status:      param.Status,
		CreatedBy:   param.CreatedBy,
		Permissions: param.Permissions,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundError(c, fmt.Errorf("group Not Found"))
			return
		}
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// DeleteGroup performs a soft delete on a user group (sets status to 0)
func DeleteGroup(c *gin.Context) {
	var param struct {
		GroupId int64 `json:"groupId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	if err := dao.OBCursor.Model(&types.ElnUserGroup{}).Where("group_id = ?", param.GroupId).Update("status", 0).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// AssignPermissionsToGroup assigns permissions to a user group
func AssignPermissionsToGroup(c *gin.Context) {
	var param PermissionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	// First, delete existing permissions for the group
	if err := dao.OBCursor.Where("group_id = ?", param.GroupId).Delete(&types.ElnGroupPermission{}).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// Batch insert new permissions
	var groupPerms []types.ElnGroupPermission
	for _, permId := range param.PermissionIds {
		groupPerms = append(groupPerms, types.ElnGroupPermission{
			GroupId:      param.GroupId,
			PermissionId: permId,
			CreateAt:     time.Now(),
		})
	}
	if len(groupPerms) > 0 {
		if err := dao.OBCursor.Model(&types.ElnGroupPermission{}).Create(&groupPerms).Error; err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}
	utils.Success(c, "Success")
}

// GetGroupPermissions retrieves the list of permissions for a user group
func GetGroupPermissions(c *gin.Context) {
	var param struct {
		GroupId int64 `json:"groupId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var perms []types.ElnGroupPermission
	err := dao.OBCursor.Table("eln_permission").
		Select("eln_permission.*").
		Joins("JOIN eln_group_permission ON eln_permission.permission_id = eln_group_permission.permission_id").
		Where("eln_group_permission.group_id = ?", param.GroupId).
		Find(&perms).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", perms)
}

// AddUserToGroup adds a user to a user group
func AddUserToGroup(c *gin.Context) {
	var param struct {
		GroupId int64 `json:"groupId" binding:"required"`
		UserId  int64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	// Check if the user is already in the group
	var count int64
	dao.OBCursor.Model(&types.ElnUserGroup{}).Where("group_id = ? AND user_id = ?", param.GroupId, param.UserId).Count(&count)
	if count > 0 {
		utils.BadRequestErr(c, errors.New("user already in group"))
		return
	}
	userGroup := types.ElnUserGroup{
		GroupId:  param.GroupId,
		UserId:   param.UserId,
		CreateAt: time.Now(),
	}
	if err := dao.OBCursor.Model(&types.ElnUserGroup{}).Create(&userGroup).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// RemoveUserFromGroup removes a user from a user group
func RemoveUserFromGroup(c *gin.Context) {
	var param struct {
		GroupId int64 `json:"groupId" binding:"required"`
		UserId  int64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	if err := dao.OBCursor.Where("group_id = ? AND user_id = ?", param.GroupId, param.UserId).Delete(&types.ElnUserGroup{}).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// GetGroupUsers retrieves the list of users in a user group
func GetGroupUsers(c *gin.Context) {
	var param struct {
		GroupId int64 `json:"groupId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var users []types.ElnUsers
	err := dao.OBCursor.Table("eln_user").
		Select("eln_user.*").
		Joins("JOIN eln_user_group ON eln_user.user_id = eln_user_group.user_id").
		Where("eln_user_group.group_id = ?", param.GroupId).
		Find(&users).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", users)
}
