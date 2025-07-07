// Package system coding=utf-8
// @Project : elabx-api
// @Time    : 2025/7/7 14:23
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : role.go
// @Software: GoLand
package system

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Role struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
	Name     string `json:"name,omitempty"`
	Status   int    `json:"status,omitempty"`
}

func GetRoleList(c *gin.Context) {
	var role Role
	err := c.ShouldBind(&role)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	var roles []types.ElnRoles
	err = dao.OBCursor.Limit(role.PageSize).Offset((role.Page - 1) * role.PageSize).Find(&roles).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"items": roles})
	return
}

func RoleAssign(c *gin.Context) {
	type roleAssign struct {
		ID     int64   `json:"id"`
		AuthID []int64 `json:"authId"`
		ApiID  []int64 `json:"apiId"`
	}
	var role roleAssign
	err := c.ShouldBind(&role)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func RoleAdd(c *gin.Context) {
	var role types.ElnRoles
	err := c.ShouldBind(&role)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	err = dao.OBCursor.Create(&role).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func RoleInfo(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	utils.SuccessWithData(c, "", gin.H{"apiId": []int64{2698}})
	return
}
