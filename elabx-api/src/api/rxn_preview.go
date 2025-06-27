// Package api coding=utf-8
// @Project : server
// @Time    : 2024/6/28 9:53
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : rxn_preview.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
)

func GetRxnForAuthority(c *gin.Context) {
	var user struct {
		IsQc      bool `json:"isQc"`
		UserId    int  `json:"userId"`
		ChemistId int  `json:"chemistId"`
	}
	err := c.ShouldBind(&user)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var authProject []struct {
		ReactionId int64  `json:"reactionId" db:"reaction_id"`
		PageName   string `json:"pageName" db:"page_name"`
	}

	// 如果是QC，可以看别人的，如果不是，则只能看自己的，不论你输入的chemistID是啥
	if user.IsQc {
		err = dao.OBCursor.Select(&authProject, `select epp.reaction_id,epp.page_name from eln_project_page epp left join eln_rxn_basicinfo erb on epp.reaction_id = erb.reaction_id 
                                     where erb.rxn_status != 'open' and epp.user_id = ?`, user.ChemistId)
	} else {
		err = dao.OBCursor.Select(&authProject, `select epp.reaction_id,epp.page_name from eln_project_page epp left join eln_rxn_basicinfo erb on epp.reaction_id = erb.reaction_id 
                                     where erb.rxn_status != 'open' and erb.witness_id = ?`, user.UserId)
	}
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"rxnIds": authProject})
	return
}
