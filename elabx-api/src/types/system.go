// Package types coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/3 14:11
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : system.go
// @Software: GoLand
package types

import "gorm.io/gorm"

type SystemApiParam struct {
	ID          int64  `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    int64  `json:"parentId,omitempty"`
	Method      string `json:"method,omitempty"`
	Path        string `json:"path"`
}

type SystemMenuParam struct {
	gorm.Model
	ID        int64  `json:"id,omitempty" db:"id" gorm:"id"`
	Type      string `json:"type,omitempty"`
	Component string `json:"component,omitempty"`
	Status    int    `json:"status,omitempty"`
	Meta      Meta   `gorm:"embedded" json:"meta,omitempty"`
	RouteName string `json:"routeName,omitempty"`
	ParentID  int    `json:"parentId,omitempty"`
	Path      string `json:"path,omitempty"`
}
