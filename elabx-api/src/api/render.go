// Package api coding=utf-8
// @Project : server
// @Time    : 2024/6/27 11:15
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : render.go
// @Software: GoLand
package api

import (
	"eLabX/src/common"
	"eLabX/src/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

func GenImgBase64String(smi string) (error, string) {
	noneImg := "data:image/png;base64," + "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"

	if smi == "" {
		return nil, noneImg
	}

	//	优先使用 'http://192.168.2.139:8015/rest-v1/jws/molconvert'
	// example https://docs.chemaxon.com/display/lts-lithium/png.md

	var inputFormat string
	if strings.Contains(smi, "<?xml") {
		inputFormat = "mrv"
	} else {
		inputFormat = "smiles"
	}
	imgFromMrv, err := GenImgFromMrv(smi, inputFormat)
	if err != nil {
		return err, noneImg
	}

	return nil, imgFromMrv
}

func GenImgFromMrv(smi string, inputFormat string) (string, error) {
	s := CmpdInfo{
		Structure:   smi,
		InputFormat: inputFormat,
		Parameters:  "png:w850,h320,b32",
	}
	img, err := MolConvert(s)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + string(img.BinaryStructure), nil
}

func GenImgFromIndigo(smi string) (string, error) {
	s := struct {
		OutputFormat string `json:"output_format"`
		Struct       string `json:"struct"`
	}{
		"image/png;base64",
		smi,
	}
	jsData, _ := json.Marshal(s)
	resp, err := utils.SendDataByApi(jsData, common.Indigo.Render)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + string(resp), nil
}

func GenImg(c *gin.Context) {
	var smi struct {
		Smiles      string `json:"smiles,omitempty"`
		CdStructure string `json:"cdStructure,omitempty"`
	}
	err := c.ShouldBind(&smi)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	var imgStr, structure string
	if smi.CdStructure != "" {
		structure = smi.CdStructure
	} else {
		structure = smi.Smiles
	}
	err, imgStr = GenImgBase64String(structure)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"img": imgStr})
	return
}
