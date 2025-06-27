// Package api coding=utf-8
// @Project : server
// @Time    : 2024/6/20 14:03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : file_server.go
// @Software: GoLand
package api

import (
	"crypto/md5"
	"eLabX/src/utils"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
)

func getFileMd5(file *multipart.FileHeader) (string, error) {
	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()
	// 复制文件内容到本地文件
	hash := md5.New()
	if _, err := io.Copy(hash, openedFile); err != nil {
		return "", err
	}

	fileMD5 := hash.Sum(nil)
	// 保存文件到本地
	calculatedMD5 := hex.EncodeToString(fileMD5)
	return calculatedMD5, nil
}

func UploadFile(c *gin.Context) {
	// 确保请求是 multipart/form-data 类型
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		utils.BadRequestErr(c, errors.New("Invalid form data."+err.Error()))
		return
	}

	// 从请求中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestErr(c, errors.New("No file provided."+err.Error()))
		return
	}

	// 保存文件到本地
	calculatedMD5, err := getFileMd5(file)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	defer openedFile.Close()

	// todo: 需要把文件传到分布式文件系统，返回文件系统文件的url
	dst, err := os.Create(utils.GlobalConfig.Service.UploadFileDir + file.Filename)
	defer dst.Close()

	if _, err := io.Copy(dst, openedFile); err != nil {
		utils.InternalRequestErr(c, errors.New("Failed to save file"+err.Error()))
		return
	}

	utils.SuccessWithData(c, "File uploaded successfully", gin.H{"fileMd5": calculatedMD5, "fileName": file.Filename, "spectrumUrl": file.Filename})
	return
}

// SaveFilenameToDb
// todo: 需要把文件传到分布式文件系统，返回文件系统文件的url
func SaveFilenameToDb(c *gin.Context) {
	var f struct {
		ReportRank   int64  `json:"reportRank,omitempty"`
		ProductId    int64  `json:"productId,omitempty"`
		ReactionId   int64  `json:"reactionId,omitempty"`
		FileType     string `json:"fileType,omitempty"`
		FileName     string `json:"fileName,omitempty"`
		ProductAlias string `json:"productAlias,omitempty"`
		Comments     string `json:"comments"`
	}
	err := c.ShouldBind(&f)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	switch f.FileType {
	case "lcms":
		return
	case "nmr":
		return
	case "other":
		return
	}
	utils.Success(c, "")
	return
}
