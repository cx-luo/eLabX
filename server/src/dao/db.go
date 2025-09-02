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
	"time"

	localMysql "github.com/go-sql-driver/mysql"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var OBCursor *gorm.DB
var MinioClient *minio.Client
var RedisClient *redis.Client

// GetMysqlCursor initializes and returns a MySQL gorm.DB connection with optimized settings.
func GetMysqlCursor(host string, port int, username, passwd, dbname string) {
	conf := localMysql.Config{
		User:                 username,
		Passwd:               passwd,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", host, port),
		DBName:               dbname,
		Timeout:              30 * time.Second,
		ReadTimeout:          10 * time.Second,
		WriteTimeout:         5 * time.Second,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	dsn := conf.FormatDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		AllowGlobalUpdate:      false,
		QueryFields:            true,
		DisableAutomaticPing:   false,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open MySQL connection: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get sql.DB from gorm.DB: %v", err))
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping MySQL: %v", err))
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	OBCursor = db
}

// GetRedisClient initializes and returns a Redis client with optimized settings.
func GetRedisClient(addr, password string, db int) {
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		PoolSize:        30,
		MinIdleConns:    5,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 10 * time.Minute,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
	RedisClient = client
}

// ConnectToMinio initializes and returns a MinIO client.
func ConnectToMinio(endpoint, accessKeyID, secretAccessKey string, useSSL bool) {
	// minio "github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7/pkg/credentials"
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MinIO: %v", err))
	}
	MinioClient = client
}
