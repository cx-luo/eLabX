// Package middleware coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/8 10:23
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : jwt_middleware.go
// @Software: GoLand
package middleware

import (
	utils2 "eLabX/src/utils"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

const TokenExpireDuration = time.Hour * 36

var MySecret = []byte("zH8jim3nbCU96ZGb")

// GenToken 生成JWT
func GenToken(username string, password string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		Username: username, // 自定义字段
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    "pharmaron-ai",                                          // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return tokenStruct.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	jwtTokenStr := strings.Split(tokenString, " ")
	if len(jwtTokenStr) < 2 {
		return nil, errors.New("invalid token")
	}
	jwtToken := jwtTokenStr[1]
	// 解析token
	token, err := jwt.ParseWithClaims(jwtToken, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 白名单,列表里的路径不进行jwt验证
var wightList = []string{
	"/api/user/login",
	"/api/user/register",
	"/api/user/logout",
	"/api/user/genVerifCode",
	"/api/user/forgetPwd",
	"/api/user/sentCodeOnly",
}

// JwtAuth jwt中间件
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, s := range wightList {
			if c.Request.URL.Path == s {
				c.Next()
				return
			}
		}
		tokenHeader := c.Request.Header.Get("Authorization")
		//
		//tokenHeader := c.Request.Header.Get("accessToken")
		if tokenHeader == "" {
			utils2.BadRequestErr(c, errors.New("token not exist"))
			return
		}

		claims, err := ParseToken(tokenHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils2.BaseResponse{
				StatusCode: -2, Msg: "inactivity timeout, please close browser and re-login", Data: gin.H{},
			})
			c.Abort()
			return
		}

		//判断token是否过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.JSON(http.StatusBadRequest, utils2.BaseResponse{
				StatusCode: -2, Msg: "inactivity timeout, please close browser and re-login", Data: gin.H{},
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // 允许所有来源，或者使用config.AllowOrigins = []string{"https://yourdomain.com"}来指定具体来源
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}

	// 如果需要支持凭证（cookies, HTTP authentication...）
	// 请注意：当AllowCredentials为true时，不能使用AllowAllOrigins: true
	// config.AllowCredentials = true

	// 如果需要限制特定的headers
	// config.AllowHeaders = []string{"Authorization"}

	// 设置最大预检存活时间
	config.MaxAge = 12 * time.Hour
	return cors.New(config)
}
