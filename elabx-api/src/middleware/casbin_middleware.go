// Package middleware coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/28 13:33
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : casbin_middleware.go
// @Software: GoLand
package middleware

import (
	"database/sql"
	"eLabX/src/dao"
	"errors"
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
		user, _ := c.Get("username") // 假设你通过 JWT 或 Session 获取了用户信息
		//roles := user.(map[string]interface{})["roles"].([]string)
		var permissions string
		err := dao.OBCursor.Table("eln_users").Select("permissions").Where("user_id = ?", user).Find(&permissions).Error
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not find permission."})
			return
		}
		roles := strings.Split(permissions, ",")
		enforcer, err := setupCasbin(dao.OBCursor)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
