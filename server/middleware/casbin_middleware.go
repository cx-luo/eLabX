// Package middleware coding=utf-8
// @Project : eLabX
// @Time    : 2025/6/28 13:33
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : casbin_middleware.go
// @Software: GoLand
package middleware

import (
	"database/sql"
	"eLabX/src/dao"
	"eLabX/src/types"
	"eLabX/src/utils"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var GlobalCasBin *casbin.Enforcer

func setupCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.act, p.act) && keyMatch(r.obj, p.obj)
`)
	a, err := gormadapter.NewAdapterByDBUseTableName(db, "eln_", "policies")
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}

	GlobalCasBin = e
	return e, nil
}

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, s := range wightList {
			if c.Request.URL.Path == s {
				c.Next()
				return
			}
		}
		//roles := user.(map[string]interface{})["roles"].([]string)
		roles, err := getUserPermissions(c)
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not find permission."})
			return
		}
		enforcer, err := setupCasbin(dao.OBCursor)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("setupCasbin error: %s", err.Error())})
			return
		}
		for _, role := range roles {
			allowed, err := enforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "The permission validation failed."})
				return
			}
			if allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No access permission."})
	}
}

func getUserPermissions(c *gin.Context) ([]string, error) {
	username, exists := c.Get("username")
	if !exists {
		return nil, errors.New("User does not exist or is unavailable.\n")
	}
	var user types.ElnUsers
	err := dao.OBCursor.Where("status = 1 and user_id = ?", username).Find(&user).Error
	if err != nil {
		utils.NotFoundError(c, fmt.Errorf("User does not exist or is unavailable. %s\n", err))
		return nil, fmt.Errorf("User does not exist or is unavailable. %s\n", err)
	}

	// Get group IDs from user.GroupId (assuming it's a string of comma-separated IDs)
	var userGroups []types.ElnUserGroup
	err = dao.OBCursor.Model(&types.ElnUserGroup{}).Where("user_id = ? AND status = 1", user.UserId).Find(&userGroups).Error
	if err != nil {
		return nil, err
	}
	var groupIds []string
	for _, ug := range userGroups {
		groupIds = append(groupIds, fmt.Sprintf("%d", ug.GroupId))
	}

	// Query permissions for these group IDs
	var permissions []string
	var permissionIds []int64
	if len(groupIds) > 0 && groupIds[0] != "" {
		// First, get permissions string(s) for the user's groups
		var groupPermissions []string
		err := dao.OBCursor.Model(&types.ElnGroup{}).Select("permissions").Where("group_id IN (?)", groupIds).Pluck("permissions", &groupPermissions).Error
		if err != nil {
			return nil, err
		}
		// Collect all permission IDs from all groups
		permissionIdSet := make(map[int64]struct{})
		for _, perms := range groupPermissions {
			for _, pidStr := range strings.Split(perms, ",") {
				pidStr = strings.TrimSpace(pidStr)
				if pidStr == "" {
					continue
				}
				var pid int64
				_, err := fmt.Sscan(pidStr, &pid)
				if err == nil {
					permissionIdSet[pid] = struct{}{}
				}
			}
		}
		permissionIds = make([]int64, 0, len(permissionIdSet))
		for pid := range permissionIdSet {
			permissionIds = append(permissionIds, pid)
		}
		// Now, get permission names
		if len(permissionIds) > 0 {
			err = dao.OBCursor.Model(&types.ElnPermission{}).Where("permission_id in (?)", permissionIds).Pluck("permission_name", &permissions).Error
			if err != nil {
				return nil, err
			}
		}
	}

	return permissions, nil
}
