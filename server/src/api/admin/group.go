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
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Name      string `json:"name,omitempty"`
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
	var groups []types.ElnGroup
	var total int64

	// Count total number of groups
	if err := dao.OBCursor.Model(&types.ElnGroup{}).Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	db := dao.OBCursor.Model(&types.ElnGroup{})
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Name != "" {
		db = db.Where("group_name LIKE ?", "%"+req.Name+"%")
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
	var group types.ElnGroup
	err := dao.OBCursor.Model(&types.ElnGroup{}).Where("group_id = ?", param.GroupId).First(&group).Error
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

	group := types.ElnGroup{
		GroupId:     utils.GenRandInt(10),
		GroupName:   param.GroupName,
		Description: param.Description,
		Status:      1,
		CreatedBy:   param.CreatedBy,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}
	if err := dao.OBCursor.Model(&types.ElnGroup{}).Create(&group).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", group)
}

// UpdateGroup updates the information of a user group
func UpdateGroup(c *gin.Context) {
	var param types.ElnGroup
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	if err := dao.OBCursor.Where("group_id = ?", param.GroupId).Updates(&types.ElnGroup{
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
	if err := dao.OBCursor.Model(&types.ElnGroup{}).Where("group_id = ?", param.GroupId).Update("status", 0).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// AssignPermissionsToGroup assigns permissions to a user group.
// One user can have several groups, and each group can have multiple permissions.
func AssignPermissionsToGroup(c *gin.Context) {
	// PermissionParam should contain Groups and PermissionIds ([]int64)
	var param struct {
		GroupId       int64   `json:"groupId" binding:"required"`
		PermissionIds []int64 `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	// Convert permission IDs to comma-separated string
	var permStrs []string
	for _, permId := range param.PermissionIds {
		permStrs = append(permStrs, strconv.FormatInt(permId, 10))
	}
	permissions := strings.Join(permStrs, ",")

	// Update the permissions field for the group
	if err := dao.OBCursor.Model(&types.ElnGroup{}).
		Where("group_id = ?", param.GroupId).
		Update("permissions", permissions).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
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
	// Get the group to fetch its permissions field
	var group types.ElnGroup
	if err := dao.OBCursor.Model(&types.ElnGroup{}).Where("group_id = ?", param.GroupId).First(&group).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	// If no permissions, return empty array
	if strings.TrimSpace(group.Permissions) == "" {
		utils.SuccessWithData(c, "Success", gin.H{"items": []types.ElnPermission{}})
		return
	}
	// Parse permissions string to []int64
	permStrs := strings.Split(group.Permissions, ",")
	var permIds []int64
	for _, s := range permStrs {
		if id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
			permIds = append(permIds, id)
		}
	}
	if len(permIds) == 0 {
		utils.SuccessWithData(c, "Success", gin.H{"items": []types.ElnPermission{}})
		return
	}
	var perms []types.ElnPermission
	if err := dao.OBCursor.Model(&types.ElnPermission{}).Where("permission_id IN ?", permIds).Find(&perms).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", gin.H{"items": perms})
}

// AddUserToGroup adds a user to a user group
func AddUserToGroup(c *gin.Context) {
	var param struct {
		GroupId int64   `json:"groupId" binding:"required"`
		UserId  []int64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	// Add each user to the group in eln_user_group
	for _, userId := range param.UserId {
		userGroup := types.ElnUserGroup{
			UserId:    userId,
			GroupId:   param.GroupId,
			Status:    1,
			CreatedBy: 0, // Optionally set CreatedBy if available
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		}
		// Use Create with OnConflict to avoid duplicate primary key error
		if err := dao.OBCursor.Clauses(clause.OnConflict{DoNothing: true}).Create(&userGroup).Error; err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}

	utils.Success(c, "Success")
}

// RemoveUserFromGroup removes a user from a user group
func RemoveUserFromGroup(c *gin.Context) {
	var param struct {
		GroupId int64   `json:"groupId" binding:"required"`
		UserId  []int64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	if err := dao.OBCursor.Where("group_id = ? AND user_id in (?)", param.GroupId, param.UserId).Delete(&types.ElnUserGroup{}).Error; err != nil {
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

	var userGroups []types.ElnUserGroup
	if err := dao.OBCursor.Where("group_id = ? AND status = 1", param.GroupId).Find(&userGroups).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if len(userGroups) == 0 {
		utils.SuccessWithData(c, "Success", gin.H{"items": []types.ElnUsers{}})
		return
	}

	userIds := make([]int64, 0, len(userGroups))
	for _, ug := range userGroups {
		userIds = append(userIds, ug.UserId)
	}

	var users []types.ElnUsers
	if err := dao.OBCursor.Where("user_id IN (?) AND status = 1", userIds).Find(&users).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Success", gin.H{"items": users})
}
