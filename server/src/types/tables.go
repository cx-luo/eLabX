// Package types coding=utf-8
// @Project : eLabX
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
	CreateTime      time.Time `json:"create_time" db:"create_time" gorm:"create_time;default:(-)"` // 创建时间
	UpdateTime      time.Time `json:"update_time" db:"update_time" gorm:"update_time;default:(-)"` // 更新时间
}

// TableName 表名称
func (*ElnUsers) TableName() string {
	return "eln_users"
}

// ElnRouteMenus undefined
type ElnRouteMenus struct {
	ID        int64     `json:"id" db:"id" gorm:"id"`
	RouteName string    `json:"routeName" db:"route_name" gorm:"route_name"`
	Path      string    `json:"path" db:"path" gorm:"path"`
	Type      string    `json:"type" db:"type" gorm:"type"`
	Component string    `json:"component" db:"component" gorm:"component"`
	Status    int8      `json:"status" db:"status" gorm:"status"`
	Meta      Meta      `json:"meta,omitempty" gorm:"embedded"`
	ParentId  int64     `json:"parentId" db:"parent_id" gorm:"parent_id"`
	CreateAt  time.Time `json:"createAt" db:"create_at" gorm:"create_at;default:(-)"` // 创建时间
	UpdateAt  time.Time `json:"updateAt" db:"update_at" gorm:"update_at;default:(-)"` // 更新时间
}

type Meta struct {
	Name               string `json:"name,omitempty" db:"name" gorm:"name"`
	Icon               string `json:"icon,omitempty" db:"icon" gorm:"icon"`
	Order              int64  `json:"order,omitempty" gorm:"order"`
	AffixTab           int8   `json:"affixTab,omitempty" db:"affix_tab" gorm:"affix_tab"`
	HideChildrenInMenu int8   `json:"hideChildrenInMenu,omitempty" db:"hide_children_in_menu" gorm:"hide_children_in_menu"`
	HideInBreadcrumb   int8   `json:"hideInBreadcrumb,omitempty" db:"hide_in_breadcrumb" gorm:"hide_in_breadcrumb"`
	HideInMenu         int8   `json:"hideInMenu,omitempty" db:"hide_in_menu" gorm:"hide_in_menu"`
	HideInTab          int8   `json:"hideInTab,omitempty" db:"hide_in_tab" gorm:"hide_in_tab"`
	KeepAlive          int8   `json:"keepAlive,omitempty" db:"keep_alive" gorm:"keep_alive"`
}

// TableName 表名称
func (*ElnRouteMenus) TableName() string {
	return "eln_route_menus"
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
	CreateAt    time.Time `json:"createAt" db:"create_at" gorm:"create_at;default:(-)"` // 创建时间
	UpdateAt    time.Time `json:"updateAt" db:"update_at" gorm:"update_at;default:(-)"` // 更新时间
}

// TableName 表名称
func (*ElnApis) TableName() string {
	return "eln_apis"
}

// ElnRoles undefined
type ElnRoles struct {
	ID       int64     `json:"id" db:"id" gorm:"id"`
	Name     string    `json:"name" db:"name" gorm:"name"`
	Code     int64     `json:"code" db:"code" gorm:"code"`
	Status   int8      `json:"status" db:"status" gorm:"status"`
	Sort     int64     `json:"sort" db:"sort" gorm:"sort"`
	ApiId    string    `json:"apiId" db:"api_id" gorm:"api_id"`
	AuthId   string    `json:"authId" db:"auth_id" gorm:"auth_id"`
	Remark   string    `json:"remark" db:"remark" gorm:"remark"`                     // comments
	CreateAt time.Time `json:"createAt" db:"create_at" gorm:"create_at;default:(-)"` // 创建时间
	UpdateAt time.Time `json:"updateAt" db:"update_at" gorm:"update_at;default:(-)"` // 更新时间
}

// TableName 表名称
func (*ElnRoles) TableName() string {
	return "eln_roles"
}

// ElnProject undefined
type ElnProject struct {
	ProjectId   int64     `json:"projectId" db:"project_id" gorm:"project_id"`
	ProjectName string    `json:"projectName" db:"project_name" gorm:"project_name"`
	CreatedBy   int64     `json:"createdBy" db:"created_by" gorm:"created_by"`
	Description string    `json:"description" db:"description" gorm:"description"`
	Status      int8      `json:"status" db:"status" gorm:"status"`
	Permissions string    `json:"permissions" db:"permissions" gorm:"permissions"` // 用户权限数组
	CreateAt    time.Time `json:"createAt" db:"create_at" gorm:"create_at"`        // 创建时间
	UpdateAt    time.Time `json:"updateAt" db:"update_at" gorm:"update_at"`        // 更新时间
}

// TableName 表名称
func (*ElnProject) TableName() string {
	return "eln_project"
}
