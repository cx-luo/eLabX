// Package casbin coding=utf-8
// @Project : eLabX
// @Time    : 2025/6/30 16:53
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : casbin_manager.go
// @Software: GoLand
package casbin

import (
	"eLabX/middleware"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
)

func AddPolicy(c *gin.Context) {
	var roles struct {
		RoleName  string `json:"roleName,omitempty" `
		ApiPath   string `json:"apiPath,omitempty"`
		ApiMethod string `json:"apiMethod"`
		Action    string `json:"action"`
	}
	err := c.ShouldBind(&roles)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	_, err = middleware.GlobalCasBin.AddPolicy(roles.RoleName, roles.ApiPath, roles.ApiMethod, roles.Action)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func AddRoleForUser(c *gin.Context) {
	var user struct {
		UserId   string `json:"userId"`
		RoleName string `json:"roleName"`
	}
	err := c.ShouldBind(&user)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = middleware.GlobalCasBin.AddRolesForUser(user.UserId, []string{user.RoleName})
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}
