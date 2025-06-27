// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/4/20 15:22
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : eln.go
// @Software: GoLand
package api

import (
	"database/sql"
	"eLabX/src/common"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func GetBasicInfo(c *gin.Context) {
	var rxn struct {
		ReactionId int64 `json:"reactionId" db:"reaction_id"`
	}
	err := c.ShouldBind(&rxn)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var rxnBasicInfo common.ElnRxnBasicInfo
	err = dao.OBCursor.Get(&rxnBasicInfo, `select a.reaction_id, a.project_name, a.batch,
	a.step_id, a.author_id, a.author_name, a.witness_id, a.witness_name, a.rxn_type,
	a.reference, a.doi,a.creation_date, a.start_date, rxn_status, a.comment, a.rxn_type_from_ai, a.rxn_type_code, b.page_name
	from eln_rxn_basicinfo a join eln_project_page b on
	a.reaction_id = b.reaction_id and a.reaction_id = ?`, rxn.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"rxnBasicInfo": rxnBasicInfo})
	return
}

func SaveBasicInfo(c *gin.Context) {
	var rxn common.ElnRxnBasicInfo
	err := c.ShouldBind(&rxn)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.NamedExec(`update eln_rxn_basicinfo set batch=:batch, 
                             step_id=:step_id, witness_id=:witness_id, witness_name=:witness_name, rxn_type=:rxn_type, 
                             reference=:reference, doi=:doi, start_date=:start_date, comment=:comment 
                             where reaction_id=:reaction_id`, rxn)

	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func BackupReactionBeforeDelete(tx *sqlx.Tx, rxnId int64) error {
	_, err := tx.Exec(`replace into elabx_deleted.eln_reaction_note_bak select id, reaction_id, reaction_smiles, daylight_smiles, cd_structure, temperature, rxn_molarity, pressure, time, get_target_product, comments, procedure_txt, procedure_html, gmt_create, gmt_modified from eln_reaction_note where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_reaction_note_bak," + err.Error())
	}
	_, err = tx.Exec(`replace into elabx_deleted.eln_rxn_basicinfo_bak select reaction_id, project_name, batch, step_id, author_id, author_name, witness_id, witness_name, rxn_type, reference, doi, creation_date, start_date, rxn_status, gmt_create, gmt_modified, comment, rxn_type_from_ai, commit_date from eln_rxn_basicinfo where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_rxn_basicinfo_bak," + err.Error())
	}
	_, err = tx.Exec(`replace into elabx_deleted.eln_rxn_reagents_bak select reagent_id, reaction_id, reagent_name, reagent_smiles, mw, exactmass, monoisotopic_mass, formula, reagent_role, equiv, cas, concentration, cd_structure, inchi, inchikey, cxsmiles, density, quantity, quantity_unit, purity, compound_id, is_chiral, reagent_hash, is_limiting, stereo_centers_cnt, product_alias, chiral_descriptor, moles, moles_unit, volume, volume_unit from eln_rxn_reagents where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_rxn_reagents_bak," + err.Error())
	}
	_, err = tx.Exec(`replace into elabx_deleted.eln_rxn_samples_bak select reaction_id, reagent_id, reagent_name, sample_status, is_test_samples, sample_type, sample_id, mf, mass, purity, synthesised, synthesised_unit, sample_yield, sample_reference, moles, moles_unit, use_in_yield, amount_submitted, amount_submitted_unit, barcode, enantiomeric_purity, product_name, qa_reason, product_alias, color, qualifier from eln_rxn_samples where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_rxn_samples_bak," + err.Error())
	}
	_, err = tx.Exec(`replace into elabx_deleted.eln_rxn_report_lcms_bak select report_rank, reaction_id, reagent_id, technique, m_type, m_value, calculate_mw, find_mw, purity, rt, result, solvent, has_spectrum, spectrum_url, product_alias, file_name, file_md5, created_at, updated_at from eln_rxn_report_lcms where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_rxn_report_lcms_bak," + err.Error())
	}
	_, err = tx.Exec(`replace into elabx_deleted.eln_rxn_report_nmr_bak select report_rank, reaction_id, reagent_id, technique, rt, result, solvent, spectrum_url, product_alias, file_name, file_md5, created_at, updated_at from eln_rxn_report_nmr where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_rxn_report_nmr_bak," + err.Error())
	}
	_, err = tx.Exec(`replace into elabx_deleted.eln_rxn_report_other_bak select report_rank, reaction_id, reagent_id, technique, rt, comments, spectrum_url, created_at, updated_at, file_name, file_md5, product_alias from eln_rxn_report_other where reaction_id = ?`, rxnId)
	if err != nil {
		return errors.New("Backup table error: table - eln_rxn_report_other_bak," + err.Error())
	}
	return nil
}

func DeleteReaction(tx *sqlx.Tx, reactionId int64) error {
	_, err := tx.Exec(`delete from eln_reaction_note where reaction_id = ?`, reactionId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`delete from eln_rxn_reagents where reaction_id = ?`, reactionId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from eln_rxn_samples where reaction_id = ?`, reactionId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from eln_rxn_report_lcms where reaction_id = ?`, reactionId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from eln_rxn_report_nmr where reaction_id = ?`, reactionId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from eln_rxn_report_other where reaction_id = ?`, reactionId)
	if err != nil {
		return err
	}

	return nil
}

func getPageName(reactionId int64) (string, error) {
	var s string
	err := dao.OBCursor.Get(&s, `select page_name from eln_project_page where
                                         reaction_id = ?`, reactionId)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(s, "-", 2)

	// 检查是否有第二个部分
	if len(parts) > 1 {
		afterFirstDash := strings.Join(parts[1:], "-")
		return afterFirstDash, nil
	}
	return "", nil
}

func getReagentNameByShell(structure string, inputFmt string, outputFmt string) (string, error) {
	var toname = StructForName{inputFmt, outputFmt, structure}
	iupacName, err := GetIupacName(toname)

	if err != nil {
		cmd := exec.Command("molconvert", "name:i", "-f mrv", "-s", fmt.Sprintf("'%s'", structure))

		outStr, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(outStr), nil
	}

	return iupacName, nil
}

// SaveReactionNote 比较复杂，放到了单独的文件中

func SaveReactionConditions(c *gin.Context) {
	var rc struct {
		ReactionID       int64   `json:"reactionId"`
		ReactionMolarity int     `json:"reactionMolarity"`
		Temperature      float64 `json:"temperature"`
		Pressure         string  `json:"pressure"`
		PressureUnit     string  `json:"pressureUnit,omitempty"`
		Time             float64 `json:"time"`
	}
	err := c.ShouldBind(&rc)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	_, err = dao.OBCursor.Exec(`update eln_reaction_note set rxn_molarity = ?, 
                                       temperature = ?, pressure = ?, pressure_unit = ?, time = ? where reaction_id = ?`, rc.ReactionMolarity,
		rc.Temperature, rc.Pressure, rc.PressureUnit, rc.Time, rc.ReactionID)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "", Data: gin.H{"reactionId": rc.ReactionID},
	})

	return
}

