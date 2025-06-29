// Package types coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/28 16:58
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : elabx_user.go
// @Software: GoLand
package types

import "time"

type ElnUsers struct {
	ID           int64     `json:"id" db:"id" gorm:"id"`
	Userid       int64     `json:"userid" db:"userid" gorm:"userid"`
	Username     string    `json:"username" db:"username" gorm:"username"`
	Email        string    `json:"email" db:"email" gorm:"email"`
	Phone        string    `json:"phone" db:"phone" gorm:"phone"`
	PasswordHash string    `json:"password_hash" db:"password_hash" gorm:"password_hash"`
	Status       int64     `json:"status" db:"status" gorm:"status"`
	GroupId      int64     `json:"group_id" db:"group_id" gorm:"group_id"`
	Avatar       string    `json:"avatar" db:"avatar" gorm:"avatar"`
	Permissions  string    `json:"permissions" db:"permissions" gorm:"permissions"`
	Authn        string    `json:"authn" db:"authn" gorm:"authn"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"updated_at"`
}

// TableName 表名称
func (*ElnUsers) TableName() string {
	return "eln_users"
}
