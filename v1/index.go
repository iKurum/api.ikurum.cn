package v1

import (
	"encoding/json"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"

	_ "api.ikurum.cn/v1/article"
	_ "api.ikurum.cn/v1/user"
)

func init() {
	route.Mux.GET("/v1/index", func(rw http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 200, Data: "This is api.ikurum.cn"}))
		rw.Write(msg)
	})
}