func GetSamples(c *gin.Context) {
	var rid struct {
		ReactionId int64 `json:"reactionId"`
	}
	err := c.ShouldBind(&rid)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	var s []common.Samples
	err = dao.OBCursor.Select(&s, `select reaction_id, reagent_id, sample_status, is_test_samples, color,
       sample_type, sample_id, mf, mass, purity, synthesised, synthesised_unit, sample_yield, sample_reference, 
       moles, moles_unit, use_in_yield, amount_submitted,amount_submitted_unit, barcode, enantiomeric_purity, product_name, product_alias, 
       qa_reason, qualifier from eln_rxn_samples where reaction_id = ?`, rid.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "", Data: gin.H{"samples": s},
	})
	return
}

func CommitRxn(c *gin.Context) {
	var s struct {
		ReactionId uint64 `json:"reactionId" db:"reaction_id"`
		UserId     uint   `json:"userId" db:"user_id"`
		UserName   string `json:"userName" db:"user_name"`
		Status     string `json:"status" db:"status"`
		Project    string `json:"project" db:"project"`
	}
	err := c.ShouldBind(&s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
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

	var rxnCode string
	err = tx.Get(&rxnCode, `select rxn_type_code from eln_rxn_basicinfo where reaction_id = ? `, s.ReactionId)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	if rxnCode == "-1" || rxnCode == "-3" || rxnCode == "-6" || rxnCode == "-7" {
		utils.BadRequestErr(c, errors.New("there is an error in the reaction, please check comments and correct it"))
		return
	}
	_, err = tx.Exec(`update eln_rxn_basicinfo set rxn_status = ?,commit_date = ? where reaction_id = ?`, s.Status, time.Now().Format("2006-01-02"), s.ReactionId)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = tx.Exec(`insert into eln_commit_log(reaction_id, project, user_id, user_name, status ) VALUE (?,?,?,?,?)`, s.ReactionId,
		s.Project, s.UserId, s.UserName, s.Status)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

func SaveSamples(c *gin.Context) {
	var s common.Samples
	err := c.ShouldBind(&s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.NamedExec(`update eln_rxn_samples set
                                    sample_status=:sample_status,
                                    is_test_samples=:is_test_samples, sample_type=:sample_type, sample_id=:sample_id,
                                    mf=:mf, purity=:purity, mass=:mass, synthesised=:synthesised, synthesised_unit=:synthesised_unit,
                                    sample_yield=:sample_yield, sample_reference=:sample_reference, moles=:moles, moles_unit=:moles_unit, use_in_yield=:use_in_yield,
                                    amount_submitted=:amount_submitted, amount_submitted_unit=:amount_submitted_unit ,barcode=:barcode, enantiomeric_purity=:enantiomeric_purity,
                                    qa_reason=:qa_reason, color=:color, qualifier=:qualifier
                                where reaction_id=:reaction_id and reagent_id=:reagent_id`, s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func GetReagentInfoByRole(c *gin.Context) {
	var r struct {
		//ReagentRole   string `json:"reagentRole"`
		ReagentPrefix string `json:"reagentPrefix"`
	}
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	type chem struct {
		Cid         int64   `json:"cid,omitempty" db:"cid"`
		CmpdName    string  `json:"cmpdname,omitempty" db:"cmpdname"`
		Smiles      string  `json:"smiles,omitempty" db:"smiles"`
		Mw          float64 `json:"mw,omitempty" db:"mw"`
		Formula     string  `json:"mf,omitempty" db:"mf"`
		Cas         string  `json:"cas,omitempty" db:"cas"`
		InChI       string  `json:"inchi,omitempty" db:"InChI"`
		InChIKey    string  `json:"inchikey,omitempty" db:"InChIKey"`
		Iupac       string  `json:"iupacName,omitempty" db:"iupacname"`
		ReagentRole string  `json:"reagentRole,omitempty" db:"reagent_role"`
	}
	var s []chem
	err = dao.OBCursor.Select(&s, `select cid, cmpdname, smiles, mw, mf, cas, InChI, InChIKey, iupacname, reagent_role from eln_std_reagents where
                                            cmpdname like ? or cas like ?`, fmt.Sprintf(`%%%s%%`, r.ReagentPrefix), fmt.Sprintf(`%%%s%%`, r.ReagentPrefix))
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	s = append(s, chem{
		Cid:      0,
		CmpdName: r.ReagentPrefix,
		Smiles:   "",
		Mw:       0,
		Formula:  "",
		Cas:      "",
		InChI:    "",
		InChIKey: "",
		Iupac:    "",
	})
	utils.SuccessWithData(c, "", gin.H{"reagents": s})
	return
}

func DeepCopyByJson(src common.ReagentInfo) (*common.ReagentInfo, error) {
	var dst = new(common.ReagentInfo)
	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, dst)
	return dst, err
}

func LoadWorkbook(c *gin.Context) {
	type rxn struct {
		ReactionId int64 `json:"reactionId" db:"reaction_id"`
	}
	type smilesInfo struct {
		ReactionId       int64   `json:"reactionId,omitempty" db:"reaction_id"`
		ReactionSmiles   string  `json:"reactionSmiles,omitempty" db:"reaction_smiles"`
		CdStructure      string  `json:"cdStructure,omitempty" db:"cd_structure"`
		ReactionMolarity string  `json:"reactionMolarity,omitempty" db:"rxn_molarity"`
		Temperature      float64 `json:"temperature" db:"temperature"`
		Time             float64 `json:"time" db:"time"`
		Pressure         string  `json:"pressure" db:"pressure"`
		PressureUnit     string  `json:"pressureUnit,omitempty" db:"pressure_unit"`
		Procedure        string  `json:"procedure,omitempty" db:"procedure_txt"`
		GetTargetProduct int8    `json:"getTargetProduct,omitempty" db:"get_target_product"` // 是否获取到目标产物
		Comments         string  `json:"comments,omitempty" db:"comments"`
		RxnImg           string  `json:"rxnImg,omitempty"`
	}

	var w rxn
	err := c.ShouldBind(&w)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	var wbInfo smilesInfo
	err = dao.OBCursor.Get(&wbInfo, `SELECT
				a.reaction_id, a.reaction_smiles,a.cd_structure, a.rxn_molarity,
				a.temperature, a.pressure, a.pressure_unit, a.time, a.get_target_product, a.comments
			FROM
				eln_reaction_note a
			WHERE a.reaction_id = ?`, w.ReactionId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundError(c, sql.ErrNoRows)
			return
		}
		utils.BadRequestErr(c, err)
		return
	}

	err, wbInfo.RxnImg = GenImgBase64String(wbInfo.CdStructure)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	lcmsData, err := GetLcmsData(wbInfo.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	nmrData, err := GetNmrData(wbInfo.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	otherReportData, err := GetOtherReport(wbInfo.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"workbooks": wbInfo, "lcms": lcmsData, "nmr": nmrData, "otherReport": otherReportData})
	return
}

func LoadReagents(c *gin.Context) {
	var rid struct {
		ReactionId int64 `json:"reactionId"`
	}
	err := c.ShouldBind(&rid)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var reagent []common.ReagentInfo
	err = dao.OBCursor.Select(&reagent, `SELECT
    			a.reagent_id, a.reagent_role, a.reagent_name , a.reagent_smiles ,
				a.cd_structure, a.density, a.purity, a.equiv,
				a.quantity, quantity_unit, a.moles, a.moles_unit, a.volume, a.volume_unit,
				a.mw, a.concentration, a.chiral_descriptor,
				a.formula , a.cas, a.compound_id, a.is_limiting, a.is_chiral, a.stereo_centers_cnt, a.product_alias
			FROM
				eln_rxn_reagents a
			WHERE a.reaction_id = ? order by reagent_id`, rid.ReactionId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundError(c, sql.ErrNoRows)
			return
		}
		utils.BadRequestErr(c, err)
		return
	}
	// 使用深拷贝拷贝数据，返回结果才正确
	var newReagentsWithImg []common.ReagentInfo
	for _, ra := range reagent {
		var structure string
		if ra.CdStructure != "" {
			structure = ra.CdStructure
		} else {
			structure = ra.ReagentSmiles
		}
		err, ra.ReagentImg = GenImgBase64String(structure)
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
		newRa, err := DeepCopyByJson(ra)
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
		newReagentsWithImg = append(newReagentsWithImg, *newRa)
	}

	utils.SuccessWithData(c, "", gin.H{"reagents": newReagentsWithImg})
	return
}

func SaveNewReagent(c *gin.Context) {
	type ri struct {
		WorkgroupID   int64   `json:"workgroupId,omitempty" db:"workgroup_id"`
		WorkbookID    int64   `json:"workbookId,omitempty" db:"workbook_id"`
		ReactionId    int64   `json:"reactionId,omitempty" db:"reaction_id"`
		ReagentName   string  `json:"reagentName,omitempty" db:"reagent_name"`
		ReagentSmiles string  `json:"reagentSmiles,omitempty" db:"reagent_smiles"`
		CdStructure   string  `json:"cdStructure,omitempty" db:"cd_structure"`
		Mw            float64 `json:"mw,omitempty" db:"mw"`
		ReagentRole   string  `json:"reagentRole,omitempty" db:"reagent_role"`
		Formula       string  `json:"mf,omitempty" db:"formula"`
		Cas           string  `json:"cas,omitempty" db:"cas"`
		Inchi         string  `json:"inchi,omitempty" db:"inchi"`
		Inchikey      string  `json:"inchikey,omitempty" db:"inchikey"`
		CompoundId    int64   `json:"compoundId,omitempty" db:"compound_id"`
	}
	var r ri
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	if r.Formula == "" || r.Mw == float64(0) {
		utils.BadRequestErr(c, errors.New("formula or mw can not be null, please click the search button"))
		return
	}

	reagentId := utils.GenerateSnowflakeID()
	insertSql := `insert into eln_rxn_reagents(reagent_id, reaction_id,
                                      reagent_name, reagent_smiles, mw, cxsmiles, formula, reagent_role, cas, inchi, 
                                      inchikey, compound_id) VALUE (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = dao.OBCursor.Exec(insertSql, reagentId, r.ReactionId,
		r.ReagentName, r.ReagentSmiles, r.Mw, r.ReagentSmiles, r.Formula, r.ReagentRole, r.Cas, r.Inchi, r.Inchikey, r.CompoundId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	var imgStr, structure string
	if r.CdStructure != "" {
		structure = r.CdStructure
	} else {
		structure = r.ReagentSmiles
	}
	err, imgStr = GenImgBase64String(structure)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "", Data: gin.H{"reagentId": reagentId, "reagentImg": imgStr},
	})
	return
}

func DeleteReagent(c *gin.Context) {
	var rid = struct {
		ReagentId int64 `json:"reagentId,omitempty" db:"reagent_id"`
	}{}
	err := c.ShouldBind(&rid)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`delete from eln_rxn_reagents where reagent_id = ?`, rid.ReagentId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

func SaveAdditionalInfo(c *gin.Context) {
	var r common.ReagentInfo
	err := c.ShouldBind(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BaseResponse{
			StatusCode: 400, Msg: err.Error(), Data: gin.H{},
		})
		return
	}
	sqlStr := `update eln_rxn_reagents set cas=:cas, purity=:purity, quantity=:quantity, quantity_unit=:quantity_unit,
                                     concentration=:concentration, density=:density, equiv=:equiv,
                                     compound_id=:compound_id, is_limiting=:is_limiting, moles=:moles, reagent_role=:reagent_role,
                                     moles_unit=:moles_unit, chiral_descriptor=:chiral_descriptor, volume=:volume, volume_unit=:volume_unit
                              where reagent_id=:reagent_id`
	_, err = dao.OBCursor.NamedExec(sqlStr, r)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BaseResponse{
			StatusCode: 400, Msg: err.Error(), Data: gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		StatusCode: 200, Msg: "", Data: gin.H{},
	})

	return
}

func SetReagentCid(c *gin.Context) {
	var r struct {
		ReagentId  int64 `json:"reagentId,omitempty" db:"reagentId"`
		CompoundID int64 `json:"compoundId,omitempty" db:"compoundID"`
	}
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.NamedExec(`update eln_rxn_reagents set compound_id = :compoundID                                    
                                      where reagent_id = :reagentId`, r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func SaveNewStdReagent(c *gin.Context) {
	type chem struct {
		Cid         int64   `json:"cid,omitempty" db:"cid"`
		CmpdName    string  `json:"cmpdname,omitempty" db:"cmpdname"`
		Smiles      string  `json:"smiles,omitempty" db:"smiles"`
		Mw          float64 `json:"mw,omitempty" db:"mw"`
		Formula     string  `json:"formula,omitempty" db:"mf"`
		Cas         string  `json:"cas,omitempty" db:"cas"`
		InChI       string  `json:"inchi,omitempty" db:"inchi"`
		InChIKey    string  `json:"inchikey,omitempty" db:"inchikey"`
		Iupac       string  `json:"iupac,omitempty" db:"iupacname"`
		Role        string  `json:"reagentRole,omitempty" db:"reagent_role"`
		UserId      int     `json:"userId,omitempty" db:"user_id"`
		CdStructure string  `json:"cdStructure,omitempty" db:"cd_structure"`
	}
	var n chem
	err := c.ShouldBind(&n)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	//if n.Iupac != "" {
	//	_, err = dao.OBCursor.NamedExec(`insert into eln_std_reagents(cmpdname, smiles, cas, InChI, InChIKey, mw, mf, iupacname, cid, reagent_role, cd_structure)
	//VALUE (:cmpdname,:smiles,:cas,:inchi,:inchikey,:mw,:mf,:iupacname,:cid,:reagent_role, :cd_structure)`, n)
	//	if err != nil {
	//		utils.BadRequestErr(c, err)
	//		return
	//	}
	//} else {
	_, err = dao.OBCursor.NamedExec(`insert into eln_std_to_review (cmpdname, smiles, cas, InChI, InChIKey, mw, mf, iupacname, cid, reagent_role, user_id, cd_structure) 
    VALUE (:cmpdname,:smiles,:cas,:inchi,:inchikey,:mw,:mf,:iupacname,:cid,:reagent_role, :user_id, :cd_structure)`, n)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	//}

	utils.Success(c, "")
	return
}

type ElnRxnReportLcms struct {
	ReportRank   int64   `json:"reportRank,omitempty" db:"report_rank"`
	ReactionId   int64   `json:"reactionId,omitempty" db:"reaction_id"`
	ReagentId    int64   `json:"reagentId,omitempty" db:"reagent_id"`
	Technique    string  `json:"technique,omitempty" db:"technique"`
	MType        string  `json:"mType,omitempty" db:"m_type"`
	MValue       string  `json:"mValue,omitempty" db:"m_value"`
	CalculateMw  float32 `json:"calculateMw,omitempty" db:"calculate_mw"`
	FindMw       float32 `json:"findMw,omitempty" db:"find_mw"`
	Purity       float32 `json:"purity,omitempty" db:"purity"`
	Rt           float32 `json:"rt,omitempty" db:"rt"`
	Result       string  `json:"result,omitempty" db:"result"`
	Solvent      string  `json:"solvent,omitempty" db:"solvent"`
	HasSpectrum  int8    `json:"hasSpectrum,omitempty" db:"has_spectrum"`
	SpectrumUrl  string  `json:"spectrumUrl,omitempty" db:"spectrum_url"`
	ProductAlias string  `json:"productAlias,omitempty" db:"product_alias"`
	FileName     string  `json:"fileName,omitempty" db:"file_name"`
	FileMd5      string  `json:"fileMd5,omitempty" db:"file_md5"`
}

func GetLcmsData(reactionId int64) ([]ElnRxnReportLcms, error) {
	var l []ElnRxnReportLcms
	err := dao.OBCursor.Select(&l, `select report_rank, reaction_id, reagent_id, technique, m_type, m_value, 
       calculate_mw, find_mw, purity, rt, result, solvent, has_spectrum, spectrum_url, product_alias, file_name, file_md5 from 
                                                                                           eln_rxn_report_lcms 
                                                                                       where reaction_id = ? order by reagent_id`, reactionId)
	if err != nil {
		return []ElnRxnReportLcms{}, err
	}
	return l, nil
}

func SaveLcms(c *gin.Context) {
	var l ElnRxnReportLcms
	err := c.ShouldBind(&l)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	autoId, err := dao.OBCursor.NamedExec(`insert into eln_rxn_report_lcms(reaction_id, reagent_id, report_rank, technique,
                               m_type, m_value, calculate_mw, find_mw, purity, rt, result, solvent, has_spectrum, 
                               spectrum_url, product_alias,file_name,file_md5) value (:reaction_id, :reagent_id, :report_rank, :technique,
                               :m_type, :m_value, :calculate_mw, :find_mw, :purity, :rt, :result, :solvent, :has_spectrum, 
                               :spectrum_url, :product_alias, :file_name, :file_md5)  on duplicate key update m_type=:m_type, m_value=:m_value, calculate_mw=:calculate_mw,
                            	find_mw=:find_mw, purity=:purity, rt=:rt, result=:result, solvent=:solvent, spectrum_url=:spectrum_url, file_name=:file_name, file_md5=:file_md5`, l)

	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	v, err := autoId.LastInsertId()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"reportRank": v})
	return
}

type ElnRxnReportNmr struct {
	ReportRank   int64   `json:"reportRank" db:"report_rank"`
	ReactionId   int64   `json:"reactionId,omitempty" db:"reaction_id"`
	ReagentId    int64   `json:"reagentId,omitempty" db:"reagent_id"`
	Technique    string  `json:"technique,omitempty" db:"technique"`
	Rt           float32 `json:"rt,omitempty" db:"rt"`
	Result       string  `json:"result,omitempty" db:"result"`
	Solvent      string  `json:"solvent,omitempty" db:"solvent"`
	HasSpectrum  int8    `json:"hasSpectrum,omitempty" db:"has_spectrum"`
	SpectrumUrl  string  `json:"spectrumUrl,omitempty" db:"spectrum_url"`
	ProductAlias string  `json:"productAlias,omitempty" db:"product_alias"`
	FileName     string  `json:"fileName,omitempty" db:"file_name"`
	FileMd5      string  `json:"fileMd5,omitempty" db:"file_md5"`
}

func GetNmrData(reactionId int64) ([]ElnRxnReportNmr, error) {
	var l []ElnRxnReportNmr
	err := dao.OBCursor.Select(&l, `select report_rank, reaction_id, reagent_id, technique, rt, result, solvent,
       spectrum_url, product_alias, file_name, file_md5 from eln_rxn_report_nmr where reaction_id = ? order by reagent_id`, reactionId)
	if err != nil {
		return []ElnRxnReportNmr{}, err
	}
	return l, nil
}

func SaveNmr(c *gin.Context) {
	var l ElnRxnReportNmr
	err := c.ShouldBind(&l)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	autoId, err := dao.OBCursor.NamedExec(`insert into eln_rxn_report_nmr(reaction_id, reagent_id, report_rank, technique, rt, 
                                     result, solvent, spectrum_url, product_alias, file_name, file_md5) VALUE (:reaction_id, :reagent_id, :report_rank, :technique, :rt, 
                                     :result, :solvent, :spectrum_url, :product_alias, :file_name, :file_name) on duplicate 
                                         key update technique=:technique, result=:result, solvent=:solvent, rt=:rt, spectrum_url=:spectrum_url, file_name=:file_name, file_md5=:file_md5`, l)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	v, err := autoId.LastInsertId()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"reportRank": v})
	return
}

