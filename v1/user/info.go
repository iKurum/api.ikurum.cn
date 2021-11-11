package v1

import (
	"encoding/json"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/user/info", func(rw http.ResponseWriter, r *http.Request) {
		DB := global.OpenDB()

		var (
			name  string
			email string
		)
		err := DB.QueryRow("select name, email from user where uid=1").Scan(&name, &email)
		global.CheckErr(err, "")

		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: map[string]string{
				"name":  name,
				"email": email,
			},
		}))
		rw.Write(msg)
	})
}
