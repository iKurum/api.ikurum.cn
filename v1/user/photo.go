package v1

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/user/photo", func(rw http.ResponseWriter, r *http.Request) {
		global.SetHeader(rw)

		var msg []byte
		data, err := os.Open("photo.jpg")
		println("获取头像：", err == nil)
		d, _ := ioutil.ReadAll(data)
		if err != nil {
			msg, _ = json.Marshal(global.NewResult(&global.Result{Code: 0, Data: err}))
			rw.Write(msg)
		} else {
			msg, _ = json.Marshal(global.NewResult(&global.Result{Code: 200, Data: base64.StdEncoding.EncodeToString(d)}))
			rw.Write(msg)
		}
	})
}
