package main

import (
	"api.ikurum.cn/route"
	"api.ikurum.cn/util"
	_ "api.ikurum.cn/v1"
)

func main() {
	go util.StartToken()

	r := &route.Router{}
	r.Listen("9091")
}
