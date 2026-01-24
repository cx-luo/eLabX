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
	UserId       int64     `json:"userId" db:"user_id" gorm:"column:user_id;primaryKey"` // 用户ID
	Username     string    `json:"username" db:"username" gorm:"username"`               // 用户名
	RealName     string    `json:"realName" db:"real_name" gorm:"real_name"`             // 用户昵称
	Email        string    `json:"email" db:"email" gorm:"email"`
	IpAddr       string    `json:"ipAddr" db:"ip_addr" gorm:"ip_addr"`
	PasswordHash string    `json:"passwordHash" db:"password_hash" gorm:"password_hash"`
	Avatar       string    `json:"avatar" db:"avatar" gorm:"avatar"` // 头像URL
	Roles        string    `json:"roles" db:"roles" gorm:"roles"`    // 用户角色数组
	Status       int64     `json:"status" db:"status" gorm:"status"`
	CreateAt     time.Time `json:"createAt" db:"create_at" gorm:"create_at"` // 创建时间
	UpdateAt     time.Time `json:"updateAt" db:"update_at" gorm:"update_at"` // 更新时间
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
	Status      int16     `json:"status" db:"status" gorm:"status"`
	Permissions string    `json:"permissions" db:"permissions" gorm:"permissions"`    // 用户权限数组
	CreateAt    time.Time `json:"createAt,omitempty" db:"create_at" gorm:"create_at"` // 创建时间
	UpdateAt    time.Time `json:"updateAt,omitempty" db:"update_at" gorm:"update_at"` // 更新时间
}

// TableName 表名称
func (*ElnProject) TableName() string {
	return "eln_project"
}

// ElnGroup represents a user group in the system
type ElnGroup struct {
	GroupId     int64     `json:"groupId" db:"group_id" gorm:"group_id"`
	GroupName   string    `json:"groupName" db:"group_name" gorm:"group_name"`
	Description string    `json:"description" db:"description" gorm:"description"`
	Status      int8      `json:"status" db:"status" gorm:"status"`
	CreatedBy   int64     `json:"createdBy" db:"created_by" gorm:"created_by"`
	Permissions string    `json:"permissions" db:"permissions" gorm:"permissions"`    // permission array
	CreateAt    time.Time `json:"createAt,omitempty" db:"create_at" gorm:"create_at"` // creation time
	UpdateAt    time.Time `json:"updateAt,omitempty" db:"update_at" gorm:"update_at"` // update time
}

// TableName returns the table name for ElnGroup
func (*ElnGroup) TableName() string {
	return "eln_group"
}

// ElnUserGroup represents the association between a user and a group
type ElnUserGroup struct {
	UserId    int64     `json:"userId" db:"user_id" gorm:"column:user_id;primaryKey"`    // 用户ID
	GroupId   int64     `json:"groupId" db:"group_id" gorm:"column:group_id;primaryKey"` // 组ID
	Status    int8      `json:"status" db:"status" gorm:"column:status;default:1"`       // 状态 (1:启用, 0:禁用)
	CreatedBy int64     `json:"createdBy" db:"created_by" gorm:"column:created_by"`      // 创建人ID
	CreateAt  time.Time `json:"createAt" db:"create_at" gorm:"column:create_at"`         // 创建时间
	UpdateAt  time.Time `json:"updateAt" db:"update_at" gorm:"column:update_at"`         // 更新时间
}

// TableName returns the table name for ElnUserGroup
func (*ElnUserGroup) TableName() string {
	return "eln_user_group"
}

// // ElnGroupPermission represents the association between a group and a permission
// type ElnGroupPermission struct {
// 	GroupPermissionId int64     `json:"groupPermissionId" db:"group_permission_id" gorm:"primaryKey;column:group_permission_id"`
// 	Groups           int64     `json:"groupId" db:"group_id" gorm:"column:group_id"`
// 	PermissionId      int64     `json:"permissionId" db:"permission_id" gorm:"column:permission_id"`
// 	CreateAt          time.Time `json:"createAt,omitempty" db:"create_at" gorm:"column:create_at"`
// 	UpdateAt          time.Time `json:"updateAt,omitempty" db:"update_at" gorm:"update_at"` // update time
// }

// // TableName returns the table name for ElnGroupPermission
// func (*ElnGroupPermission) TableName() string {
// 	return "eln_group_permission"
// }

// ElnPermission represents a permission in the system
type ElnPermission struct {
	PermissionId   int64     `json:"permissionId" db:"permission_id" gorm:"primaryKey;column:permission_id"`
	PermissionName string    `json:"permissionName" db:"permission_name" gorm:"column:permission_name"`
	PermissionType string    `json:"permissionType" db:"permission_type" gorm:"column:permission_type"`
	Description    string    `json:"description" db:"description" gorm:"column:description"`
	Status         int8      `json:"status" db:"status" gorm:"column:status;default:1"`
	CreatedBy      int64     `json:"createdBy" db:"created_by" gorm:"column:created_by"`
	CreateAt       time.Time `json:"createAt,omitempty" db:"create_at" gorm:"column:create_at"`
	UpdateAt       time.Time `json:"updateAt,omitempty" db:"update_at" gorm:"column:update_at"`
}

