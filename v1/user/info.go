package v1

import (
	"encoding/json"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/user/info", func(rw http.ResponseWriter, r *http.Request) {
		global.SetHeader(rw)

		info := global.GetByBucket("userInfo")
		d := map[string]interface{}{
			"surname": info["surname"],
			"mail":    info["mail"],
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 200, Data: d}))
		rw.Write(msg)
	})
}
