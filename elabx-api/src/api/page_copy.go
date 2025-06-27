// Package api coding=utf-8
// @Project : server
// @Time    : 2025/4/8 14:17
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : page_copy.go
// @Software: GoLand
package api

import (
	"eLabX/src/common"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CopyPage(c *gin.Context) {
	type rxn struct {
		ReactionId int64 `json:"reactionId" db:"reaction_id"`
	}
	var oldReactionId rxn
	err := c.ShouldBind(&oldReactionId)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	// 开始事务
	tx, err := dao.OBCursor.Beginx()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 确保在函数结束时正确地提交或回滚事务
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var rxnBasic common.ElnRxnBasicInfo
	err = tx.Get(&rxnBasic, `select project_name, author_id, author_name, witness_id, witness_name, rxn_type, reference, doi, rxn_type_from_ai from eln_rxn_basicinfo where reaction_id = ?`, oldReactionId.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	newRxnId := utils.GenerateSnowflakeID()
	rxnBasic.ReactionId = newRxnId

	// 插入eln_rxn_basicinfo表
	if err = InsertElnRxnBasicInfo(tx, rxnBasic); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 查询或插入eln_project表，并获取projectId和rankId
	projectId, rankId, err := GetOrCreateProjectIdAndRankId(tx, rxnBasic.ProjectName, rxnBasic.AuthorId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 生成pageName
	pageName := fmt.Sprintf("%s-%03d", rxnBasic.ProjectName, rankId+1)

	// 插入eln_project_page表
	if err = InsertElnProjectPage(tx, newRxnId, rxnBasic.AuthorId, projectId, pageName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 更新eln_project表的rank_id
	if err = UpdateProjectRankId(tx, rxnBasic.ProjectName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 插入eln_commit_log表
	if err = InsertElnCommitLog(tx, newRxnId, rxnBasic.ProjectName, rxnBasic.AuthorId, rxnBasic.AuthorName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if err = copyRxn(tx, oldReactionId.ReactionId, newRxnId); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if err = copyReagents(tx, oldReactionId.ReactionId, newRxnId, pageName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"pageName": pageName, "reactionId": newRxnId})
	return
}

func copyRxn(tx *sqlx.Tx, oldRxnId int64, newRxnId int64) error {
	var rxn common.ElnReactionNote
	err := tx.Get(&rxn, `select reaction_smiles, daylight_smiles, cd_structure, temperature, pressure, time, 
       procedure_txt, procedure_html from eln_reaction_note where reaction_id = ?`, oldRxnId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`insert into eln_reaction_note(reaction_id, 
                                       reaction_smiles, daylight_smiles, cd_structure, temperature, pressure, time, 
                              procedure_txt, procedure_html) value (?,?,?,?,?,?,?,?,?)`, newRxnId, rxn.ReactionSmiles, rxn.DaylightSmiles, rxn.CdStructure, rxn.Temperature, rxn.Pressure, rxn.Time,
		rxn.ProcedureTxt, rxn.ProcedureHtml)
	return err
}

func copyReagents(tx *sqlx.Tx, oldRxnId int64, newRxnId int64, pageName string) error {
	var reagents []common.ElnRxnReagents
	err := tx.Select(&reagents, `select reagent_name, reagent_smiles, mw, exactmass, monoisotopic_mass, formula,
       reagent_role, equiv, cas, concentration, cd_structure, inchi, inchikey, cxsmiles, density, quantity,
       quantity_unit, purity, compound_id, is_chiral, reagent_hash, is_limiting, stereo_centers_cnt,
       product_alias, chiral_descriptor, moles, moles_unit, volume, volume_unit from eln_rxn_reagents where reaction_id = ?`, oldRxnId)
	if err != nil {
		return err
	}

	query := `INSERT INTO eln_rxn_reagents (reagent_id, reaction_id, reagent_name, reagent_smiles,
                              mw, exactmass, monoisotopic_mass, formula, reagent_role, equiv, cas, concentration, 
                              cd_structure, inchi, inchikey, cxsmiles, density, quantity, quantity_unit, purity, 
                              compound_id, is_chiral, reagent_hash, is_limiting, stereo_centers_cnt, product_alias, 
                              chiral_descriptor, moles, moles_unit, volume, volume_unit
        ) VALUES (:reagent_id, :reaction_id,
            :reagent_name, :reagent_smiles, :mw, :exactmass, :monoisotopic_mass, :formula,
            :reagent_role, :equiv, :cas, :concentration, :cd_structure, :inchi, :inchikey, 
            :cxsmiles, :density, :quantity, :quantity_unit, :purity, :compound_id, 
            :is_chiral, :reagent_hash, :is_limiting, :stereo_centers_cnt, :product_alias, 
            :chiral_descriptor, :moles, :moles_unit, :volume, :volume_unit
        )`

	for _, reagent := range reagents {
		reagent.ReagentID = utils.GenerateSnowflakeID()
		reagent.ReactionId = newRxnId
		_, err := tx.NamedExec(query, reagent)
		if err != nil {
			return err
		}

		if strings.ToLower(reagent.ReagentRole) == "product" {
			_, err = tx.Exec(`insert into eln_rxn_samples(reaction_id, reagent_id, reagent_name, mf, mass,
                            sample_id, product_alias) value (?, ?, ?, ?, ?, ?, ?)`, newRxnId, reagent.ReagentID,
				reagent.ReagentName, reagent.Formula, reagent.Mw, pageName+"-"+reagent.ProductAlias, reagent.ProductAlias)
			if err != nil {
				return err
			}
		}
	}

	return err
}
