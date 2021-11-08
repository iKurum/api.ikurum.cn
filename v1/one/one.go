package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/one", func(rw http.ResponseWriter, r *http.Request) {
		global.SetHeader(rw)

		var data string
		var code int
		// m := global.GetByEssay("oneBucket")

		// if m["data"] != "" {
		// 	code = 200
		// 	a := strings.Split(m["data"], "*_*")

		// 	rand.Seed(time.Now().Unix())
		// 	data = a[rand.Intn(len(a))]
		// } else {
		// 	code = 0
		// 	data = "get something error"
		// }

		fmt.Println("一言:", data)
		msg, _ := json.Marshal(global.NewResult(&global.Result{Code: code, Data: data}))
		rw.Write(msg)
	})
}
