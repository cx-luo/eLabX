// Package middleware coding=utf-8
// @Project : eLabX
// @Time    : 2025/8/4 10:07
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : idempotency.go
// @Software: GoLand
package middleware

// Idempotency middleware for HTTP requests.
// Ensures that requests with the same idempotency key are only processed once.

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory store for idempotency keys and their responses.
// In production, use a persistent store like Redis.
// Optimized idempotency middleware for Gin

type idempotencyEntry struct {
	response   []byte
	statusCode int
	header     http.Header
	timestamp  time.Time
}

var (
	idempotencyStore = make(map[string]*idempotencyEntry)
	idempotencyMutex sync.Mutex
	idempotencyTTL   = 10 * time.Minute // configurable TTL for idempotency keys
)

// responseRecorder captures the response for storage.
type responseRecorder struct {
	gin.ResponseWriter
	body        []byte
	statusCode  int
	wroteHeader bool
}

func (r *responseRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.statusCode = code
		r.ResponseWriter.WriteHeader(code)
		r.wroteHeader = true
	}
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

// GinIdempotencyMiddleware is a Gin middleware that enforces idempotency for POST, PUT, PATCH requests with an Idempotency-Key header.
func GinIdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
			c.Next()
			return
		}

		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.Next()
			return
		}

		// Clean up expired entries and check for existing entry
		var entry *idempotencyEntry
		var exists bool
		now := time.Now()
		idempotencyMutex.Lock()
		for k, v := range idempotencyStore {
			if now.Sub(v.timestamp) > idempotencyTTL {
				delete(idempotencyStore, k)
			}
		}
		entry, exists = idempotencyStore[key]
		idempotencyMutex.Unlock()

		if exists && now.Sub(entry.timestamp) <= idempotencyTTL {
			// Replay the stored response
			for k, vv := range entry.header {
				for _, v := range vv {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.WriteHeader(entry.statusCode)
			_, _ = c.Writer.Write(entry.response)
			c.Abort()
			return
		}

		// Record the response
		rec := &responseRecorder{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
		}
		c.Writer = rec

		c.Next()

		// Store the response after handler
		idempotencyMutex.Lock()
		idempotencyStore[key] = &idempotencyEntry{
			response:   rec.body,
			statusCode: rec.statusCode,
			header:     c.Writer.Header().Clone(),
			timestamp:  time.Now(),
		}
		idempotencyMutex.Unlock()
	}
}
