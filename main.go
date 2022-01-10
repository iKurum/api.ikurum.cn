package main

import (
	_ "api.ikurum.cn/global"
	"api.ikurum.cn/route"
	_ "api.ikurum.cn/v1"
)

func main() {
	r := &route.Router{}
	r.Listen("9091")
}
