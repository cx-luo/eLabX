// Package middleware coding=utf-8
// @Project : server
// @Time    : 2024/10/22 10:27
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : encrypt.go
// @Software: GoLand
package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// 定义一个自定义的ResponseWriter，用于捕获响应数据
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 方法实现写入数据到body Buffer
func (w *responseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

// WriteHeader 方法设置响应状态码
func (w *responseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
}

// 加密密钥，实际应用中应该是一个安全的密钥
var secretKey = []byte("68f5e589b5646a1a8edb5001b87780b7")
var ivKey = []byte("11f503747061573db62ce666332f3266")

func EncryptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解密请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Request.Body.Close()
		plaintext, err := Decrypt(secretKey, string(body))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fmt.Println(plaintext)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(plaintext))

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           new(bytes.Buffer),
		}
		c.Writer = rw
		// 处理请求
		c.Next()

		// 加密响应体
		encodedResponse, err := Encrypt(secretKey, rw.body.Bytes(), ivKey)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Writer.WriteHeader(c.Writer.Status())
		c.Writer.Write([]byte(encodedResponse))
	}
}

// Encrypt 加密数据
func Encrypt(key []byte, plaintext []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// 将密文和 IV 一起编码，以便解密时使用
	encoded := make([]byte, len(iv)+len(ciphertext))
	copy(encoded, iv)
	copy(encoded[len(iv):], ciphertext)

	return base64.StdEncoding.EncodeToString(encoded), nil
}

// Decrypt 解密数据
func Decrypt(key []byte, ciphertext string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	if len(decoded) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := decoded[:aes.BlockSize]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(decoded)-aes.BlockSize)
	stream.XORKeyStream(plaintext, decoded[aes.BlockSize:])
	fmt.Println(plaintext)

	return plaintext, nil
}