type ElnRxnOtherReports struct {
	ReportRank   int64  `json:"reportRank" db:"report_rank"`
	ReactionId   int64  `json:"reactionId" db:"reaction_id"`
	ReagentId    int64  `json:"reagentId" db:"reagent_id"`
	Technique    string `json:"technique" db:"technique"`
	Rt           string `json:"rt,omitempty," db:"rt"`
	Comments     string `json:"comments,omitempty" db:"comments"`
	SpectrumUrl  string `json:"spectrumUrl,omitempty" db:"spectrum_url"`
	FileName     string `json:"fileName,omitempty" db:"file_name"`
	FileMd5      string `json:"fileMd5,omitempty" db:"file_md5"`
	ProductAlias string `json:"productAlias,omitempty" db:"product_alias"`
}

func GetOtherReport(rid int64) ([]ElnRxnOtherReports, error) {
	var o []ElnRxnOtherReports
	err := dao.OBCursor.Select(&o, `select report_rank, reaction_id, reagent_id, technique, rt, comments, spectrum_url, file_name, file_md5, product_alias from 
                                                                        eln_rxn_report_other where
                                              reaction_id = ?`, rid)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	return o, nil
}

func SaveOtherReport(c *gin.Context) {
	var r ElnRxnOtherReports
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	autoId, err := dao.OBCursor.NamedExec(`insert into eln_rxn_report_other(report_rank, reaction_id, reagent_id,
                                       technique, rt, comments, spectrum_url, file_name, file_md5, product_alias) 
    						VALUE (:report_rank, :reaction_id, :reagent_id, :technique, :rt, :comments, :spectrum_url, :file_name, :file_md5, :product_alias) 
 							on duplicate key update technique =:technique, rt = :rt, comments = :comments, spectrum_url = :spectrum_url, file_name=:file_name, file_md5=:file_md5, product_alias=:product_alias `,
		r)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	v, err := autoId.LastInsertId()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"reportRank": v})
	return
}

