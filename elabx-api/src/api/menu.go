// Package api coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/28 13:20
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : menu.go
// @Software: GoLand
package api

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Meta represents the meta section of the JSON structure.
type Meta struct {
	Icon     string `json:"icon" gorm:"column:meta_icon"`
	Order    int64  `json:"order" gorm:"column:meta_order"`
	Title    string `json:"title" gorm:"column:meta_title"`
	AffixTab bool   `json:"affixTab,omitempty" gorm:"column:meta_affix_tab"` // Optional field
}

// Child represents a child route within the main route.
//type Child struct {
//	ID        uint   `gorm:"primaryKey;autoIncrement"`
//	Name      string `json:"name" gorm:"column:child_name"`
//	Path      string `json:"path" gorm:"column:child_path"`
//	Component string `json:"component" gorm:"column:child_component"`
//	Icon      string `json:"icon" gorm:"column:child_meta_icon"`
//	Order     int    `json:"order" gorm:"column:child_meta_order"`
//	Title     string `json:"title" gorm:"column:child_meta_title"`
//	AffixTab  bool   `json:"affixTab,omitempty" gorm:"column:child_meta_affix_tab"` // Optional field
//	ParentID  uint   `gorm:"column:parent_id"`
//}

// Route represents the entire JSON structure.
type Route struct {
	ID        int64   `gorm:"primaryKey;autoIncrement"`
	Meta      Meta    `json:"meta" gorm:"-"`
	Name      string  `json:"name" gorm:"column:name"`
	Path      string  `json:"path" gorm:"column:path"`
	Component string  `json:"component" gorm:"column:component"`
	Children  []Route `json:"children,omitempty" gorm:"-"`
	ParentID  int64   `gorm:"column:parent_id"`
}

func genRoutesFromTable(routes []types.ElnRoutes) []Route {
	var components []Route
	for _, elnRoutes := range routes {
		components = append(components, Route{
			ID:        elnRoutes.ID,
			Meta:      Meta{elnRoutes.MetaIcon, elnRoutes.MetaOrder, elnRoutes.MetaTitle, elnRoutes.MetaAffixTab == 1},
			Name:      elnRoutes.Name,
			Path:      elnRoutes.Path,
			Component: elnRoutes.Component,
			Children:  nil,
			ParentID:  elnRoutes.ParentId,
		})
	}
	return components
}
func getParentRoutes(db *gorm.DB) ([]Route, error) {
	var routes []types.ElnRoutes
	err := db.Table("eln_routes").Select("id", "name", "path", "component", "meta_icon", "meta_order",
		"meta_title", "meta_affix_tab").Where(`parent_id = 0`).Find(&routes).Error
	if err != nil {
		return nil, err
	}
	r := genRoutesFromTable(routes)
	return r, nil
}

func getChildrenRoutes(db *gorm.DB, parentId int64) ([]Route, error) {
	var children []types.ElnRoutes
	err := db.Table("eln_routes").Select("id", "name", "path", "component", "meta_icon", "meta_order",
		"meta_title", "meta_affix_tab", "parent_id").Where(`parent_id = ?`, parentId).Find(&children).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	r := genRoutesFromTable(children)

	return r, nil
}

func defaultRoute() []Route {
	parentRoutes, err := getParentRoutes(dao.OBCursor)
	if err != nil {
		zap.L().Error(fmt.Sprintf("genrouter error: %s", err.Error()))
		return nil
	}

	for i, _ := range parentRoutes {
		children, err := getChildrenRoutes(dao.OBCursor, parentRoutes[i].ID)
		if err != nil {
			zap.L().Error(fmt.Sprintf("failed to get children routes for parent ID %d: %s", parentRoutes[i].ID, err.Error()))
			continue
		}
		parentRoutes[i].Children = children
	}

	return parentRoutes
}

func GetRouteTree(c *gin.Context) {
	var routes []types.ElnRoutes
	err := dao.OBCursor.Table("eln_routes").Find(&routes).Error
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to get all route: %s", err.Error()))
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"items": routes, "total": len(routes)})
	return
}

func GetUserRouteList(c *gin.Context) {
	var roles struct {
		Userid       int    `json:"userid,omitempty" db:"userid"`
		AuthorityIds string `json:"authorityIds,omitempty" db:"permissions"`
	}
	err := c.ShouldBind(&roles)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"items": defaultRoute()})
	return
}

func UpdateMenu(c *gin.Context) {
	var menu types.SystemMenuParam
	err := c.ShouldBind(&menu)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Model(&types.ElnRouteMenus{}).Where(`id = ?`, menu.ID).Updates(menu).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"items": defaultRoute()})
	return
}
