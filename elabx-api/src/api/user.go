// Package api coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/8 10:15
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : user.go
// @Software: GoLand
package api

import (
	"database/sql"
	"eLabX/src/dao"
	"eLabX/src/middleware"
	"eLabX/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	goTools "github.com/cx-luo/go-toolkit"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sema = goTools.NewSemaphore(32)

type CzUsers struct {
	ID           int64     `json:"id" db:"id" gorm:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at" db:"deleted_at" gorm:"deleted_at"`
	Email        string    `json:"email" db:"email" gorm:"email"`
	Username     string    `json:"username" db:"username" gorm:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash" gorm:"password_hash"`
	Status       int64     `json:"status" db:"status" gorm:"status"`
	GroupId      int64     `json:"group_id" db:"group_id" gorm:"group_id"`
	Storage      int64     `json:"storage" db:"storage" gorm:"storage"`
	TwoFactor    string    `json:"two_factor" db:"two_factor" gorm:"two_factor"`
	//Avatar       string    `json:"avatar" db:"avatar" gorm:"avatar"`
	Permissions string `json:"permissions" db:"permissions" gorm:"permissions"`
	Authn       string `json:"authn" db:"authn" gorm:"authn"`
}

type Users struct {
	Email    string `json:"email" db:"email" gorm:"email"`
	Username string `json:"username" db:"username" gorm:"username"`
	GroupId  int64  `json:"group_id" db:"group_id" gorm:"group_id"`
	//Avatar      string         `json:"avatar" db:"avatar" gorm:"avatar"`
	Permissions sql.NullString `json:"permissions" db:"permissions" gorm:"permissions"`
}

func UserLogin(c *gin.Context) {
	type user struct {
		UserId   string `json:"userid"`
		Password string `json:"password"`
	}
	var u user
	err := c.ShouldBind(&u)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}

	var passwordHash string
	err = dao.OBCursor.Get(
		&passwordHash,
		`SELECT password_hash from eln_users WHERE status = 1 AND userid = ?`,
		u.UserId,
	)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"msg": "User authentication failed, username or password incorrect."})
		return
	}

	if passwordHash == u.Password {
		token, _ := middleware.GenToken(u.UserId, u.Password)
		response := utils.BaseResponse{
			StatusCode: 200,
			Msg:        "success",
			Data:       gin.H{"accessToken": "Bearer " + token},
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User authentication failed, username or password incorrect."})
	}
}

func UserLogout(c *gin.Context) {
	utils.Success(c, "Logged out successfully")
	return
}

