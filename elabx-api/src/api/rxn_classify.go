// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/11/19 16:02
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : rxn_classify.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"encoding/json"
	"strconv"
)

type ElnRxnClassify struct {
	Rid                    string `json:"rid,omitempty" db:"rid"`
	InputSmiles            string `json:"input_smiles,omitempty" db:"input_smiles"`
	FormattedSmiles        string `json:"formatted_smiles,omitempty" db:"formatted_smiles"`
	MappedSmiles           string `json:"mapped_smiles,omitempty" db:"mapped_smiles"`
	CleanMappedSmiles      string `json:"clean_mapped_smiles,omitempty" db:"clean_mapped_smiles"`
	CanonicalSmiles        string `json:"canonical_smiles,omitempty" db:"canonical_smiles"`
	RetroTemplate0         string `json:"retro_template_0,omitempty" db:"retro_template_0"`
	RetroTemplate1         string `json:"retro_template_1,omitempty" db:"retro_template_1"`
	RetroTemplate2         string `json:"retro_template_2,omitempty" db:"retro_template_2"`
	CleanReagents          string `json:"clean_reagents,omitempty" db:"clean_reagents"`
	MappedAtomsDifferError string `json:"mapped_atoms_differ_error,omitempty" db:"mapped_atoms_differ_error"`
	TemplateFailed         bool   `json:"template_failed,omitempty" db:"template_failed"`
	TemplateFailedReason   string `json:"template_failed_reason,omitempty" db:"template_failed_reason"`
	ReactionType           any    `json:"reaction_type,omitempty" db:"reaction_type"`
	ReactionTypeComments   string `json:"reaction_type_comments,omitempty" db:"reaction_type_comments"`
	ErrorComments          string `json:"error_comments,omitempty" db:"error_comments"`
}

func GetClassifyRes(reactionId int64, cdStructure string, inputFormat string) {
	rxn := struct {
		ReactionId  string `json:"rid" db:"reaction_id"`
		CdStructure string `json:"data" db:"cd_structure"`
		InputFormat string `json:"input_format" db:"input_format" default:"mrv"`
		UserId      int    `json:"user_id" db:"user_id" default:"0"`
	}{
		ReactionId:  strconv.FormatInt(reactionId, 10),
		CdStructure: cdStructure,
		InputFormat: inputFormat,
	}
	jsData, _ := json.Marshal(rxn)
	dataByApi, err := utils.SendDataByApi(jsData, "http://192.168.2.139:6020/api/rxn/classify")
	if err != nil {
		_, _ = dao.OBCursor.Exec(`insert into eln_rxn_classify(rid, rxn_smiles, error_comments) value (?, ?, 'python api error') on duplicate key update rxn_smiles = ?, error_comments = 'python api error' `, reactionId, cdStructure, cdStructure)
	}
	var c ElnRxnClassify
	err = json.Unmarshal(dataByApi, &c)
	if err != nil {
		_, _ = dao.OBCursor.Exec(`insert into eln_rxn_classify(rid, rxn_smiles, error_comments) value (?, ?, 'Unmarshal error') on duplicate key update rxn_smiles = ?, error_comments = 'Unmarshal error' `, reactionId, cdStructure, cdStructure)
	}
	_, err = dao.OBCursor.NamedExec(`insert into eln_rxn_classify(rid, rxn_smiles, std_smiles, cano_smiles, cano_smiles_mapped, retro_template_0, clean_reagents, changed_atom_error, template_failed, template_failed_reason, reaction_type, template_failed_reason, retro_template_1, retro_template_2, error_comments) VALUE (
    								:rid, :rxn_smiles, :std_smiles, :cano_smiles, :cano_smiles_mapped, :retro_template, :clean_reagents, :changed_atom_error, :template_failed, :template_failed_reason, :reaction_type, :classification_comments, :retro_template_1, :retro_template_2, :error_comments) on duplicate key update 
                                    rxn_smiles=:rxn_smiles, std_smiles=:std_smiles, cano_smiles=:cano_smiles,cano_smiles_mapped=:cano_smiles_mapped,
                                    retro_template_0=:retro_template_0, retro_template_1=:retro_template_1,retro_template_2=:retro_template_2,clean_reagents=:clean_reagents,changed_atom_error=changed_atom_error,
                                    template_failed=:template_failed,template_failed_reason=:template_failed_reason,reaction_type=:reaction_type,
                                    reaction_type_comments=:reaction_type_comments, error_comments=:error_comments`, c)
	if err != nil {
		_, _ = dao.OBCursor.Exec(`insert into eln_rxn_classify(rid, rxn_smiles, template_failed_reason) value (?, ?, 'insert to table error') on duplicate key update rxn_smiles = ?, error_comments = 'insert to table error' `, reactionId, cdStructure, cdStructure)
	}
	//var typeName string
	//_ = dao.OBCursor.Get(&typeName, `select classification_comments from eln_rxn_classify where rid = ?`, c.Rid)
	_, err = dao.OBCursor.Exec(`update eln_rxn_basicinfo  set rxn_type_from_ai = ?, rxn_type_code = ? where reaction_id = ?`, c.ReactionType, c.ReactionTypeComments, reactionId)
	if err != nil {
		return
	}
}
