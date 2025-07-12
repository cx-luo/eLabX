// Package dao coding=utf-8
// @Project : eLabX
// @Time    : 2024/2/23 10:54
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : redis_ops.go
// @Software: GoLand
package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//type predictResList []common.PredictRes
//
//var _ encoding.BinaryMarshaler = new(predictResList)
//var _ encoding.BinaryUnmarshaler = new(predictResList)
//
//func (m *predictResList) MarshalBinary() (data []byte, err error) {
//	return json.Marshal(m)
//}
//
//func (m *predictResList) UnmarshalBinary(data []byte) error {
//	return json.Unmarshal(data, m)
//}

var ctx = context.Background()

// 定义自定义错误类型
var (
	ErrKeyNotExist = errors.New("key not exist")
	ErrValueIsNull = errors.New("value is nullString")
)

func GetRSearchResult(key string) (error, string) {
	val, err := RedisClient.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return ErrKeyNotExist, ""
	case err != nil:
		return err, ""
	case val == "":
		return ErrValueIsNull, ""
	}
	return nil, val
}

func WriteSearchResult(key string, val interface{}) error {
	// 序列化结构体为JSON
	//middleware.Logger.Info("Processing and writing to Redis", zap.String("searchKey", key))
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		fmt.Println(duration)
		//middleware.Logger.Info("Processing completed", zap.Duration("duration", duration))
	}()

	jsonData, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = RedisClient.Set(ctx, key, string(jsonData), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

type RedisKeyData struct {
	TotalCount int         `json:"totalCount"`
	Rows       interface{} `json:"rows"`
}

func CleanUpKeys(rdb *redis.Client) error {
	var cursor uint64
	pattern := "*" // 根据实际情况调整模式匹配

	for {
		keys, newCursor, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("error scanning keys: %v", err)
		}

		for _, key := range keys {
			value, err := rdb.Get(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				continue // 键不存在，跳过
			} else if err != nil {
				return fmt.Errorf("error getting value for key %s: %v", key, err)
			}

			var data RedisKeyData
			if err := json.Unmarshal([]byte(value), &data); err != nil {
				return fmt.Errorf("error unmarshalling JSON for key %s: %v", key, err)
			}

			if data.TotalCount == 0 {
				if err := rdb.Del(ctx, key).Err(); err != nil {
					return fmt.Errorf("error deleting key %s: %v", key, err)
				}
				//middleware.Logger.Info(fmt.Sprintf("Deleted key: %s\n", key))
			}
		}

		if newCursor == 0 {
			break // 扫描完成
		}
		cursor = newCursor
	}
	return nil
}
