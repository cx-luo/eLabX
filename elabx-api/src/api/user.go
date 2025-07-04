// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/8 10:15
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : tables.go
// @Software: GoLand
package api

import (
	"database/sql"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"errors"
	"fmt"
	goTools "github.com/cx-luo/go-toolkit"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

var sema = goTools.NewSemaphore(32)

type Users struct {
	Username string `json:"username" db:"username" gorm:"username"`
	Roles    string `json:"roles" db:"roles" gorm:"roles"`
	//Avatar      string         `json:"avatar" db:"avatar" gorm:"avatar"`
	Permissions sql.NullString `json:"permissions" db:"permissions" gorm:"permissions"`
}

type ElnUsers struct {
	ID           int64     `json:"id,omitempty" db:"id" gorm:"id"`
	CreatedAt    time.Time `json:"createdAt,omitempty" db:"created_at" gorm:"created_at"`
	Userid       int64     `json:"userid,omitempty" db:"userid" gorm:"userid"`
	Username     string    `json:"username,omitempty" db:"username" gorm:"username"`
	Email        string    `json:"email,omitempty" db:"email" gorm:"email"`
	Phone        string    `json:"phone,omitempty" db:"phone" gorm:"phone"`
	PasswordHash string    `json:"passwordHash,omitempty" db:"password_hash" gorm:"password_hash"`
	Status       int64     `json:"status,omitempty" db:"status" gorm:"status"`
	GroupId      int64     `json:"groupId,omitempty" db:"group_id" gorm:"group_id"`
	Permissions  string    `json:"permissions,omitempty" db:"permissions" gorm:"permissions"`
	Authn        string    `json:"authn,omitempty" db:"authn" gorm:"authn"`
	AuthorityIds []string  `json:"authorityIds,omitempty"`
}

func GetUserList(c *gin.Context) {
	var users []ElnUsers
	err := dao.OBCursor.Select(`select userid, username, phone, email, permissions, status, created_at from eln_users order by userid`).Scan(&users).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	// 采用 range 获取的下标值，然后用下标方式引用的数组项也可以直接修改
	// 采用 range 获取数组项不能修改数组中结构体的值
	for i, user := range users {
		users[i].AuthorityIds = strings.Split(user.Permissions, ",")
	}

	utils.SuccessWithData(c, "", gin.H{"total": len(users), "users": users})
	return
}

type route struct {
	BookName   string  `json:"bookName"`
	Icon       string  `json:"icon"`
	Path       string  `json:"path"`
	Layout     string  `json:"layout,omitempty"`
	PageName   []route `json:"pageName,omitempty"`
	ReactionId int64   `json:"reactionId,omitempty"`
}

func UserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.BadRequestErr(c, errors.New("User does not exist or is unavailable.\n"))
		return
	}

	var user Users
	err := dao.OBCursor.Table("eln_users").Select("username", "roles", "permissions").
		Where("status = 1 and user_id = ?", username).Find(&user).Error
	if err != nil {
		utils.NotFoundError(c, fmt.Errorf("User does not exist or is unavailable. %s\n", err))
		return
	}

	utils.SuccessWithData(c, "success", gin.H{"permissions": strings.Split(user.Permissions.String, ","), "username": user.Username})
	return
}

type userId struct {
	WitnessId int `json:"witnessId"`
}

func FetchUserName(c *gin.Context) {
	var uid userId
	err := c.ShouldBind(&uid)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	type userInfo struct {
		ChemistId   int    `json:"chemistId" db:"chemist_id"`
		ChemistName string `json:"chemistName" db:"chemist_name"`
	}
	var u userInfo
	err = dao.OBCursor.Table("eln_company_user_list").Where(`chemist_id = ?`, uid.WitnessId).
		Find(&u).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SuccessWithData(c, "", gin.H{"user": userInfo{uid.WitnessId, strconv.Itoa(uid.WitnessId)}})
		} else {
			utils.InternalRequestErr(c, err)
		}
		return
	}

	utils.SuccessWithData(c, "", gin.H{"user": u})
	return
}

type pwdForm struct {
	UserID      int    `json:"userId"`
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

func ChangePwd(c *gin.Context) {
	var p pwdForm
	err := c.ShouldBind(&p)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var oldPwd string
	err = dao.OBCursor.Table("eln_users").Select("password_hash").
		Where("userid = ?", p.UserID).Find(&oldPwd).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	if oldPwd != p.OldPassword {
		utils.BadRequestErr(c, errors.New("the old password was incorrectly entered"))
		return
	}
	err = dao.OBCursor.Table("eln_users").Where("userid = ?", p.UserID).
		Update("password_hash", p.Password).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Password changed")
	return
}

type ForgetPwdForm struct {
	UserID   string `json:"userId"`
	SmsCode  string `json:"smsCode"`
	Password string `json:"password"`
}

func ForgetPwd(c *gin.Context) {
	now := time.Now().Unix()
	var p ForgetPwdForm
	err := c.ShouldBind(&p)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var s struct {
		VerifCode           string `json:"verifCode" db:"verif_code"`
		VerifCodeExpireTime int64  `json:"verifCodeExpireTime" db:"verif_code_expire_time"`
	}
	err = dao.OBCursor.Table("eln_register_records").Select("verif_code", "verif_code_expire_time").
		Where(`userid = ?`, p.UserID).Find(&s).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	duration := s.VerifCodeExpireTime - now
	if duration < 125 && duration > 0 {
		err = dao.OBCursor.Table("eln_users").
			Update("password_hash", p.Password).
			Where(`userid = ?`, p.UserID).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}
	} else {
		utils.BadRequestErr(c, errors.New("timeout"))
		return
	}
	utils.Success(c, "change passcode successfully")
}

func ModifyUserInfo(c *gin.Context) {
	var u struct {
		UserId   int    `json:"userId"`
		Username string `json:"username"`
	}
	err := c.ShouldBind(&u)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = dao.OBCursor.Table("eln_users").
		Update("username", u.Username).
		Where(`userid = ?`, u.UserId).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Username changed")
	return
}

func ResetPwd(c *gin.Context) {
	var s struct {
		RequestId int `json:"requestId"`
		UserID    int `json:"userId"`
	}

	err := c.ShouldBind(&s)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	var permission string
	err = dao.OBCursor.Table("eln_users").Select("permissions").Where(`userid = ?`, s.RequestId).Find(&permission).Error
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if strings.Contains(permission, "admin") {
		err = dao.OBCursor.Table("eln_users").
			Update("password_hash", "e10adc3949ba59abbe56e057f20f883e").
			Where(`userid = ?`, s.RequestId).Error
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}

		utils.Success(c, "change passcode successfully")
		return
	}

	utils.InternalRequestErr(c, errors.New("no permission"))
}
