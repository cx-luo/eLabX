// Package dao coding=utf-8
// @Project : eLabX
// @Time    : 2023/10/23 11:33
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : db.go
// @Software: GoLand
package dao

import (
	"context"
	"fmt"
	localMysql "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func GetMysqlCursor(host string, port int, username string, passwd string, dbname string) *gorm.DB {
	conf := localMysql.Config{
		User:                 username,
		Passwd:               passwd,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", host, port),
		DBName:               dbname,
		Timeout:              time.Second * 300,
		ReadTimeout:          time.Second * 10,
		WriteTimeout:         time.Second * 30,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := gorm.Open(mysql.Open(conf.FormatDSN()), &gorm.Config{
		PrepareStmt:            true, // 可以改为 false，全局禁用预编译
		SkipDefaultTransaction: true,
		AllowGlobalUpdate:      false, // 添加：禁止全局更新
		QueryFields:            true,  // 添加：显式指定查询字段
		DisableAutomaticPing:   false, // 启用自动 ping，及时发现连接问题
	})

	gormMysql, _ := db.DB()

	//尝试与数据库进行连接
	err = gormMysql.Ping()
	if err != nil {
		fmt.Println("数据库连接失败", err)
		panic("Ping db failed!")
	}
	gormMysql.SetMaxOpenConns(30)
	gormMysql.SetMaxIdleConns(4)
	return db
}

var OBCursor *gorm.DB

func GetRedisClusterClient() *redis.Client {
	redisDb := redis.NewClient(&redis.Options{
		Addr:            "127.0.0.1：4000",
		Password:        "",
		PoolSize:        20,
		MaxIdleConns:    4,
		ConnMaxIdleTime: time.Second * 600,
	})

	ctx := context.Background()
	_, err := redisDb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return redisDb
}

//var RedisClient = GetRedisClusterClient()
