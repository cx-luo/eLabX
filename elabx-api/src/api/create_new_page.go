// Package api coding=utf-8
// @Project : server
// @Time    : 2025/4/8 16:23
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : create_new_page.go
// @Software: GoLand
package api

import (
	"database/sql"
	"eLabX/src/common"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateNewPage(c *gin.Context) {
	var rxn common.ElnRxnBasicInfo
	if err := c.ShouldBind(&rxn); err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	// 设置ReactionId
	rxn.ReactionId = utils.GenerateSnowflakeID()

	// 开始事务
	tx, err := dao.OBCursor.Beginx()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 确保在函数结束时正确地提交或回滚事务
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 插入eln_rxn_basicinfo表
	if err = InsertElnRxnBasicInfo(tx, rxn); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 查询或插入eln_project表，并获取projectId和rankId
	projectId, rankId, err := GetOrCreateProjectIdAndRankId(tx, rxn.ProjectName, rxn.AuthorId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 生成pageName
	pageName := fmt.Sprintf("%s-%03d", rxn.ProjectName, rankId+1)

	// 插入eln_project_page表
	if err = InsertElnProjectPage(tx, rxn.ReactionId, rxn.AuthorId, projectId, pageName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 更新eln_project表的rank_id
	if err = UpdateProjectRankId(tx, rxn.ProjectName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 插入eln_commit_log表
	if err = InsertElnCommitLog(tx, rxn.ReactionId, rxn.ProjectName, rxn.AuthorId, rxn.AuthorName); err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"pageName": pageName, "reactionId": rxn.ReactionId})
	return
}

// InsertElnRxnBasicInfo 辅助函数：插入eln_rxn_basicinfo表
func InsertElnRxnBasicInfo(tx *sqlx.Tx, rxn common.ElnRxnBasicInfo) error {
	_, err := tx.NamedExec(`insert into eln_rxn_basicinfo(reaction_id, project_name, author_id, author_name, creation_date, start_date) VALUE (:reaction_id, :project_name, :author_id, :author_name, CURRENT_DATE, CURRENT_DATE)`, rxn)
	return err
}

// GetOrCreateProjectIdAndRankId 辅助函数：查询或插入eln_project表，并获取projectId和rankId
func GetOrCreateProjectIdAndRankId(tx *sqlx.Tx, projectName string, userId int64) (int64, int64, error) {
	type rpr struct {
		ProjectId int64 `db:"project_id"`
		RankId    int64 `db:"rank_id"`
	}
	var rxnPidRid rpr

	err := tx.Get(&rxnPidRid, `select project_id, rank_id from eln_project where project_name = binary ?`, projectName)

	if errors.Is(err, sql.ErrNoRows) {
		rxnPidRid = rpr{
			ProjectId: utils.GenerateSnowflakeID(),
			RankId:    0,
		}
		_, err = tx.Exec(`insert into eln_project(project_id, project_name) VALUE (?, ?)`, rxnPidRid.ProjectId, projectName)
		if err != nil {
			return 0, 0, err
		}
	} else if err != nil {
		return 0, 0, err
	}

	return rxnPidRid.ProjectId, rxnPidRid.RankId, nil
}

// InsertElnProjectPage 辅助函数：插入eln_project_page表
func InsertElnProjectPage(tx *sqlx.Tx, reactionId, userId, projectId int64, pageName string) error {
	_, err := tx.Exec(`insert into eln_project_page(reaction_id, user_id, project_id, page_name, time_creation) VALUE (?, ?, ?, ?, current_timestamp) `, reactionId, userId, projectId, pageName)
	return err
}

// UpdateProjectRankId 辅助函数：更新eln_project表的rank_id
func UpdateProjectRankId(tx *sqlx.Tx, projectName string) error {
	_, err := tx.Exec(`update eln_project set rank_id = rank_id+1 where project_name = binary ?`, projectName)
	return err
}

// InsertElnCommitLog 辅助函数：插入eln_commit_log表
func InsertElnCommitLog(tx *sqlx.Tx, reactionId int64, project string, userId int64, userName string) error {
	_, err := tx.Exec(`insert into eln_commit_log(reaction_id, project, user_id, user_name, status ) VALUE (?,?,?,?,?)`, reactionId, project, userId, userName, "create")
	return err
}
