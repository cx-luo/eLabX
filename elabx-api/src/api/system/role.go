// Package system coding=utf-8
// @Project : eLabX
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
	"github.com/gin-gonic/gin"
	"strings"
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
	var role struct {
		ID     int64   `json:"id"`
		AuthId []int64 `json:"authId"`
		ApiId  []int64 `json:"apiId"`
	}
	err := c.ShouldBind(&role)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Model(&types.ElnRoles{}).Where(`id = ?`, role.ID).Updates(map[string]interface{}{
		"auth_id": strings.Join(utils.Int64SliceToStringSlice(role.AuthId), ","),
		"api_id":  strings.Join(utils.Int64SliceToStringSlice(role.ApiId), ","),
	}).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
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
	var role types.ElnRoles
	err := dao.OBCursor.Select(`id,name,status,auth_id,api_id`).Where(`id = ?`, id).First(&role).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	roleAuthId := strings.Split(role.AuthId, ",")
	roleApiId := strings.Split(role.ApiId, ",")
	utils.SuccessWithData(c, "", gin.H{"apiId": roleApiId, "id": role.ID, "authId": roleAuthId, "name": role.Name, "status": role.Status})
	return
}

func UpdateRole(c *gin.Context) {
	var role types.ElnRoles
	err := c.ShouldBind(&role)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Model(&types.ElnRoles{}).Where(`id = ?`, role.ID).Updates(&role).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

func DeleteRole(c *gin.Context) {
	var role struct {
		ID int `json:"id"`
	}
	err := c.ShouldBind(&role)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Where(`id = ?`, role.ID).Delete(&types.ElnRoles{}).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}
