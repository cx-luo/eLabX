// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2025/7/18 16:09
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : emails.go
// @Software: GoLand
package utils

import (
	"errors"
	"net/smtp"
	"strings"
)

// SendEmail sends an email
func SendEmail(to, subject, body string) error {
	// Please modify the following configuration according to your actual situation
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	smtpUser := "your_email@example.com"
	smtpPass := "your_email_password"

	from := smtpUser
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// 支持多个收件人
	toList := strings.Split(to, ",")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, []byte(msg))
	if err != nil {
		return errors.New("failed to send email: " + err.Error())
	}
	return nil
}
