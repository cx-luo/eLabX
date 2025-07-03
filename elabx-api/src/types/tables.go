// Package types coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/28 16:58
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : tables.go
// @Software: GoLand
package types

import "time"

type ElnUsers struct {
	UserId          string    `json:"user_id" db:"user_id" gorm:"user_id"`       // 用户ID
	Username        string    `json:"username" db:"username" gorm:"username"`    // 用户名
	RealName        string    `json:"real_name" db:"real_name" gorm:"real_name"` // 用户昵称
	PasswordHash    string    `json:"password_hash" db:"password_hash" gorm:"password_hash"`
	Avatar          string    `json:"avatar" db:"avatar" gorm:"avatar"`                // 头像URL
	Roles           string    `json:"roles" db:"roles" gorm:"roles"`                   // 用户角色数组
	UserPermissions string    `json:"permissions" db:"permissions" gorm:"permissions"` // 用户权限数组
	Status          int64     `json:"status" db:"status" gorm:"status"`
	CreateTime      time.Time `json:"create_time" db:"create_time" gorm:"create_time"` // 创建时间
	UpdateTime      time.Time `json:"update_time" db:"update_time" gorm:"update_time"` // 更新时间
}

// TableName 表名称
func (*ElnUsers) TableName() string {
	return "eln_users"
}

// ElnRoutes undefined
type ElnRoutes struct {
	ID           int64  `json:"id" db:"id" gorm:"id"`
	Name         string `json:"name" db:"name" gorm:"name"`
	Path         string `json:"path" db:"path" gorm:"path"`
	Component    string `json:"component" db:"component" gorm:"component"`
	MetaIcon     string `json:"metaIcon" db:"meta_icon" gorm:"meta_icon"`
	MetaOrder    int64  `json:"metaOrder" db:"meta_order" gorm:"meta_order"`
	MetaTitle    string `json:"metaTitle" db:"meta_title" gorm:"meta_title"`
	MetaAffixTab int8   `json:"metaAffixTab" db:"meta_affix_tab" gorm:"meta_affix_tab"`
	ParentId     int64  `json:"parentId" db:"parent_id" gorm:"parent_id"`
	Status       int64  `json:"status" db:"status" gorm:"status"`
}

// TableName 表名称
func (*ElnRoutes) TableName() string {
	return "eln_routes"
}

// ElnApis undefined
type ElnApis struct {
	ID          int64     `json:"id" db:"id" gorm:"id"`
	ApiName     string    `json:"apiName" db:"api_name" gorm:"api_name"`
	ApiPath     string    `json:"apiPath" db:"api_path" gorm:"api_path"`
	Method      string    `json:"method" db:"method" gorm:"method"`
	ApiGroup    string    `json:"apiGroup" db:"api_group" gorm:"api_group"`
	ParentId    int64     `json:"parentId" db:"parent_id" gorm:"parent_id"`
	Description string    `json:"description" db:"description" gorm:"description"`
	CreateAt    time.Time `json:"createAt" db:"create_at" gorm:"create_at"` // 创建时间
	UpdateAt    time.Time `json:"updateAt" db:"update_at" gorm:"update_at"` // 更新时间
}

// TableName 表名称
func (*ElnApis) TableName() string {
	return "eln_apis"
}