func UserInfo(c *gin.Context) {
	var user Users
	userid, ok := c.Get("username")
	if !ok {
		response := utils.BaseResponse{
			StatusCode: 400, Msg: "Please log in again.", Data: gin.H{"error": "Please log in again."},
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//
	err := dao.OBCursor.Get(
		&user,
		`SELECT username, email, group_id, permissions from eln_users WHERE status = 1 and userid = ?`, userid,
	)
	if err != nil {
		u := fmt.Sprintf("User does not exist or is unavailable.")
		response := utils.BaseResponse{
			StatusCode: 505, Msg: err.Error(), Data: gin.H{"error": u},
		}
		c.JSON(505, response)
		return
	}

	err, elnRoutes, authorityRoutes := userRouter(utils.GetInterfaceToInt(userid))
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	response := utils.BaseResponse{
		StatusCode: 200, Msg: "success", Data: gin.H{"permissions": strings.Split(user.Permissions.String, ","), "username": user.Username,
			"routes": elnRoutes, "authorityRoutes": authorityRoutes},
	}
	c.JSON(http.StatusOK, response)
	return
}

func UserRegister(c *gin.Context) {
	type userForm struct {
		Data struct {
			UserId          string `json:"userid"`
			Name            string `json:"name"`
			Email           string `json:"email"`
			SmsCode         string `json:"smsCode"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"ConfirmPassword"`
		} `json:"data"`
	}
	var uForm userForm
	err := c.ShouldBind(&uForm)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.Abort()
	}
	unixTimestamp := time.Now().Unix()
	var isExists int
	err = dao.OBCursor.Get(
		&isExists,
		`select count(*) from eln_register_records where email = ? and username = ? and verif_code = ? and verif_code_expire_time > ? and is_del = 0;`,
		uForm.Data.Email,
		uForm.Data.Name,
		uForm.Data.SmsCode,
		unixTimestamp,
	)
	if err != nil {
		utils.BadRequestErr(c, errors.New(err.Error()+"Get user info with error."))
		return
	}
	if isExists == 1 {
		// 如果在注册记录里有，开始写记录到用户表里
		_, err := dao.OBCursor.Exec(
			`insert into eln_users(userid,email,username,password_hash,status,permissions) values (?, ?, ?, ?, 1, ?)`,
			uForm.Data.UserId,
			uForm.Data.Email,
			uForm.Data.Name,
			uForm.Data.Password,
			"user",
		)
		if err != nil {
			utils.Logger.Error(err.Error())
			response := utils.BaseResponse{
				StatusCode: 500, Msg: err.Error(), Data: gin.H{"error": err.Error()},
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := utils.BaseResponse{
			StatusCode: 200, Msg: "User registration successful.", Data: gin.H{},
		}
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, utils.BaseResponse{
			StatusCode: 400, Msg: "Get user info with error.Please fill in the registration information again.", Data: gin.H{"error": "Get user info with error.Please fill in the registration information again."},
		})
		return
	}
}

type ElnUsers struct {
	ID        int64     `json:"id,omitempty" db:"id" gorm:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at" gorm:"created_at"`
	//UpdatedAt    time.Time `json:"updatedAt,omitempty" db:"updated_at" gorm:"updated_at"`
	Userid       int64  `json:"userid,omitempty" db:"userid" gorm:"userid"`
	Username     string `json:"username,omitempty" db:"username" gorm:"username"`
	Email        string `json:"email,omitempty" db:"email" gorm:"email"`
	Phone        string `json:"phone,omitempty" db:"phone" gorm:"phone"`
	PasswordHash string `json:"passwordHash,omitempty" db:"password_hash" gorm:"password_hash"`
	Status       int64  `json:"status,omitempty" db:"status" gorm:"status"`
	GroupId      int64  `json:"groupId,omitempty" db:"group_id" gorm:"group_id"`
	//Avatar       string   `json:"avatar,omitempty" db:"avatar" gorm:"avatar"`
	Permissions  string   `json:"permissions,omitempty" db:"permissions" gorm:"permissions"`
	Authn        string   `json:"authn,omitempty" db:"authn" gorm:"authn"`
	AuthorityIds []string `json:"authorityIds,omitempty"`
}

func GetUserList(c *gin.Context) {
	var users []ElnUsers
	err := dao.OBCursor.Select(&users, `select userid, username, phone, email, permissions, status, created_at from eln_users order by userid`)
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

func SetUserAuthorities(c *gin.Context) {
	var roles struct {
		Userid       int    `json:"userid,omitempty" db:"userid"`
		AuthorityIds string `json:"authorityIds,omitempty" db:"permissions"`
	}
	err := c.ShouldBind(&roles)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	_, err = dao.OBCursor.Exec(`update eln_users set permissions = ? where userid = ?`, roles.AuthorityIds, roles.Userid)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "")
	return
}

type data struct {
	Data struct {
		UserId string `json:"userid"`
		Name   string `json:"username"`
		Email  string `json:"email"`
	} `json:"data"`
}

func SentCodeToWechat(userId string) error {
	code := utils.GenVerifCode()
	unixTimestamp := time.Now().Unix()
	_, err := dao.OBCursor.Exec(
		`update eln_register_records set verif_code = ?, verif_code_expire_time = ? where userid = ?`,
		code,
		unixTimestamp+125,
		utils.GetInterfaceToInt(userId),
	)
	if err != nil {
		return err
	}
	msg, err := utils.SendQwxMsg(userId, code)
	if err != nil {
		return err
	}
	type result struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		ResultType    string `json:"resultType"`
	}
	var r result
	err = json.Unmarshal(msg, &r)
	if err != nil {
		return err
	}
	if r.ResultType == "0" {
		return nil
	} else {
		return errors.New("status BadRequest")
	}
}

func SentCodeOnly(c *gin.Context) {
	var u struct {
		UserId string `json:"userId"`
	}
	err := c.ShouldBind(&u)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	err = SentCodeToWechat(u.UserId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.Success(c, "Please check the verification code through Enterprise WeChat")
}

func SendVerifCode(c *gin.Context) {
	var d data
	err := c.ShouldBind(&d)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.BaseResponse{
			StatusCode: 400, Msg: err.Error(), Data: gin.H{"username": d.Data.Name},
		})
	}
	var isExists int
	err = dao.OBCursor.Get(
		&isExists,
		`select count(*) from eln_users where userid = ? and status = 1`,
		d.Data.UserId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BaseResponse{
			StatusCode: 400, Msg: err.Error(), Data: gin.H{},
		})
		return
	}

	if isExists == 1 {
		c.JSON(http.StatusBadRequest, utils.BaseResponse{
			StatusCode: 400, Msg: "User already exists, please do not register again.", Data: gin.H{},
		})
		return
	}

	_, err = dao.OBCursor.Exec(
		`replace INTO eln_register_records(userid,username,email) values (?, ?, ?)`,
		d.Data.UserId,
		d.Data.Name,
		d.Data.Email,
	)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	err = SentCodeToWechat(d.Data.UserId)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	utils.Success(c, "please check the verification code through Enterprise WeChat")
}