// TableName returns the table name for ElnPermission
func (*ElnPermission) TableName() string {
	return "eln_permission"
}

// ElnUserMessage UserMessage represents a message or notification for a user
type ElnUserMessage struct {
	ID       int64     `json:"id" db:"id" gorm:"primaryKey;column:id"`
	UserId   int64     `json:"userId" db:"user_id" gorm:"column:user_id"`
	Title    string    `json:"title" db:"title" gorm:"column:title"`
	Content  string    `json:"content" db:"content" gorm:"column:content"`
	IsRead   bool      `json:"isRead" db:"is_read" gorm:"column:is_read;default:false"`
	Type     string    `json:"type" db:"type" gorm:"column:type"` // e.g., "notification", "alert", "reminder"
	CreateAt time.Time `json:"createAt" db:"create_at" gorm:"column:create_at"`
}

// TableName returns the table name for UserMessage
func (*ElnUserMessage) TableName() string {
	return "eln_user_message"
}

// ElnReaction undefined
type ElnReaction struct {
	ReactionId     int64     `json:"reactionId" gorm:"reaction_id"`
	ReactionSmiles string    `json:"reactionSmiles" gorm:"reaction_smiles"`
	CxSmiles       string    `json:"cxSmiles" gorm:"cx_smiles"`
	CdStructure    string    `json:"cdStructure" gorm:"cd_structure"`
	GmtCreate      time.Time `json:"gmtCreate" gorm:"gmt_create"`
	GmtModified    time.Time `json:"gmtModified" gorm:"gmt_modified"`
}

// TableName 表名称
func (*ElnReaction) TableName() string {
	return "eln_reaction"
}

// ElnRxnReagents undefined
type ElnRxnReagents struct {
	ReagentId        int64     `json:"reagentId" gorm:"reagent_id"`
	ReactionId       int64     `json:"reactionId" gorm:"reaction_id"`
	ReagentName      string    `json:"reagentName" gorm:"reagent_name"`
	ReagentSmiles    string    `json:"reagentSmiles" gorm:"reagent_smiles"`
	Mw               float32   `json:"mw" gorm:"mw"`
	MonoisotopicMass float32   `json:"monoisotopicMass" gorm:"monoisotopic_mass"`
	Formula          string    `json:"formula" gorm:"formula"`
	ReagentRole      string    `json:"reagentRole" gorm:"reagent_role"`    // role
	Equiv            float32   `json:"equiv" gorm:"equiv"`                 // eq
	Cas              string    `json:"cas" gorm:"cas"`                     // CAS#
	Concentration    float32   `json:"concentration" gorm:"concentration"` // Conv. 浓度
	CdStructure      string    `json:"cdStructure" gorm:"cd_structure"`
	Inchi            string    `json:"inchi" gorm:"inchi"`
	Inchikey         string    `json:"inchikey" gorm:"inchikey"`
	Cxsmiles         string    `json:"cxsmiles" gorm:"cxsmiles"`
	Density          float32   `json:"density" gorm:"density"`
	Quantity         float32   `json:"quantity" gorm:"quantity"`
	QuantityUnit     string    `json:"quantityUnit" gorm:"quantity_unit"`
	Purity           float32   `json:"purity" gorm:"purity"`
	CompoundId       int64     `json:"compoundId" gorm:"compound_id"`
	IsChiral         int8      `json:"isChiral" gorm:"is_chiral"`
	ReagentHash      int64     `json:"reagentHash" gorm:"reagent_hash"`
	IsLimiting       int8      `json:"isLimiting" gorm:"is_limiting"`
	StereoCentersCnt int64     `json:"stereoCentersCnt" gorm:"stereo_centers_cnt"`
	ProductAlias     string    `json:"productAlias" gorm:"product_alias"`
	ChiralDescriptor string    `json:"chiralDescriptor" gorm:"chiral_descriptor"`
	Moles            float32   `json:"moles" gorm:"moles"`
	MolesUnit        string    `json:"molesUnit" gorm:"moles_unit"`
	Volume           float32   `json:"volume" gorm:"volume"`
	VolumeUnit       string    `json:"volumeUnit" gorm:"volume_unit"`
	IsSolution       int8      `json:"isSolution" gorm:"is_solution"`
	GmtCreate        time.Time `json:"gmtCreate" gorm:"gmt_create"`
	GmtModified      time.Time `json:"gmtModified" gorm:"gmt_modified"`
}

// TableName 表名称
func (*ElnRxnReagents) TableName() string {
	return "eln_rxn_reagents"
}
