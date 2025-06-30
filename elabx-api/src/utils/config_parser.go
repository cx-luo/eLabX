// Package utils coding=utf-8
// @Project : eLabX
// @Time    : 2024/1/16 9:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : config_parser.go
// @Software: GoLand
package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
	Service Service `json:"service" yaml:"service"`
	Redis   Redis   `json:"redis" yaml:"redis"`
	Output  Output  `json:"output" yaml:"output"`
}
type Mysql struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Passwd   string `json:"passwd" yaml:"passwd"`
	Database string `json:"database" yaml:"database"`
}
type Service struct {
	Port          int    `json:"port" yaml:"port"`
	Backup        bool   `json:"backup" yaml:"backup"`
	UploadFileDir string `json:"uploadFileDir" yaml:"uploadFileDir"`
}
type Redis struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type Output struct {
	Loglevel string `json:"loglevel" yaml:"loglevel"`
	Logfile  string `json:"logfile" yaml:"logfile"`
	//ErrFile    string `json:"errFile" yaml:"errFile"`
	ServiceLog string `json:"serviceLog" yaml:"serviceLog"`
	MaxSize    int    `json:"maxSize" yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
}

type Indigo struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

func GetConfigData(configPath string) (*Config, error) {
	//导入配置文件
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config *Config
	//将配置文件读到结构体中
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err = viper.Unmarshal(&config)
		if err != nil {
			panic(err)
		}
	})

	return config, nil
}

var GlobalConfig *Config
