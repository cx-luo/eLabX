// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2024/5/7 15:01
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : response.go
// @Software: GoLand
package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type BaseResponse struct {
	StatusCode int         `json:"statusCode"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data,omitempty"`
}

func Error(c *gin.Context, status int, code int, message string) {
	zap.L().Error(message)
	c.AbortWithStatusJSON(status, BaseResponse{
		StatusCode: code,
		Msg:        message,
		Data:       nil,
	})
}
func Success(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, BaseResponse{
		StatusCode: 200, Msg: msg,
	})
}

func BadRequestErr(c *gin.Context, err error) {
	//zap.L().Error("BadRequest: " + err.Error())
	Error(c, http.StatusBadRequest, http.StatusBadRequest, err.Error())
}

func InternalRequestErr(c *gin.Context, err error) {
	//zap.L().Error("InternalServerError: " + err.Error())
	Error(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
}

func SuccessWithData(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, BaseResponse{
		StatusCode: 200, Msg: msg, Data: data,
	})
}

func NotFoundError(c *gin.Context, err error) {
	//Logger.Error("NotFound: " + err.Error())
	Error(c, http.StatusNotFound, http.StatusNotFound, err.Error())
}
