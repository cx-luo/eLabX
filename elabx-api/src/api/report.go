// Package api coding=utf-8
// @Project : server
// @Time    : 2025/3/19 16:59
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : report.go
// @Software: GoLand
package api

import (
	"database/sql"
	"eLabX/src/common"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func SignatureReport(c *gin.Context) {
	var rxn struct {
		ReactionId int64 `json:"reactionId" db:"reaction_id"`
	}
	err := c.ShouldBind(&rxn)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var basicInfo common.ElnRxnBasicInfo
	err = dao.OBCursor.Get(&basicInfo, `select a.reaction_id, a.project_name, a.batch, a.step_id, a.author_id, a.author_name, a.witness_name, a.rxn_type, 
       a.reference, a.start_date, a.commit_date, b.page_name from eln_rxn_basicinfo a join eln_project_page b on 
           a.reaction_id = b.reaction_id and a.reaction_id = ?`, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, errors.New("get ElnRxnBasicInfo error:"+err.Error()))
		return
	}

	var reagents []common.ReagentInfo
	err = dao.OBCursor.Select(&reagents, `select reagent_id, reagent_name, formula, mw, quantity, quantity_unit, 
       equiv, moles, moles_unit, purity, density from eln_rxn_reagents where reaction_id = ? and reagent_role != 'product' and reagent_role != 'solvent'`, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, errors.New("get reagents error:"+err.Error()))
		return
	}

	var solvents []common.ReagentInfo
	err = dao.OBCursor.Select(&solvents, `select reagent_id, reagent_name, formula, mw, volume, volume_unit, 
       equiv, moles, moles_unit, purity, density from eln_rxn_reagents where reaction_id = ? and reagent_role = 'solvent'`, rxn.ReactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 处理没有找到记录的情况
			solvents = nil
		} else {
			utils.InternalRequestErr(c, errors.New("get solvents error:"+err.Error()))
			return
		}
	}
	//var products []common.ReagentInfo
	//err = dao.OBCursor.Select(&products, `select reagent_id, reagent_name, formula, mw, quantity, quantity_unit,
	//   purity, yield from eln_rxn_reagents where reaction_id = ? and reagent_role = 'product'`, rxn.ReactionId)
	//if err != nil {
	//	utils.InternalRequestErr(c, err)
	//	return
	//}

	var sps []common.Samples
	err = dao.OBCursor.Select(&sps, `select product_alias ,reagent_name ,mf ,mass , synthesised ,synthesised_unit,
       purity ,sample_yield  from eln_rxn_samples b where b.reaction_id = ?`, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, errors.New("get Samples error:"+err.Error()))
		return
	}

	type rxnInfo struct {
		CdStructure      string `json:"cdStructure" db:"cd_structure"`
		Procedure        string `json:"procedure" db:"procedure_txt"`
		GetTargetProduct int    `json:"getTargetProduct" db:"get_target_product"`
	}
	var structProcedure rxnInfo
	err = dao.OBCursor.Get(&structProcedure, `select cd_structure, procedure_txt, get_target_product from eln_reaction_note b where b.reaction_id = ?`, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, errors.New("get rxnInfo error:"+err.Error()))
		return
	}

	var analysisData []string
	err = dao.OBCursor.Select(&analysisData, `select file_name from eln_rxn_report_lcms where reaction_id = ? and file_name != '' 
                                          union select file_name from eln_rxn_report_nmr where reaction_id = ? and file_name != '' 
                                        union select file_name from eln_rxn_report_other where reaction_id = ? and file_name != ''`, rxn.ReactionId, rxn.ReactionId, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, errors.New("get analysisData error:"+err.Error()))
		return
	}

	// SignatureDetails 结构体定义
	type SignatureDetails struct {
		BasicInfo    common.ElnRxnBasicInfo `json:"basicInfo"`
		RxnInfo      rxnInfo                `json:"rxnInfo"`
		ReagentInfo  []common.ReagentInfo   `json:"reagentInfo"`
		SolventInfo  []common.ReagentInfo   `json:"solventInfo"`
		ProductInfo  []common.Samples       `json:"productInfo"`
		AnalysisData []string               `json:"analysisData,omitempty"`
	}

	var signD = SignatureDetails{basicInfo, structProcedure, reagents, solvents, sps, analysisData}
	_, err = dao.OBCursor.Exec(`update eln_project_page set export_cnt = export_cnt + 1 where reaction_id = ?`, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", signD)
	return
}
