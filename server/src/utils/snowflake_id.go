// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2024/4/20 16:33
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : snowflake_id.go
// @Software: GoLand
package utils

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	idBitSize      = 10                // ID总位数
	seqBitSize     = 6                 // 序列号位数
	timestampShift = seqBitSize        // 时间戳左移位数
	maxSequence    = 1<<seqBitSize - 1 // 最大序列号
)

var (
	lastTimestamp int64  = -1
	sequence      uint32 = 0
)

func GenerateSnowflakeID() int64 {
	now := time.Now().UnixNano() / 1e6 // 使用毫秒级时间戳

	if now < lastTimestamp { // 避免时间回拨
		time.Sleep(time.Millisecond)
		now = time.Now().UnixNano() / 1e6
	}

	if now == lastTimestamp {
		// 同一毫秒内自增序列号
		nextSequence := atomic.AddUint32(&sequence, 1)
		if nextSequence > maxSequence {
			time.Sleep(time.Millisecond)
			now = time.Now().UnixNano() / 1e6
			atomic.StoreUint32(&sequence, 0)
		}
	} else {
		atomic.StoreUint32(&sequence, 0)
	}

	lastTimestamp = now

	// 组合时间戳和序列号
	return now<<timestampShift | int64(sequence)
}

// GenRandInt generates a random int64 with the specified number of digits.
// If n <= 0, returns 0. If n > 18, n is capped at 18 (max int64 digits).
// GenRandInt generates a unique int64 using a snowflake-like rule, truncated to n digits if needed.
func GenRandInt(n int) int64 {
	id := GenerateSnowflakeID()
	if n <= 0 {
		return 0
	}
	if n > 18 {
		n = 18
	}
	// Truncate or pad the ID to n digits
	idStr := []rune(fmt.Sprintf("%018d", id)) // pad to 18 digits
	start := 18 - n
	if start < 0 {
		start = 0
	}
	resultStr := string(idStr[start:])
	result, err := strconv.ParseInt(resultStr, 10, 64)
	if err != nil {
		return id // fallback to original id if conversion fails
	}
	return result
}
