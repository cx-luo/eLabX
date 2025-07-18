// Package user coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/8 10:15
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : tables.go
// @Software: GoLand
package user

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
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

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 60
	}

	sortField := req.SortField
	sortOrder := req.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Build the ORDER BY clause
	orderBy := "user_id"
	if sortField != "" {
		orderBy = sortField
	}
	orderBy = orderBy + " " + sortOrder

	var users []types.ElnUsers
	query := dao.OBCursor.Table("eln_users").
		Select("user_id, username, phone, email, permissions, status, created_at").
		Order(orderBy).
		Offset((page - 1) * pageSize).
		Limit(pageSize)
	err := query.Find(&users).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// Get total count for pagination
	var total int64
	if err := dao.OBCursor.Table("eln_users").Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"total": total, "items": users})
	return
}

type route struct {
	BookName   string  `json:"bookName"`
	Icon       string  `json:"icon"`
	Path       string  `json:"path"`
	Layout     string  `json:"layout,omitempty"`
	PageName   []route `json:"pageName,omitempty"`
	ReactionId int64   `json:"reactionId,omitempty"`
}

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

type userId struct {
	WitnessId int `json:"witnessId"`
}

type pwdForm struct {
	UserID      int    `json:"userId"`
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

func ChangePwd(c *gin.Context) {
	var p pwdForm
	err := c.ShouldBind(&p)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var oldPwd string
	err = dao.OBCursor.Table("eln_users").Select("password_hash").
		Where("userid = ?", p.UserID).Find(&oldPwd).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	if oldPwd != p.OldPassword {
		utils.BadRequestErr(c, errors.New("the old password was incorrectly entered"))
		return
	}
	err = dao.OBCursor.Table("eln_users").Where("userid = ?", p.UserID).
		Update("password_hash", p.Password).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Password changed")
	return
}

type ForgetPwdForm struct {
	UserID   string `json:"userId"`
	SmsCode  string `json:"smsCode"`
	Password string `json:"password"`
}

// SendEmail sends an email
func SendEmail(c *gin.Context) {
	type EmailForm struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	var form EmailForm
	if err := c.ShouldBindJSON(&form); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err := utils.SendEmail(form.To, form.Subject, form.Body)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "Email sent successfully")
	return
}

// ForgetPwd resets the password via email
func ForgetPwd(c *gin.Context) {
	type EmailForgetPwdForm struct {
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var p EmailForgetPwdForm
	if err := c.ShouldBindJSON(&p); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	// Query the user ID corresponding to the email
	var user struct {
		UserID int    `json:"userId"`
		Email  string `json:"email"`
	}
	err := dao.OBCursor.Table("eln_users").Select("userid", "email").Where("email = ?", p.Email).First(&user).Error
	if err != nil {
		utils.BadRequestErr(c, errors.New("email not found"))
		return
	}

	// Query verification code and expiration time
	var record struct {
		VerifCode           string `json:"verifCode" db:"verif_code"`
		VerifCodeExpireTime int64  `json:"verifCodeExpireTime" db:"verif_code_expire_time"`
	}
	err = dao.OBCursor.Table("eln_register_records").
		Select("verif_code", "verif_code_expire_time").
		Where("userid = ?", user.UserID).
		First(&record).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	now := time.Now().Unix()
	if record.VerifCode != p.Code {
		utils.BadRequestErr(c, errors.New("invalid verification code"))
		return
	}
	if now > record.VerifCodeExpireTime {
		utils.BadRequestErr(c, errors.New("verification code expired"))
		return
	}

	// Update password
	err = dao.OBCursor.Table("eln_users").
		Where("userid = ?", user.UserID).
		Update("password_hash", p.Password).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "Password reset successfully")
	return
}

func ModifyUserInfo(c *gin.Context) {
	var req struct {
		UserId   int    `json:"userId" binding:"required"`
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	if err := dao.OBCursor.Table("eln_users").
		Where("userid = ?", req.UserId).
		Update("username", req.Username).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "Username changed")
	return
}

func ResetPwd(c *gin.Context) {
	var req struct {
		RequestId int `json:"requestId"`
		UserID    int `json:"userId"`
	}
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	var permission string
	if err := dao.OBCursor.Table("eln_users").
		Select("permissions").
		Where("userid = ?", req.RequestId).
		Find(&permission).Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if !strings.Contains(permission, "admin") {
		// Permission denied, handled after selection
		return
	}

	if err := dao.OBCursor.Table("eln_users").
		Where("userid = ?", req.RequestId).
		Update("password_hash", "e10adc3949ba59abbe56e057f20f883e").Error; err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "change passcode successfully")
	return
}

func SetUserAuthorities(c *gin.Context) {
	var req struct {
		UserId       int    `json:"userId" binding:"required"`
		AuthorityIds string `json:"authorityIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	err := dao.OBCursor.Table("eln_users").Where("user_id = ?", req.UserId).Update("permissions", req.AuthorityIds).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "User authorities updated successfully")
	return
}
