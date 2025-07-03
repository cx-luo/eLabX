// Package types coding=utf-8
// @Project : elabx-api
// @Time    : 2025/7/3 14:11
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : system.go
// @Software: GoLand
package types

type ApiParam struct {
	ID    int64 `json:"id,omitempty"`
	Value Param `json:"param,omitempty"`
}

type Param struct {
	Description string `json:"description,omitempty"`
	ParentID    int64  `json:"parentId,omitempty"`
	Method      string `json:"method,omitempty"`
	Path        string `json:"path,omitempty"`
}
