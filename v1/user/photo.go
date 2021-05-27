package v1

import (
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/user/photo", func(rw http.ResponseWriter, r *http.Request) {
		str := global.GetByDB("photo", "str")

		rw.Write([]byte(str))
	})
}
