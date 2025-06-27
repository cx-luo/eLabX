// Package api coding=utf-8
// @Project : server
// @Time    : 2024/6/18 9:17
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : authority.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ElnUserAuthority struct {
	AuthorityId   int64     `json:"authorityId" db:"authority_id"`
	AuthorityName string    `json:"authorityName" db:"authority_name"`
	Status        int8      `json:"status" db:"status"`
	CreatedAt     time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

func GetAuthorityList(c *gin.Context) {
	var p struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	err := c.ShouldBind(&p)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	var roles []ElnUserAuthority
	err = dao.OBCursor.Select(&roles, `select * from eln_user_authority`)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"total": len(roles), "page": p.Page, "pageSize": p.PageSize, "roles": roles})
	return
}

func DeleteAuthority(c *gin.Context) {
	var p struct {
		AuthorityId int `json:"authorityId"`
	}
	err := c.ShouldBind(&p)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`delete from eln_user_authority where authority_id = ?`, p.AuthorityId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

func DisableAuthority(c *gin.Context) {
	var x struct {
		AuthorityID int `json:"authorityId,omitempty" db:"authority_id"`
		Status      int `json:"status,omitempty" db:"authorityId"`
	}
	err := c.ShouldBind(&x)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`update eln_user_authority set status = ? where authority_id = ?`, x.Status, x.AuthorityID)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

func CopyAuthority(c *gin.Context) {
	var x struct {
		AuthorityID   int    `json:"authorityId,omitempty" db:"authority_id"`
		AuthorityName string `json:"authorityName,omitempty" db:"authority_name"`
	}
	err := c.ShouldBind(&x)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`insert into eln_user_authority(authority_name) value (?)`, x.AuthorityName)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

func GetApiList(c *gin.Context) {
	type api struct {
		ApiGroup    string `json:"apiGroup"`
		Path        string `json:"path"`
		Description string `json:"description"`
		Method      string `json:"method"`
	}
	var apis []api
	for _, a := range utils.Apis {
		ap := strings.Split(a.Path, "/")[2]
		apis = append(apis, api{
			ApiGroup:    ap,
			Path:        a.Path,
			Description: a.Handler,
			Method:      a.Method,
		})
	}
	utils.SuccessWithData(c, "", gin.H{"apis": apis})
	return
}

func AuthorityStatusChange(c *gin.Context) {
	var action struct {
		ReactionId int64  `json:"reactionId"`
		Status     string `json:"status"`
		Comments   string `json:"comments"`
	}
	err := c.ShouldBind(&action)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`update eln_rxn_basicinfo set rxn_status = ?, comment =? where reaction_id = ?`, action.Status, action.Comments, action.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}
