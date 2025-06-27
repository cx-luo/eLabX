// Package middleware coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/29 10:18
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : cas.go
// @Software: GoLand
package middleware

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/cas.v2"
	"net/http"
	"net/url"
)

func CasAuthenticated(c *gin.Context) {
	u, _ := url.Parse("https://cas.pharmaron-bj.com/cas")
	client := cas.NewClient(&cas.Options{URL: u})

	h := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cas.IsAuthenticated(r) {
			client.RedirectToLogin(w, r)
		} else {
			c.String(200, "已登录")
		}
	})
	h.ServeHTTP(c.Writer, c.Request)
	c.String(200, "登录成功")
}
