package config

import (
	"time"
)

// true		打包后，连接数据库
// false	本地启项目，连接远端数据库
var Online = true

// token刷新时间
var SetTokenTime time.Duration = time.Duration(1) * time.Hour

// 数据库连接信息
var DB = map[string]string{
	"title":    "mysql",
	"user":     "",
	"pw":       "",
	"port":     "",
	"database": "",
	"ip":       "", // 在init中配置
}

func init() {
	if Online {
		DB["ip"] = "127.0.0.1"
	} else {
		DB["ip"] = "/* 远端ip */"
	}
}
