// Package api coding=utf-8
// @Project : server
// @Time    : 2024/9/29 10:21
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : project_manager.go
// @Software: GoLand
package api

import (
	"database/sql"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetProject(c *gin.Context) {
	var p struct {
		UserId      int    `json:"userId"`
		ProjectName string `json:"projectName,omitempty"`
	}
	err := c.ShouldBind(&p)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	query := `select b.user_id, b.project_id, a.project_name, b.page_name, b.reaction_id , a.rank_id 
		from eln_project a join eln_project_page b on a.project_id = b.project_id 
		where b.user_id = ?`
	if p.ProjectName != "" {
		query = query + fmt.Sprintf(" and project_name = '%s'", p.ProjectName)
	}
	var projects []ElnProject
	err = dao.OBCursor.Select(&projects, query, p.UserId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"total": len(projects), "projects": projects})
	return
}

func DeleteProject(c *gin.Context) {
	var p struct {
		UserId    int   `json:"userId"`
		ProjectId int64 `json:"projectId"`
	}
	err := c.ShouldBind(&p)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
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

	_, err = tx.Exec(`replace into elabx_deleted.eln_project_bak select project_id, project_name, user_id, rank_id from eln_project where project_id = ?`, p.ProjectId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	_, err = tx.Exec(`replace into elabx_deleted.eln_project_page_bak select reaction_id, user_id, project_id, page_name, time_creation, export_cnt from eln_project_page where project_id = ?`, p.ProjectId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	var rxnIds []int64
	err = tx.Select(&rxnIds, `select reaction_id from eln_project_page where user_id = ? and project_id = ?`, p.UserId, p.ProjectId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	if len(rxnIds) > 0 {
		for _, rxnId := range rxnIds {
			sema.Acquire(1)
			go func(rId int64) {
				defer sema.Release()
				if utils.GlobalConfig.Service.Backup {
					err = BackupReactionBeforeDelete(tx, rId)
					if err != nil {
						utils.InternalRequestErr(c, err)
						return
					}
				}
				err = DeleteReaction(tx, rId)
				if err != nil {
					utils.InternalRequestErr(c, err)
					return
				}
				_, err = tx.Exec(`delete from eln_rxn_basicinfo where reaction_id = ?`, rId)
				if err != nil {
					return
				}
			}(rxnId)
		}
		sema.Wait()
	}

	_, err = tx.Exec(`delete from eln_project_page where project_id = ?`, p.ProjectId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	_, err = tx.Exec(`delete from eln_project where project_id = ?`, p.ProjectId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.SuccessWithData(c, "", gin.H{"total": len(rxnIds)})
	return
}

func SearchProject(c *gin.Context) {
	var condition struct {
		UserId      int    `json:"userId,omitempty"`
		ProjectName string `json:"projectName,omitempty"`
		PageName    string `json:"pageName,omitempty"`
		ReactionId  uint64 `json:"reactionId,omitempty"`
	}
	err := c.ShouldBind(&condition)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	type prjInfo struct {
		UserId      int    `json:"userId,omitempty" db:"user_id"`
		ProjectId   int64  `json:"projectId,omitempty" db:"project_id"`
		ProjectName string `json:"projectName,omitempty" db:"project_name"`
		PageName    string `json:"pageName,omitempty" db:"page_name"`
		ReactionId  int64  `json:"reactionId,omitempty" db:"reaction_id"`
		ReactionImg string `json:"reactionImg,omitempty"`
		RxnStatus   string `json:"rxnStatus,omitempty"`
		ExportCnt   uint64 `json:"exportCnt" db:"export_cnt"`
	}
	var projectInfo []prjInfo
	baseSql := "select b.user_id, a.project_id, a.project_name, b.page_name ,b.reaction_id, b.export_cnt from eln_project a left join eln_project_page b on a.project_id = b.project_id where %s order by a.user_id, b.page_name "
	baseConditions := "1 = 1"
	if condition.UserId != 0 {
		baseConditions += fmt.Sprintf(" and a.user_id = %d", condition.UserId)
	}
	if condition.ProjectName != "" {
		baseConditions += fmt.Sprintf(" and a.project_name = binary '%s'", condition.ProjectName)
	}
	if condition.PageName != "" {
		baseConditions += fmt.Sprintf(" and b.page_name = binary '%s'", condition.PageName)
	}
	if condition.ReactionId != 0 {
		baseConditions += fmt.Sprintf(" and b.reaction_id = %d", condition.ReactionId)
	}

	err = dao.OBCursor.Select(&projectInfo, fmt.Sprintf(baseSql, baseConditions))
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	for idx := range projectInfo {
		sema.Acquire(1)
		go func(i int) {
			defer sema.Release()
			var smile struct {
				CdStructure string `db:"cd_structure"`
				RxnStatus   string `db:"rxn_status"`
			}
			err := dao.OBCursor.Get(&smile, `select a.cd_structure, b.rxn_status from eln_reaction_note a join eln_rxn_basicinfo b on a.reaction_id = b.reaction_id where a.reaction_id = ?`, projectInfo[i].ReactionId)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				utils.InternalRequestErr(c, errors.New(err.Error()+": "+strconv.FormatInt(projectInfo[i].ProjectId, 10)))
				return
			}
			err, projectInfo[i].ReactionImg = GenImgBase64String(smile.CdStructure)
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
			projectInfo[i].RxnStatus = smile.RxnStatus
		}(idx)
	}
	sema.Wait()

	utils.SuccessWithData(c, "", gin.H{"projectInfo": projectInfo, "total": len(projectInfo)})
	return
}

func DeletePage(c *gin.Context) {
	var p struct {
		UserId     int     `json:"userId"`
		ReactionId []int64 `json:"reactionId" db:"reaction_id"`
	}
	err := c.ShouldBind(&p)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
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

	for _, rxnId := range p.ReactionId {
		sema.Acquire(1)
		go func(rId int64) {
			defer sema.Release()

			_, err = tx.Exec(`replace into elabx_deleted.eln_project_page_bak select reaction_id, user_id, project_id, page_name, time_creation, export_cnt from eln_project_page where reaction_id = ?`, rId)
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
			_, err = tx.Exec(`delete from eln_project_page where reaction_id = ?`, rId)
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}

			if utils.GlobalConfig.Service.Backup {
				err = BackupReactionBeforeDelete(tx, rId)
				if err != nil {
					utils.InternalRequestErr(c, err)
					return
				}
			}
			err = DeleteReaction(tx, rId)
			if err != nil {
				utils.InternalRequestErr(c, err)
				return
			}
			_, err = tx.Exec(`delete from eln_rxn_basicinfo where reaction_id = ?`, rId)
			if err != nil {
				return
			}
		}(rxnId)
	}
	sema.Wait()

	utils.SuccessWithData(c, "", gin.H{"total": len(p.ReactionId)})
	return
}
