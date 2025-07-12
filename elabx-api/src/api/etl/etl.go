// Package etl coding=utf-8
// @Project : eLabX
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

// GetTableData handles GET /api/etl/table/data/:dbName/:tableName
func GetTableData(c *gin.Context) {
	dbName := c.Param("dbName")
	tableName := c.Param("tableName")

	var req struct {
		Page      int      `json:"page,omitempty"`
		PageSize  int      `json:"pageSize,omitempty"`
		Columns   []string `json:"columns,omitempty"`
		SortField string   `json:"sortField,omitempty"`
		SortOrder string   `json:"sortOrder,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, errors.New("invalid request body: "+err.Error()))
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 60
	}
	offset := (page - 1) * pageSize

	sortField := req.SortField
	sortOrder := req.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	fullTableName := dbName + "." + tableName

	// Count total number of records
	var total int64
	if err := dao.OBCursor.Table(fullTableName).Count(&total).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to count table data: "+err.Error()))
		return
	}

	// Query data
	var results []map[string]interface{}
	if len(req.Columns) == 0 {
		// 查询所有列
		req.Columns = []string{"*"}
	}
	query := dao.OBCursor.Table(fullTableName).Select(req.Columns)

	// Support sorting
	if sortField != "" {
		query = query.Order(sortField + " " + sortOrder)
	}

	// Support pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to fetch table data: "+err.Error()))
		return
	}

	utils.SuccessWithData(c, "", gin.H{
		"items":    results,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
	return
}