func DeleteLcmsReport(c *gin.Context) {
	var r ElnRxnOtherReports
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`delete from eln_rxn_report_lcms where reagent_id = ? and report_rank = ?`,
		r.ReagentId, r.ReportRank)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func DeleteNmrReport(c *gin.Context) {
	var r ElnRxnOtherReports
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`delete from eln_rxn_report_nmr where reagent_id = ? and report_rank = ?`,
		r.ReagentId, r.ReportRank)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func DeleteOtherReport(c *gin.Context) {
	var r ElnRxnOtherReports
	err := c.ShouldBind(&r)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`delete from eln_rxn_report_other where reagent_id = ? and report_rank = ?`,
		r.ReagentId, r.ReportRank)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

type ElnRxnProcedure struct {
	AutoId        int64  `json:"autoId,omitempty" db:"auto_id" gorm:"auto_id"`
	ParentId      int64  `json:"parentId,omitempty" db:"parent_id" gorm:"parent_id"`
	TabName       string `json:"tabName,omitempty" db:"tab_name" gorm:"tab_name"`
	Label         string `json:"label,omitempty" db:"label" gorm:"label"`
	EnTxt         string `json:"enTxt,omitempty" db:"en_txt" gorm:"en_txt"`
	EnHtml        string `json:"enHtml,omitempty" db:"en_html" gorm:"en_html"`
	ZhTxt         string `json:"zhTxt,omitempty" db:"zh_txt" gorm:"zh_txt"`
	UseForCompany string `json:"useForCompany,omitempty" db:"use_for_company" gorm:"use_for_company"` // 例句应用于哪些公司,多个用逗号分隔
	TheOrder      int64  `json:"theOrder,omitempty" db:"the_order" gorm:"the_order"`
	TabType       string `json:"tabType,omitempty" db:"tab_type" gorm:"tab_type"`
}

