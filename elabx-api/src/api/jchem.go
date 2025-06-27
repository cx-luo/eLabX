// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/5/17 8:58
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : jchem.go
// @Software: GoLand
package api

import (
	"eLabX/src/common"
	"eLabX/src/utils"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
	"strings"
)

type ReagentData struct {
	Structure string  `json:"structure,omitempty"`
	Prename   string  `json:"prename,omitempty"`
	Formula   string  `json:"formula,omitempty"`
	Weight    float64 `json:"weight,omitempty"`
	EMass     float64 `json:"eMass,omitempty"`
}

// SingleMolResponseInfo for GetSingleMol
type SingleMolResponseInfo struct {
	Msg  string      `json:"msg,omitempty"`
	Code int         `json:"code,omitempty"`
	Data ReagentData `json:"data,omitempty"`
}

func GetSingleMol(smiles string) (SingleMolResponseInfo, error) {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{ "structure":"%s"}`, smiles))
	req, err := http.NewRequest("POST", common.JChemApi.GetSingleMolApi, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return SingleMolResponseInfo{}, err
	}
	var m SingleMolResponseInfo
	err = json.Unmarshal(bodyText, &m)
	if err != nil {
		return SingleMolResponseInfo{}, err
	}
	return m, nil
}

type MolExportRes struct {
	// BinaryStructure : only generate image
	BinaryStructure string `json:"binaryStructure,omitempty"`
	Format          string `json:"format"`
	Structure       string `json:"structure,omitempty"`
}

type CmpdInfo struct {
	FilterChain []FilterChain `json:"filterChain,omitempty"`
	InputFormat string        `json:"inputFormat"`
	Parameters  string        `json:"parameters"`
	Structure   string        `json:"structure"`
}
type Parameters struct {
	StandardizerDefinition string `json:"standardizerDefinition,omitempty" default:"aromatize"`
}
type FilterChain struct {
	Filter     string     `json:"filter,omitempty" default:"standardizer"`
	Parameters Parameters `json:"parameters,omitempty"`
}

// MolConvert : return without "data:image/png;base64,"
func MolConvert(s CmpdInfo) (MolExportRes, error) {
	jsData, _ := json.Marshal(s)
	api, err := utils.SendDataByApi(jsData, common.JChemApi.MolExport)
	if err != nil {
		return MolExportRes{}, err
	}
	var res MolExportRes
	err = json.Unmarshal(api, &res)
	if err != nil {
		return MolExportRes{}, err
	}
	return res, nil
}

type StructForName struct {
	InputFormat string `json:"inputFormat"`
	Parameters  string `json:"parameters"`
	Structure   string `json:"structure"`
}

func GetIupacName(s StructForName) (string, error) {
	type res struct {
		Structure string `json:"structure"`
		Format    string `json:"format"`
	}
	jsData, _ := json.Marshal(s)
	api, err := utils.SendDataByApi(jsData, common.JChemApi.MolExport)
	if err != nil {
		return "", err
	}
	var r res
	err = json.Unmarshal(api, &r)
	if err != nil {
		return "", err
	}
	return r.Structure, nil
}

type RxnCompoundData struct {
	Reactants []ReagentData `json:"reactants,omitempty"`
	Product   []ReagentData `json:"product,omitempty"`
}
type ReactionDetails struct {
	Msg  string          `json:"msg,omitempty"`
	Code int             `json:"code,omitempty"`
	Data RxnCompoundData `json:"data,omitempty"`
}

func GetRxnCompound(s CmpdInfo) (RxnCompoundData, error) {
	jsData, _ := json.Marshal(s)
	api, err := utils.SendDataByApi(jsData, common.JChemApi.GetReactionDetailsApi)
	if err != nil {
		return RxnCompoundData{}, err
	}
	var c ReactionDetails
	err = json.Unmarshal(api, &c)
	if err != nil {
		return RxnCompoundData{}, err
	}

	return c.Data, nil
}

func (rd *ReagentData) GetSmiles() string {
	info := CmpdInfo{
		Structure:   rd.Structure,
		InputFormat: "mrv",
		Parameters:  "smiles",
	}
	smiles, err := MolConvert(info)
	if err != nil {
		return rd.GetCxSmiles()
	}
	return smiles.Structure
}

func (rd *ReagentData) GetCxSmiles() string {
	info := CmpdInfo{
		Structure:   rd.Structure,
		InputFormat: "mrv",
		Parameters:  "cxsmiles",
	}
	cxSmiles, err := MolConvert(info)
	if err != nil {
		return err.Error()
	}
	return cxSmiles.Structure
}

func (rd *ReagentData) GetInchiKey() string {
	info := CmpdInfo{
		Structure:   rd.Structure,
		InputFormat: "mrv",
		Parameters:  "inchiKey",
	}
	inchiKey, err := MolConvert(info)
	if err != nil {
		return err.Error()
	}
	return inchiKey.Structure
}

func (rd *ReagentData) GetInchi() string {
	info := CmpdInfo{
		Structure:   rd.Structure,
		InputFormat: "mrv",
		Parameters:  "inchi",
	}
	inchi, err := MolConvert(info)
	if err != nil {
		return err.Error()
	}

	return strings.Split(inchi.Structure, "\n")[0]
}

type StereoDetail struct {
	TetraHedral []struct {
		AtomIndex int    `json:"atomIndex,omitempty"`
		Chirality string `json:"chirality,omitempty"`
	} `json:"tetraHedral,omitempty"`
}

// GetStereo http://192.168.2.139:8015/rest-v1/jws/stereo/cip
// s手性计算
func GetStereo(rd string) (string, int, error) {
	type s struct {
		Structure   string `json:"structure"`
		InputParams string `json:"inputParams"`
	}

	cmpd := s{
		Structure:   rd,
		InputParams: "mrv",
	}
	jsData, _ := json.Marshal(cmpd)
	api, err := utils.SendDataByApi(jsData, common.JChemApi.GetStereoApi)
	if err != nil {
		return "", 0, err
	}
	var detail StereoDetail
	err = json.Unmarshal(api, &detail)
	if err != nil {
		return "", 0, err
	}
	b, err := json.Marshal(detail.TetraHedral)
	if err != nil {
		panic(err)
	}
	return string(b), len(detail.TetraHedral), nil
}
