// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2023/12/12 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@pharmaron.com
// @File    : router.go
// @Software: GoLand
package router

import (
	_ "eLabX/docs"
	"eLabX/src/api"
	middleware2 "eLabX/src/middleware"
	"eLabX/src/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// NewRouter returns a new router.
func NewRouter(outputPath string, loglevel string) *gin.Engine {
	// 设置全局 Logger
	logger := utils.SetupLogger(outputPath, loglevel)

	// 延迟关闭 logger
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}(logger)

	router := gin.New()

	// 为需要中间件的路由组注册中间件

	// 使用 Zap 中间件
	router.Use(utils.GinLogger(logger), utils.GinRecovery(logger, true))

	// 注册其他中间件
	router.Use(middleware2.CORS(), middleware2.JwtAuth(), middleware2.Oplog())

	// 注册用户相关路由
	registerUserRoutes(router)

	// 注册文件相关路由
	registerFileRoutes(router)

	// 注册 PUG 相关路由
	registerPugRoutes(router)

	// 注册权限相关路由
	registerAuthorityRoutes(router)

	// 注册 ELN 相关路由
	registerElnRoutes(router)

	// 注册项目管理相关路由
	registerReportMgeRoutes(router)

	// 注册其他路由
	registerOtherRoutes(router)

	return router
}

// 用户相关路由
func registerUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/login", api.UserLogin)
		userGroup.POST("/logout", api.UserLogout)
		userGroup.GET("/userInfo", api.UserInfo)
		userGroup.POST("/register", api.UserRegister)
		userGroup.POST("/fetchUserName", api.FetchUserName)
		userGroup.GET("/getUserList", api.GetUserList)
		userGroup.POST("/setUserAuthorities", api.SetUserAuthorities)
		userGroup.POST("/genVerifCode", api.SendVerifCode)
		userGroup.POST("/getRecentUsage", api.GetRecentUsage)
		userGroup.GET("/getAuthorities", api.GetAuthorities)
		userGroup.POST("/changePwd", api.ChangePwd)
		userGroup.POST("/changeUsername", api.ChangeUserName)
		userGroup.POST("/forgetPwd", api.ForgetPwd)
		userGroup.POST("/resetPwd", api.ResetPwd)
		userGroup.POST("/sentCodeOnly", api.SentCodeOnly)
	}
}

// 文件相关路由
func registerFileRoutes(r *gin.Engine) {
	fileGroup := r.Group("/api/filesvr")
	{
		fileGroup.POST("/upload", api.UploadFile)
		fileGroup.POST("/saveFilename", api.SaveFilenameToDb)
	}
}

// PUG 相关路由
func registerPugRoutes(r *gin.Engine) {
	pugGroup := r.Group("/api/pug")
	{
		pugGroup.POST("/getCmpdInfoBySmiles", api.GetCmpdInfoBySmiles)
		pugGroup.POST("/getCmpdInfoByName", api.GetCmpdInfoByName)
		pugGroup.POST("/cmpdComboSearch", api.CmpdComboSearch)
	}
}

// 权限相关路由
func registerAuthorityRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/authority")
	{
		authGroup.POST("/getAuthorityList", api.GetAuthorityList)
		authGroup.POST("/deleteAuthority", api.DeleteAuthority)
		authGroup.POST("/copyAuthority", api.CopyAuthority)
		authGroup.POST("/disableAuthority", api.DisableAuthority)
		authGroup.GET("/getApiList", api.GetApiList)
		authGroup.POST("/getRxnForAuthority", api.GetRxnForAuthority)
		authGroup.POST("/authorityStatusChange", api.AuthorityStatusChange)
		authGroup.POST("/getRxnDetail", api.SignatureReport)
	}
}

// ELN 相关路由
func registerElnRoutes(r *gin.Engine) {
	elnGroup := r.Group("/api/eln")
	{
		elnGroup.POST("/saveReactionNote", api.SaveReactionNote)
		elnGroup.POST("/createNewPage", api.CreateNewPage)
		elnGroup.POST("/getReagentInfoByRole", api.GetReagentInfoByRole)
		elnGroup.POST("/loadWorkbook", api.LoadWorkbook)
		elnGroup.POST("/loadReagents", api.LoadReagents)
		elnGroup.POST("/saveAdditionalInfo", api.SaveAdditionalInfo)
		elnGroup.POST("/saveNewReagent", api.SaveNewReagent)
		elnGroup.POST("/deleteReagent", api.DeleteReagent)
		elnGroup.POST("/setReagentCid", api.SetReagentCid)
		elnGroup.POST("/saveRxnCondition", api.SaveReactionConditions)
		elnGroup.POST("/getSamples", api.GetSamples)
		elnGroup.POST("/commitRxn", api.CommitRxn)
		elnGroup.POST("/saveSamples", api.SaveSamples)
		elnGroup.POST("/saveLcms", api.SaveLcms)
		elnGroup.POST("/saveNmr", api.SaveNmr)
		elnGroup.POST("/saveOtherReport", api.SaveOtherReport)
		elnGroup.POST("/deleteLcmsReport", api.DeleteLcmsReport)
		elnGroup.POST("/deleteNmrReport", api.DeleteNmrReport)
		elnGroup.POST("/deleteOtherReport", api.DeleteOtherReport)
		elnGroup.POST("/saveNewStdReagent", api.SaveNewStdReagent)
		elnGroup.POST("/getProcedureTemplate", api.GetProcedureTemplate)
		elnGroup.POST("/getProcedure", api.GetProcedure)
		elnGroup.POST("/saveProcedure", api.SaveProcedure)
		elnGroup.POST("/saveProcedureComments", api.SaveProcedureComments)
		elnGroup.POST("/getBasicInfo", api.GetBasicInfo)
		elnGroup.POST("/reopenRxn", api.ReopenRxn)
		elnGroup.POST("/saveBasicInfo", api.SaveBasicInfo)
		elnGroup.POST("/copyPage", api.CopyPage)
	}
}

// 项目管理相关路由
func registerReportMgeRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api/report")
	{
		otherGroup.POST("/getProjectPageInfo", api.GetProjectPageInfo)
		otherGroup.POST("/getProject", api.GetProject)
		otherGroup.POST("/deleteProject", api.DeleteProject)
		otherGroup.POST("/deletePage", api.DeletePage)
		otherGroup.POST("/searchProject", api.SearchProject)
	}
}

// 其他路由
func registerOtherRoutes(r *gin.Engine) {
	otherGroup := r.Group("/api")
	{
		otherGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		otherGroup.POST("/render", api.GenImg)
	}
}
