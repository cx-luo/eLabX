// Package middleware coding=utf-8
// @Project : elabx-api
// @Time    : 2025/6/28 13:33
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : casbin_middleware.go
// @Software: GoLand
package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user") // 假设你通过 JWT 或 Session 获取了用户信息
		roles := user.(map[string]interface{})["roles"].([]string)

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
