package v1

import (
	"encoding/json"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/user/photo", func(rw http.ResponseWriter, r *http.Request) {
		global.SetHeader(rw)
		DB := global.OpenDB()

		var photo string
		err := DB.QueryRow("select photo from user where uid=1").Scan(&photo)
		global.CheckErr(err, "")

		msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 200, Data: photo}))
		rw.Write(msg)
	})
}
