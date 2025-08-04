// Package admin coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/30 14:50
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : permissions.go
// @Software: GoLand
package admin

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPermissionList 获取权限列表
func GetPermissionList(c *gin.Context) {
	var req struct {
		Page     int    `json:"page,omitempty"`
		PageSize int    `json:"pageSize,omitempty"`
		Status   int8   `json:"status,omitempty"`
		Query    string `json:"query,omitempty"`
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

	var permissions []types.ElnPermission
	var total int64

	db := dao.OBCursor.Model(&types.ElnPermission{})
	if req.Query != "" {
		db = db.Where("permission_name LIKE ? OR description LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if err := db.Offset(offset).Limit(req.PageSize).Find(&permissions).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Success", gin.H{
		"items": permissions,
		"total": total,
	})
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var req types.ElnPermission
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	// 检查权限名是否已存在
	var count int64
	if err := dao.OBCursor.Model(&types.ElnPermission{}).Where("permission_name = ?", req.PermissionName).Count(&count).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	if count > 0 {
		utils.BadRequestErr(c, errors.New("permission name already exists"))
		return
	}

	if err := dao.OBCursor.Create(&req).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Permission created successfully", req)
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	var req types.ElnPermission
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	if req.PermissionId == 0 {
		utils.BadRequestErr(c, errors.New("permission ID is required"))
		return
	}

	// 检查权限是否存在
	var perm types.ElnPermission
	if err := dao.OBCursor.First(&perm, req.PermissionId).Error; err != nil {
		utils.BadRequestErr(c, errors.New("permission not found"))
		return
	}

	// 检查权限名是否冲突
	if req.PermissionName != "" && req.PermissionName != perm.PermissionName {
		var count int64
		if err := dao.OBCursor.Model(&types.ElnPermission{}).Where("permission_name = ? AND permission_id <> ?", req.PermissionName, req.PermissionId).Count(&count).Error; err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
		if count > 0 {
			utils.BadRequestErr(c, errors.New("permission name already exists"))
			return
		}
	}

	if err := dao.OBCursor.Model(&perm).Updates(req).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Permission updated successfully", perm)
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	if req.ID == 0 {
		utils.BadRequestErr(c, errors.New("permission ID is required"))
		return
	}

	if err := dao.OBCursor.Delete(&types.ElnPermission{}, req.ID).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "Permission deleted successfully")
}

// GetPermsListWithFilterApi getPermsListWithFilterApi 获取权限列表（带筛选）
func GetPermsListWithFilterApi(c *gin.Context) {
	var req struct {
		Query string `json:"query" form:"query"`
	}
	// Accept both JSON and form for flexibility
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid query params: "+err.Error()))
		return
	}

	var perms []struct {
		PermissionId   int64  `json:"permissionId" db:"permission_id" gorm:"primaryKey;column:permission_id"`
		PermissionName string `json:"permissionName" db:"permission_name" gorm:"column:permission_name"`
		Description    string `json:"description" db:"description" gorm:"column:description"`
	}

	queryBuilder := dao.OBCursor.Table("eln_permission").Select("permission_id, permission_name, description").Where("status = ?", 1)

	// Use COLLATE for case-insensitive search if needed, and escape % and _ in LIKE
	if req.Query != "" {
		likeQuery := "%" + req.Query + "%"
		queryBuilder = queryBuilder.Where("permission_name LIKE ?", likeQuery)
	}

	var total int64
	// Count must be done on a new session to avoid side effects of Order/Limit
	if err := queryBuilder.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if err := queryBuilder.Order("permission_id desc").Find(&perms).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Permissions fetched successfully", gin.H{
		"items": perms,
		"total": total,
	})
}
