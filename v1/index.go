package v1

import (
	"encoding/json"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"github.com/thinkeridea/go-extend/exnet"

	_ "api.ikurum.cn/v1/article"
	_ "api.ikurum.cn/v1/ocr"
	_ "api.ikurum.cn/v1/one"
	_ "api.ikurum.cn/v1/user"
)

func init() {
	route.Mux.GET("/v1/index", func(rw http.ResponseWriter, r *http.Request) {
		global.OpenDB()
		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: "database connect",
			Msg:  "This is api.ikurum.cn",
		}))
		rw.Write(msg)
	})

	route.Mux.GET("/v1/getip", func(rw http.ResponseWriter, r *http.Request) {
		ip := exnet.ClientPublicIP(r)
		if ip == "" {
			ip = exnet.ClientIP(r)
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 200, Data: ip}))
		rw.Write(msg)
	})
}
