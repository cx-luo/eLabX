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
)

// GetPermissionList 获取权限列表
func GetPermissionList(c *gin.Context) {
	var req struct {
		Query string `json:"query,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	var permissions []types.ElnPermission
	var total int64

	db := dao.OBCursor.Model(&types.ElnPermission{})
	if req.Query != "" {
		db = db.Where("description LIKE ?", "%"+req.Query+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if err := db.Find(&permissions).Error; err != nil {
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
