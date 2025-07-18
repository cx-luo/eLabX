// Package admin coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/17 16:46
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : project.go
// @Software: GoLand
package admin

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"time"
)

// ProjectParam 用于创建和更新项目的参数
type ProjectParam struct {
	ProjectId   int64  `json:"projectId,omitempty"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
}

// GetProjectList 获取项目列表
func GetProjectList(c *gin.Context) {
	// 增加分页功能
	var req struct {
		Page      int    `json:"page,omitempty"`
		PageSize  int    `json:"pageSize,omitempty"`
		SortField string `json:"sortField,omitempty"`
		SortOrder string `json:"sortOrder,omitempty"`
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
	var projects []types.ElnProject
	var total int64

	// 先统计总数
	if err := dao.OBCursor.Model(&types.ElnProject{}).Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 分页查询
	if err := dao.OBCursor.Model(&types.ElnProject{}).Limit(req.PageSize).Offset(offset).Find(&projects).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Success", gin.H{"items": projects, "total": total})
}

// GetProjectDetail 获取单个项目详情
func GetProjectDetail(c *gin.Context) {
	var param struct {
		ProjectId int64 `json:"projectId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var project types.ElnProject
	err := dao.OBCursor.Model(&types.ElnProject{}).Where("project_id = ?", param.ProjectId).First(&project).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", project)
}

// CreateProject 新建项目
func CreateProject(c *gin.Context) {
	var param ProjectParam
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	project := types.ElnProject{
		ProjectId:   node.Generate().Int64(),
		ProjectName: param.Name,
		Description: param.Description,
		Status:      param.Status,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}
	if err := dao.OBCursor.Model(&types.ElnProject{}).Create(&project).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "Success", project)
}

// UpdateProject 更新项目信息
func UpdateProject(c *gin.Context) {
	var param ProjectParam
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	update := map[string]interface{}{
		"name":        param.Name,
		"description": param.Description,
		"status":      param.Status,
		"update_at":   time.Now(),
	}
	if err := dao.OBCursor.Model(&types.ElnProject{}).Where("project_id = ?", param.ProjectId).Updates(update).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	var param struct {
		ProjectId int64 `json:"projectId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	// 将项目的 status 字段设置为 0（软删除）
	if err := dao.OBCursor.Model(&types.ElnProject{}).Where("project_id = ?", param.ProjectId).Update("status", 0).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}
