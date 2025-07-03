// Package system coding=utf-8
// @Project : elabx-api
// @Time    : 2025/7/3 10:09
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : apis.go
// @Software: GoLand
package system

import (
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"hash/crc32"
	"strings"
)

func letterToNumber(s string) int64 {
	hash := crc32.ChecksumIEEE([]byte(s))
	return int64(hash) % 1000000000
}

func RefreshApis(c *gin.Context) {
	for _, a := range utils.Apis {
		segments := strings.Split(a.Path, "/") // 最多切分为3段
		group := "default"
		if len(segments) > 2 {
			group = segments[2]
		}

		funcSegments := strings.Split(a.Handler, "/")
		node, _ := snowflake.NewNode(1)
		id := node.Generate().Time() % 10000
		api := &types.ElnApis{
			ID:          id,
			ApiName:     funcSegments[len(funcSegments)-1],
			ApiPath:     a.Path,
			Method:      a.Method,
			Description: segments[len(segments)-1],
			ApiGroup:    group,
			ParentId:    letterToNumber(group),
		}

		// 在 `id` 冲突时将列更新为新值
		result := dao.OBCursor.Select("api_path", "api_name", "id", "method", "description", "api_group", "parent_id").Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "api_path"}}, // key colume
			// 也可以用 map [ string ] interface {} { "role" : "user" }
			DoUpdates: clause.AssignmentColumns([]string{"id", "api_name", "api_group"}), // column needed to be updated
		}).Create(&api)

		if result.Error != nil {
			zap.L().Error(fmt.Sprintf("failed to update apis: %v\n", result.Error))
		}
	}
	utils.Success(c, "")
	return
}

func GetApiList(c *gin.Context) {
	var apiList []types.ElnApis
	//RefreshApis(c)
	err := dao.OBCursor.Model(&types.ElnApis{}).Find(&apiList).Error
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to update apis: %v\n", err))
	}

	utils.SuccessWithData(c, "", gin.H{"items": apiList, "total": len(apiList)})
	return
}

func AddApi(c *gin.Context) {
	var apiData types.Param
	err := c.ShouldBind(&apiData)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	node, _ := snowflake.NewNode(1)
	id := node.Generate().Time() % 1000000000

	err = dao.OBCursor.Select("api_path", "id", "method", "description", "parent_id").Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "api_path"}}, // key colume
	}).Create(&types.ElnApis{
		ID:          id,
		ApiPath:     apiData.Path,
		ParentId:    apiData.ParentID,
		Description: apiData.Description,
		Method:      apiData.Method,
	}).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"id": id})
	return
}

func UpdateAPi(c *gin.Context) {
	var apiData types.ApiParam
	err := c.ShouldBind(&apiData)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Where(`id = ?`, apiData.ID).Updates(&types.ElnApis{
		ApiPath:     apiData.Value.Path,
		ParentId:    apiData.Value.ParentID,
		Description: apiData.Value.Description,
		Method:      apiData.Value.Method,
	}).Error

	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}

func DeleteAPi(c *gin.Context) {
	var apiId struct {
		ID int `json:"id"`
	}
	err := c.ShouldBind(&apiId)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Where(`id = ?`, apiId.ID).Delete(&types.ElnApis{}).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.Success(c, "")
	return
}
