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

// GetTableColumnsWithPK handles GET /api/etl/table/columns/:dbName/:tableName and adds primary key info
func GetTableColumnsWithPK(c *gin.Context) {
	dbName := c.Param("dbName")
	tableName := c.Param("tableName")

	type ColumnResult struct {
		ColumnName    string  `gorm:"column:column_name" json:"columnName"`
		DataType      string  `gorm:"column:data_type" json:"dataType"`
		IsNullable    string  `gorm:"column:is_nullable" json:"isNullable"`
		ColumnDefault *string `gorm:"column:column_default" json:"columnDefault"`
		ColumnComment string  `gorm:"column:column_comment" json:"columnComment"`
		IsPrimaryKey  bool    `json:"isPrimaryKey"`
	}

	var columns []ColumnResult

	// 查询所有列
	if err := dao.OBCursor.Table("information_schema.columns").
		Select("column_name, data_type, is_nullable, column_default, column_comment").
		Where("table_schema = ? AND table_name = ?", dbName, tableName).
		Scan(&columns).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to fetch columns: "+err.Error()))
		return
	}

	// 查询主键列
	var pkColumns []string
	if err := dao.OBCursor.Table("information_schema.key_column_usage").
		Select("column_name").
		Where("table_schema = ? AND table_name = ? AND constraint_name = 'PRIMARY'", dbName, tableName).
		Pluck("column_name", &pkColumns).Error; err != nil {
		utils.InternalRequestErr(c, errors.New("failed to fetch primary key columns: "+err.Error()))
		return
	}
	pkSet := make(map[string]struct{}, len(pkColumns))
	for _, pk := range pkColumns {
		pkSet[pk] = struct{}{}
	}

	// 标记主键
	for i := range columns {
		if _, ok := pkSet[columns[i].ColumnName]; ok {
			columns[i].IsPrimaryKey = true
		}
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

func UpdateTableDataApi(c *gin.Context) {
	var req struct {
		DBName     string                 `json:"dbName"`
		TableName  string                 `json:"tableName"`
		PrimaryKey []string               `json:"primaryKey"`
		Data       map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	if req.DBName == "" || req.TableName == "" || len(req.PrimaryKey) == 0 || req.Data == nil {
		utils.BadRequestErr(c, errors.New("missing required fields"))
		return
	}

	// Collect primary key values from data
	pkValues := make([]interface{}, 0, len(req.PrimaryKey))
	for _, pk := range req.PrimaryKey {
		val, ok := req.Data[pk]
		if !ok {
			utils.BadRequestErr(c, errors.New("primary key value missing in data: "+pk))
			return
		}
		pkValues = append(pkValues, val)
	}

	fullTableName := req.DBName + "." + req.TableName

	// Remove primary keys from update data
	updateData := make(map[string]interface{})
	for k, v := range req.Data {
		isPK := false
		for _, pk := range req.PrimaryKey {
			if k == pk {
				isPK = true
				break
			}
		}
		if !isPK {
			updateData[k] = v
		}
	}

	if len(updateData) == 0 {
		utils.BadRequestErr(c, errors.New("no fields to update"))
		return
	}

	// Build where clause for composite primary key
	query := dao.OBCursor.Table(fullTableName)
	for i, pk := range req.PrimaryKey {
		if i == 0 {
			query = query.Where(pk+" = ?", pkValues[i])
		} else {
			query = query.Where(pk+" = ?", pkValues[i])
		}
	}

	result := query.Updates(updateData)
	if result.Error != nil {
		utils.InternalRequestErr(c, errors.New("failed to update table data: "+result.Error.Error()))
		return
	}
	utils.SuccessWithData(c, "update success", gin.H{"rowsAffected": result.RowsAffected})
	return
}
