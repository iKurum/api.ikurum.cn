package main

import (
	_ "api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util"
	"api.ikurum.cn/util/logs"
	_ "api.ikurum.cn/v1"
)

func main() {
	logs.Init()

	go util.StartToken()
	r := &route.Router{}
	r.Listen("9091")
}
