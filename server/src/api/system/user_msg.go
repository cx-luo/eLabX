// Package system coding=utf-8
// @Project : eLabX
// @Time    : 2025/8/1 16:04
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : user_msg.go
// @Software: GoLand
package system

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserMessages retrieves a paginated list of messages for a user
func GetUserMessages(c *gin.Context) {
	var req struct {
		UserId   int64  `json:"userId" binding:"required"`
		Page     int    `json:"page,omitempty"`
		PageSize int    `json:"pageSize,omitempty"`
		Type     string `json:"type,omitempty"`
		IsRead   *bool  `json:"isRead,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	offset := (req.Page - 1) * req.PageSize

	db := dao.OBCursor.Model(&types.ElnUserMessage{}).Where("user_id = ?", req.UserId)
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.IsRead != nil {
		db = db.Where("is_read = ?", *req.IsRead)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	var messages []types.ElnUserMessage
	if err := db.Order("create_at desc").Limit(req.PageSize).Offset(offset).Find(&messages).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "Success", gin.H{
		"items": messages,
		"total": total,
	})
}

// MarkMessageAsRead marks a specific message as read
func MarkMessageAsRead(c *gin.Context) {
	var req struct {
		UserId int64 `json:"userId" binding:"required"`
		MsgId  int64 `json:"msgId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}
	if err := dao.OBCursor.Model(&types.ElnUserMessage{}).
		Where("id = ? AND user_id = ?", req.MsgId, req.UserId).
		Update("is_read", true).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// MarkAllMessagesAsRead marks all messages as read for a user
func MarkAllMessagesAsRead(c *gin.Context) {
	var req struct {
		UserId int64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}
	if err := dao.OBCursor.Model(&types.ElnUserMessage{}).
		Where("user_id = ? AND is_read = ?", req.UserId, false).
		Update("is_read", true).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Success")
}

// CreateUserMessage creates a new message/notification for a user
func CreateUserMessage(userId int64, title, content, msgType string) error {
	msg := types.ElnUserMessage{
		UserId:   userId,
		Title:    title,
		Content:  content,
		IsRead:   false,
		Type:     msgType,
		CreateAt: time.Now(),
	}
	if err := dao.OBCursor.Model(&types.ElnUserMessage{}).Create(&msg).Error; err != nil {
		return fmt.Errorf("failed to create user message: %w", err)
	}
	return nil
}
