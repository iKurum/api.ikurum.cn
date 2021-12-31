package main

import (
	"log"
	"os"

	_ "api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util"
	_ "api.ikurum.cn/v1"
)

func main() {
	log.SetPrefix("[IKURUM]~")
	log.SetFlags(2)
	log.SetOutput(os.Stdout)

	go util.StartToken()
	r := &route.Router{}
	r.Listen("9091")
}
