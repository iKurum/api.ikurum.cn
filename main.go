package main

import (
	"api.ikurum.cn/config"
	_ "api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util/logs"
	_ "api.ikurum.cn/v1"
)

func init() {
	// 初始化log
	logs.Init()
	// 读取 ApiConfig.yaml 配置
	config.GetApiConfig()
}

func main() {
	r := &route.Router{}
	r.Listen("9091")
}
