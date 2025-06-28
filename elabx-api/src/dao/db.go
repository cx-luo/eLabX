// Package dao coding=utf-8
// @Project : eLabX
// @Time    : 2023/10/23 11:33
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@pharmaron.com
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
	db, err := gorm.Open(mysql.Open(conf.FormatDSN()), &gorm.Config{})

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
		//Addrs:           []string{"192.168.2.139:4000"},
		Addr:            "192.168.2.139:4000",
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

var RedisClient = GetRedisClusterClient()