type ElnProject struct {
	ProjectId   int64  `json:"projectId" db:"project_id"`
	ProjectName string `json:"projectName" db:"project_name"`
	UserId      int64  `json:"userId" db:"user_id"`
	RankId      int64  `json:"rankId,omitempty" db:"rank_id"`
	ReactionId  int64  `json:"reactionId" db:"reaction_id"`
	PageName    string `json:"pageName" db:"page_name"`
}

type ElnProjectPage struct {
	ReactionId   int64  `json:"reactionId" db:"reaction_id"`
	UserId       int64  `json:"userId" db:"user_id"`
	ProjectId    int64  `json:"projectId" db:"project_id"`
	PageName     string `json:"pageName" db:"page_name"`
	TimeCreation string `json:"timeCreation" db:"time_creation"`
}

type UnionElnProjectPage struct {
	//ReactionId   int64  `json:"reactionId" db:"reaction_id"`
	//UserId       int64  `json:"userId" db:"user_id"`
	ProjectId   int64  `json:"projectId" db:"project_id"`
	ProjectName string `json:"projectName" db:"project_name"`
	//PageName     string `json:"pageName" db:"page_name"`
	//TimeCreation string `json:"timeCreation" db:"time_creation"`
}

type route struct {
	BookName   string  `json:"bookName"`
	Icon       string  `json:"icon"`
	Path       string  `json:"path"`
	Layout     string  `json:"layout,omitempty"`
	PageName   []route `json:"pageName,omitempty"`
	ReactionId int64   `json:"reactionId,omitempty"`
}

