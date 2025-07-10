// Package router coding=utf-8
// @Project : eLabX
// @Time    : 2023/12/12 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : main.go
// @Software: GoLand

package main

import (
	"eLabX/middleware"
	"eLabX/router"
	"eLabX/src/dao"
	"eLabX/src/utils"
	"flag"
	"fmt"
)

var (
	BuildVersion   = ""
	BuildGoVersion = ""
	BuildCommit    = ""
	BuildTime      = ""
)

func main() {
	// 构建信息，golang版本 commit id 时间
	var version bool
	var configPath string
	flag.BoolVar(&version, "V", false, "version")
	flag.StringVar(&configPath, "c", "conf/.env.yaml", "config")
	flag.Parse()

	if version {
		fmt.Printf("go version: %s\nBuild version: %s\n Build commit: %s\nBuild time: %s\n",
			BuildGoVersion, BuildVersion, BuildCommit, BuildTime)
		return
	}

	config, err := utils.GetConfigData(configPath)
	if err != nil {
		panic(err)
	}
	utils.GlobalConfig = config
	err = middleware.InitLogger()
	if err != nil {
		panic(err)
	}
	dao.OBCursor = dao.GetMysqlCursor(utils.GlobalConfig.Mysql.Host, utils.GlobalConfig.Mysql.Port, utils.GlobalConfig.Mysql.User, utils.GlobalConfig.Mysql.Passwd, utils.GlobalConfig.Mysql.Database)
	r := router.NewRouter(config.Output.Logfile, config.Output.Loglevel)
	utils.GetAllRoutes(r)
	err = r.Run(fmt.Sprintf(":%d", utils.GlobalConfig.Service.Port))
	if err != nil {
		panic(err)
	}
	//utils.Logger.Info(r.Run(fmt.Sprintf(":%d", utils.GlobalConfig.Service.Port)))
}
