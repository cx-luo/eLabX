// Package api coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/28 12:00
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : auth.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/middleware"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
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
	err = dao.OBCursor.Table("eln_users").Select("password_hash").
		Where(`status = 1 AND userid = ?`, u.Username).Find(&passwordHash).Error

	if err != nil {
		utils.Logger.Error(err.Error())
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"msg": "User authentication failed, username or password incorrect."})
		return
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

func SetUserAuthorities(c *gin.Context) {
	var roles struct {
		Userid       int    `json:"userid,omitempty" db:"userid"`
		AuthorityIds string `json:"authorityIds,omitempty" db:"permissions"`
	}
	err := c.ShouldBind(&roles)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	err = dao.OBCursor.Exec(`update eln_users set permissions = ? where userid = ?`, roles.AuthorityIds, roles.Userid).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}
