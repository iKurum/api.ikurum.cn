package v1

import (
	"encoding/json"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util/logs"

	_ "api.ikurum.cn/v1/article"
	_ "api.ikurum.cn/v1/ocr"
	_ "api.ikurum.cn/v1/one"
	_ "api.ikurum.cn/v1/server"
	_ "api.ikurum.cn/v1/user"
	"github.com/thinkeridea/go-extend/exnet"
)

func init() {
	route.GET("/v1/index", func(rw http.ResponseWriter, r *http.Request) {
		global.OpenDB()
		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: "database connected",
			Msg:  "This is api.ikurum.cn",
		}))
		rw.Write(msg)
	})

	route.GET("/v1/getip", func(rw http.ResponseWriter, r *http.Request) {
		ip := exnet.ClientPublicIP(r)
		if ip == "" {
			ip = exnet.ClientIP(r)
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 200, Data: ip}))
		rw.Write(msg)
	})

	logs.Info("初始化路由")
}
