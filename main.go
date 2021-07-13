package main

import (
	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util"

	_ "api.ikurum.cn/v1"
)

func main() {
	err := global.BoltInit()
	if err == nil {
		global.SetOne()
	}

	go util.StartToken()

	r := &route.Router{}
	r.Listen("9091")
}