func GetProcedureTemplate(c *gin.Context) {
	var e []ElnRxnProcedure
	err := dao.OBCursor.Select(&e, `SELECT tab_name ,label ,en_txt ,zh_txt,en_html from eln_std_procedure order by parent_id ,auto_id`)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", e)
	return
}

func GetProcedure(c *gin.Context) {
	var s struct {
		ReactionId int64 `json:"reactionId" db:"reaction_id"`
	}
	err := c.ShouldBind(&s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var procedureInfo struct {
		ProcedureTxt     string `json:"procedureTxt,omitempty" db:"procedure_txt"`
		ProcedureHtml    string `json:"procedureHtml,omitempty" db:"procedure_html"`
		GetTargetProduct int8   `json:"getTargetProduct,omitempty" db:"get_target_product"`
		Comments         string `json:"comments,omitempty" db:"comments"`
	}
	err = dao.OBCursor.Get(&procedureInfo, `select procedure_txt, procedure_html, get_target_product, comments from eln_reaction_note where reaction_id = ?`, s.ReactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFoundError(c, sql.ErrNoRows)
			return
		}
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"procedure": procedureInfo})
	return
}

func SaveProcedure(c *gin.Context) {
	var s struct {
		ReactionId    int64  `json:"reactionId" db:"reaction_id"`
		ProcedureText string `json:"procedureText" db:"procedure_txt"`
		ProcedureHtml string `json:"procedureHtml" db:"procedure_html"`
	}
	err := c.ShouldBind(&s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	_, err = dao.OBCursor.Exec(`update eln_reaction_note set procedure_txt = ?, procedure_html = ? where reaction_id = ?`, s.ProcedureText, s.ProcedureHtml, s.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func SaveProcedureComments(c *gin.Context) {
	var s struct {
		ReactionId       int64  `json:"reactionId" db:"reaction_id"`
		GetTargetProduct int    `json:"getTargetProduct,omitempty" db:"get_target_product"`
		Comments         string `json:"comments,omitempty" db:"comments"`
	}
	err := c.ShouldBind(&s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	_, err = dao.OBCursor.Exec(`update eln_reaction_note set get_target_product = ?, comments = ? where reaction_id = ?`, s.GetTargetProduct, s.Comments, s.ReactionId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}
