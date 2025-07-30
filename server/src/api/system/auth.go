// Package system coding=utf-8
// @Project : eLabX
// @Time    : 2025/6/28 12:00
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : auth.go
// @Software: GoLand
package system

import (
	"eLabX/middleware"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func UserLogin(c *gin.Context) {
	type user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var u user
	err := c.ShouldBind(&u)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}

	var passwordHash string
	// Query the password hash for the given username and status=1
	err = dao.OBCursor.Table("eln_users").Select("password_hash").
		Where("status = 1 AND user_id = ?", u.Username).Take(&passwordHash).Error
	if err != nil {
		zap.L().Error("Failed to find user or user inactive", zap.String("user_id", u.Username), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User authentication failed, username or password incorrect."})
		return
	}

	// Record login IP to table (do not block login on failure)
	loginIP := c.ClientIP()
	errUpdate := dao.OBCursor.Table("eln_users").
		Where("user_id = ?", u.Username).
		Update("ip_addr", loginIP).Error
	if errUpdate != nil {
		zap.L().Warn("Failed to update login IP", zap.String("user_id", u.Username), zap.Error(errUpdate))
	}

	if passwordHash == u.Password {
		token, _ := middleware.GenToken(u.Username, u.Password)
		utils.SuccessWithData(c, "", gin.H{"accessToken": "Bearer " + token})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User authentication failed, username or password incorrect."})
	}
}

func UserLogout(c *gin.Context) {
	utils.Success(c, "Logged out successfully")
	return
}
