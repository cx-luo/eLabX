// Package api coding=utf-8
// @Project : server
// @Time    : 2024/8/29 17:08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : closing_report.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"strings"
)

func GetProjectPageInfo(c *gin.Context) {
	var r struct {
		ProjectName string `json:"projectName,omitempty"`
		UserId      string `json:"userId,omitempty"`
	}
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	querySql := "select b.reaction_id from eln_project a join eln_project_page b on a.project_id = b.project_id where b.user_id = ?"
	if strings.TrimSpace(r.ProjectName) != "" {
		querySql = querySql + " and a.project_name = binary ?"
	}
	var rxnId []int64
	err = dao.OBCursor.Select(&rxnId, querySql, r.UserId, r.ProjectName)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	query, args, err := sqlx.In(`select b.project_name, c.page_name, a.reaction_id ,a.reaction_smiles ,a.cd_structure ,a.procedure_txt ,b.step_id, b.batch, b.start_date from eln_reaction_note a join eln_rxn_basicinfo b join eln_project_page c on a.reaction_id = c.reaction_id and a.reaction_id = b.reaction_id 
				where a.reaction_id in (?) order by a.reaction_id desc, b.step_id desc`, rxnId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	type ri struct {
		ProjectName    string   `json:"projectName,omitempty" db:"project_name"`
		PageName       string   `json:"pageName,omitempty" db:"page_name"`
		ReactionId     int64    `json:"reactionId" db:"reaction_id"`
		ReactionSmiles string   `json:"reactionSmiles" db:"reaction_smiles"`
		CdStructure    string   `json:"cdStructure,omitempty" db:"cd_structure"`
		ProcedureTxt   string   `json:"procedureTxt" db:"procedure_txt"`
		StepId         string   `json:"stepId" db:"step_id"`
		RxnImg         string   `json:"rxnImg" db:"rxn_img"`
		Batch          int      `json:"batch" db:"batch"`
		StartDate      string   `json:"startDate" db:"start_date"`
		Lcms           []string `json:"lcms"`
		Nmr            []string `json:"nmr"`
		Other          []string `json:"other"`
	}
	var rxnInfo []ri
	query = dao.OBCursor.Rebind(query)
	err = dao.OBCursor.Select(&rxnInfo, query, args...)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	for i, rxn := range rxnInfo {
		err := dao.OBCursor.Select(&rxnInfo[i].Lcms, `select spectrum_url from eln_rxn_report_lcms where reaction_id = ?`, rxn.ReactionId)
		err = dao.OBCursor.Select(&rxnInfo[i].Nmr, `select spectrum_url from eln_rxn_report_nmr where reaction_id = ?`, rxn.ReactionId)
		err = dao.OBCursor.Select(&rxnInfo[i].Other, `select spectrum_url from eln_rxn_report_other where reaction_id = ?`, rxn.ReactionId)
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
		var structure string
		if rxnInfo[i].CdStructure != "" {
			structure = rxnInfo[i].CdStructure
		} else {
			structure = rxnInfo[i].ReactionSmiles
		}
		err, rxnInfo[i].RxnImg = GenImgBase64String(structure)
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	}

	utils.SuccessWithData(c, "", gin.H{"list": rxnInfo})
	return
}
