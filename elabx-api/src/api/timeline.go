// Package api coding=utf-8
// @Project : server
// @Time    : 2024/8/20 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : timeline.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
)

type opInfo struct {
	ReactionId uint64 `json:"reactionId" db:"reaction_id"`
	UserId     uint   `json:"userId" db:"user_id"`
	UserName   string `json:"userName" db:"user_name"`
	Status     string `json:"status" db:"status"`
	Project    string `json:"project" db:"project"`
	Comments   string `json:"comments" db:"comments"`
}

func ReopenRxn(c *gin.Context) {
	var op opInfo
	err := c.ShouldBind(&op)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	_, err = dao.OBCursor.NamedExec(`insert into eln_commit_log(reaction_id, project, user_id, user_name, status ) VALUE (:reaction_id,:project,:user_id,:user_name,:status)`, op)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`update eln_rxn_basicinfo set rxn_status = 'open' where reaction_id = ?`, op.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}