func userRouter(userId int) (error, []route, []route) {
	var subRoutes []route
	var authorityRoutes []route
	var projects []UnionElnProjectPage
	err := dao.OBCursor.Select(&projects, `select b.project_name, a.project_id 
       from eln_project_page a join eln_project b on a.project_id = b.project_id where a.user_id= ? group by project_id`, userId)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err, nil, nil
	}
	for _, project := range projects {
		r := route{
			project.ProjectName,
			"Notebook",
			fmt.Sprintf("/eln/%s", project.ProjectName),
			"",
			nil,
			0,
		}
		var docs []ElnProjectPage
		err := dao.OBCursor.Select(&docs, `select reaction_id, user_id, project_id, page_name, time_creation from eln_project_page where project_id = ? and user_id = ?`, project.ProjectId, userId)
		if err != nil {
			utils.Logger.Error(err.Error())
			return err, nil, nil
		}
		for _, doc := range docs {
			sub := route{
				doc.PageName, "Document", fmt.Sprintf("/eln/%s/%s", project.ProjectName, doc.PageName), "Workbook", nil, doc.ReactionId,
			}
			r.PageName = append(r.PageName, sub)
		}
		subRoutes = append(subRoutes, r)
	}
	//var routes = route{
	//	"route.elnModules", "Collection", "/eln", "Layout", subRoutes, 0,
	//}

	var authProject []struct {
		ProjectName string `json:"projectName" db:"project_name"`
		ReactionId  int64  `json:"reactionId,omitempty" db:"reaction_id"`
	}
	err = dao.OBCursor.Select(&authProject, `select reaction_id, project_name from eln_rxn_basicinfo where witness_id = ?`, userId)
	if err != nil {
		return err, nil, nil
	}
	for _, project := range authProject {
		r := route{
			project.ProjectName,
			"Notebook",
			fmt.Sprintf("/preview/%s", project.ProjectName),
			"",
			nil,
			0,
		}
		var docs []ElnProjectPage
		err := dao.OBCursor.Select(&docs, `select reaction_id, user_id, project_id, page_name, time_creation from eln_project_page where reaction_id = ?`, project.ReactionId)
		if err != nil {
			utils.Logger.Error(err.Error())
			return err, nil, nil
		}
		for _, doc := range docs {
			sub := route{
				doc.PageName, "Document", fmt.Sprintf("/preview/%s/%s", project.ProjectName, doc.PageName), "Preview", nil, doc.ReactionId,
			}
			r.PageName = append(r.PageName, sub)
		}
		authorityRoutes = append(authorityRoutes, r)
	}
	return nil, subRoutes, authorityRoutes
}

func GetAuthorities(c *gin.Context) {
	var x []struct {
		AuthorityId   int64  `json:"authorityId" db:"authority_id"`
		AuthorityName string `json:"authorityName" db:"authority_name"`
	}
	err := dao.OBCursor.Select(&x, `select authority_id,authority_name from eln_user_authority where status = 1`)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	utils.SuccessWithData(c, "", gin.H{"authorities": x})
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
	err = dao.OBCursor.Get(&u, `select * from eln_company_user_list where
                                              chemist_id = ?`, uid.WitnessId)
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
	err = dao.OBCursor.Get(&oldPwd, `select password_hash from eln_users where userid = ?`, p.UserID)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}
	if oldPwd != p.OldPassword {
		utils.BadRequestErr(c, errors.New("the old password was incorrectly entered"))
		return
	}
	_, err = dao.OBCursor.Exec(`update eln_users set password_hash = ? where userid = ?`, p.Password, p.UserID)
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
	err = dao.OBCursor.Get(&s, `select verif_code, verif_code_expire_time from eln_register_records where
                                       userid = ?`, p.UserID)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	duration := s.VerifCodeExpireTime - now
	if duration < 125 && duration > 0 {
		_, err = dao.OBCursor.Exec(`update eln_users set password_hash = ? where userid = ?`, p.Password, p.UserID)
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

func ChangeUserName(c *gin.Context) {
	var u struct {
		UserId   int    `json:"userId"`
		Username string `json:"username"`
	}
	err := c.ShouldBind(&u)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	_, err = dao.OBCursor.Exec(`update eln_users set username = ? where userid = ?`, u.Username, u.UserId)
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
	err = dao.OBCursor.Get(&permission, `select permissions from eln_users where userid = ?`, s.RequestId)
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	if strings.Contains(permission, "admin") {
		_, err = dao.OBCursor.Exec(`update eln_users set password_hash = ? where userid = ?`, "e10adc3949ba59abbe56e057f20f883e", s.UserID)
		if err != nil {
			utils.InternalRequestErr(c, err)
			return
		}

		utils.Success(c, "change passcode successfully")
		return
	}

	utils.InternalRequestErr(c, errors.New("no permission"))
}
