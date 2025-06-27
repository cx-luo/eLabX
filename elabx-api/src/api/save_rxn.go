// Package api coding=utf-8
// @Project : server
// @Time    : 2025/4/8 16:36
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : save_rxn.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

type SmilesInfo struct {
	DeleteOperationScope int    `json:"deleteOperationScope"`
	DaylightSmiles       string `json:"daylightSmiles" db:"daylight_smiles"`
	RxnCxsmiles          string `json:"cxsmiles" db:"reaction_smiles"`
	ReactionId           int64  `json:"reactionId" db:"reaction_id"`
	CdStructure          string `json:"cdStructure" db:"cd_structure"`
}

func SaveReactionNote(c *gin.Context) {
	var s SmilesInfo
	if err := c.ShouldBind(&s); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	reactionId := s.ReactionId
	if reactionId == 0 {
		reactionId = utils.GenerateSnowflakeID()
	}

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

	// 如果标记为删除，则先备份再删除旧反应信息
	if s.DeleteOperationScope == 1 {
		if err := handleDeletion(tx, s.ReactionId); err != nil {
			handleError(c, err, "handleDeletion", reactionId)
			return
		}
	}

	// 获取反应物信息
	rxnInfo, err := getRxnCompoundInfo(s.CdStructure)
	if err != nil {
		handleError(c, err, "GetRxnCompound", reactionId)
		return
	}

	// 处理反应物
	if err := processReactants(tx, rxnInfo.Reactants, reactionId); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 处理产物
	if err := processProducts(tx, rxnInfo.Product, reactionId); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 插入反应笔记
	if err := insertReactionNote(tx, reactionId, s); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	go func() {
		GetClassifyRes(reactionId, s.CdStructure, "mrv")
	}()

	utils.Success(c, "SaveReactionNote")
	return
}

func handleDeletion(tx *sqlx.Tx, reactionId int64) error {
	// don't delete all reagents, only delete reactant and product
	if utils.GlobalConfig.Service.Backup {
		if err := BackupReactionBeforeDelete(tx, reactionId); err != nil {
			return err
		}
	}

	_, err := tx.Exec(`delete from eln_rxn_reagents where reaction_id = ? and reagent_role in ('product', 'reactant')`, reactionId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from eln_rxn_samples where reaction_id = ? `, reactionId)
	if err != nil {
		return err
	}

	return nil
}

func getRxnCompoundInfo(mrv string) (RxnCompoundData, error) {
	r := CmpdInfo{
		Structure:   mrv,
		InputFormat: "mrv",
		Parameters:  "mrv",
	}
	rxnInfo, err := GetRxnCompound(r)
	return rxnInfo, err
}

func processReactants(tx *sqlx.Tx, reactants []ReagentData, reactionId int64) error {
	for _, reagent := range reactants {
		stereo, stereoCnt, err := GetStereo(reagent.Structure)
		if err != nil {
			return err
		}

		if reagent.Prename == "" {
			outStr, errStr := getReagentNameByShell(reagent.Structure, "mrv", "name:i")
			if errStr != nil {
				utils.Logger.Error(errStr.Error())
			} else {
				reagent.Prename = outStr
			}
		}

		_, err = tx.Exec(`insert into
				eln_rxn_reagents(reagent_id,
				reaction_id,
				reagent_name,
				reagent_smiles,
				reagent_role,
				cd_structure,
				mw, exactmass,
				formula, inchi, inchikey, cxsmiles, is_chiral, stereo_centers_cnt, chiral_descriptor
				) value (?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update is_chiral = ?`, utils.GenerateSnowflakeID(),
			reactionId, reagent.Prename, reagent.GetSmiles(), "reactant", reagent.Structure, reagent.Weight, reagent.EMass,
			reagent.Formula, reagent.GetInchi(), reagent.GetInchiKey(), reagent.GetCxSmiles(), stereoCnt > 0, stereoCnt, stereo, stereoCnt > 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func processProducts(tx *sqlx.Tx, products []ReagentData, reactionId int64) error {
	for i, reagent := range products {
		productId := utils.GenerateSnowflakeID()
		stereo, stereoCnt, err := GetStereo(reagent.Structure)
		if err != nil {
			return err
		}

		if reagent.Prename == "" {
			outStr, errStr := getReagentNameByShell(reagent.Structure, "mrv", "name:i")
			if errStr != nil {
				utils.Logger.Error(errStr.Error())
			} else {
				reagent.Prename = outStr
			}
		}

		alias := "P" + strconv.Itoa(i+1)
		_, err = tx.Exec(`insert into
				eln_rxn_reagents(reagent_id,
				reaction_id,
				equiv,
				reagent_name,
				reagent_smiles,
				reagent_role,
				cd_structure,
				mw, exactmass,
				formula, inchi, inchikey, cxsmiles, is_chiral, stereo_centers_cnt, chiral_descriptor, product_alias
				) value (?, ?, 1, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?) on duplicate key update is_chiral = ?`, productId,
			reactionId, reagent.Prename, reagent.GetSmiles(), "product", reagent.Structure, reagent.Weight, reagent.EMass,
			reagent.Formula, reagent.GetInchi(), reagent.GetInchiKey(), reagent.GetCxSmiles(),
			stereoCnt > 0, stereoCnt, stereo, alias, stereoCnt > 0)
		if err != nil {
			return err
		}

		pageName, err := getPageName(reactionId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`insert into eln_rxn_samples(reaction_id, reagent_id, reagent_name, mf, mass, sample_id, product_alias) value (?, ?, ?, ?, ?, ?, ?)`, reactionId, productId, reagent.Prename, reagent.Formula, reagent.Weight, pageName+"-"+alias, alias)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertReactionNote(tx *sqlx.Tx, reactionId int64, s SmilesInfo) error {
	_, err := tx.Exec(`insert into eln_reaction_note(reaction_id, 
                                       reaction_smiles, daylight_smiles, cd_structure, pressure, pressure_unit) value (?, ?, ?, ?, 1, 'atm') 
                                       on duplicate key update reaction_smiles = ?, daylight_smiles = ?, cd_structure = ? `,
		reactionId, s.RxnCxsmiles, s.DaylightSmiles, s.CdStructure,
		s.RxnCxsmiles, s.DaylightSmiles, s.CdStructure)
	return err
}

func handleError(c *gin.Context, err error, funcName string, reactionId int64) {
	c.JSON(http.StatusInternalServerError, utils.BaseResponse{
		StatusCode: 500, Msg: err.Error(), Data: gin.H{"func": funcName, "ReactionId": reactionId},
	})
}
