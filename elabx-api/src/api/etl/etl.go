// Package etl coding=utf-8
// @Project : elabx-vben
// @Time    : 2025/7/9 11:57
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : etl.go
// @Software: GoLand
package etl

import (
	"eLabX/src/dao"
	"eLabX/src/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

// GetDatabase handles GET /api/etl/database/list
func GetDatabase(c *gin.Context) {
	type DatabaseNameResult struct {
		Database string `gorm:"column:schema_name"`
	}
	type list struct {
		Value string `json:"value,omitempty"`
		Label string `json:"label,omitempty"`
	}
	var dbs []list
	var results []DatabaseNameResult

	// Query all database names from information_schema.schemata
	if err := dao.OBCursor.Table("information_schema.schemata").
		Select("schema_name").
		Scan(&results).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to fetch databases: "+err.Error()))
		return
	}

	for _, r := range results {
		dbs = append(dbs, list{
			Value: r.Database,
			Label: r.Database,
		})
	}

	utils.SuccessWithData(c, "", gin.H{"items": dbs})
	return

}

// GetTableList handles GET /api/etl/table/list/:dbName
func GetTableList(c *gin.Context) {
	dbName := c.Param("dbName")

	type TableNameResult struct {
		TableName string `gorm:"column:table_name"`
	}

	var tables []string
	var results []TableNameResult

	// Replace "DB" with your actual *gorm.DB instance
	if err := dao.OBCursor.Table("information_schema.tables").
		Select("table_name").
		Where("table_schema = ?", dbName).
		Scan(&results).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to fetch tables: "+err.Error()))

		return
	}

	for _, r := range results {
		tables = append(tables, r.TableName)
	}

	utils.SuccessWithData(c, "", gin.H{"items": tables})
	return
}

// GetTableColumns handles GET /api/etl/table/columns/:dbName/:tableName
func GetTableColumns(c *gin.Context) {
	dbName := c.Param("dbName")
	tableName := c.Param("tableName")

	type ColumnResult struct {
		ColumnName    string  `gorm:"column:column_name" json:"columnName"`
		DataType      string  `gorm:"column:data_type" json:"dataType"`
		IsNullable    string  `gorm:"column:is_nullable" json:"isNullable"`
		ColumnDefault *string `gorm:"column:column_default" json:"columnDefault"`
		ColumnComment string  `gorm:"column:column_comment" json:"columnComment"`
	}

	var columns []ColumnResult

	if err := dao.OBCursor.Table("information_schema.columns").
		Select("column_name, data_type, is_nullable, column_default, column_comment").
		Where("table_schema = ? AND table_name = ?", dbName, tableName).
		Scan(&columns).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to fetch columns: "+err.Error()))
		return
	}

	utils.SuccessWithData(c, "", gin.H{"items": columns})
	return
}
