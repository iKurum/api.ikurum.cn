package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.POST("/v1/article/detail", func(rw http.ResponseWriter, r *http.Request) {
		// 根据请求body创建一个json解析器实例
		decoder := json.NewDecoder(r.Body)

		// 用于存放参数key=value数据
		var params map[string]string
		// 解析参数 存入map
		decoder.Decode(&params)
		fmt.Println("params:", params)

		id := params["id"]
		if id != "" {
			m := global.GetByEssay(id)
			fmt.Println(m.Err)
			if m.Err != "" {
				msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 0, Msg: m.Err}))
				rw.Write(msg)
			} else {
				msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 200, Data: m}))
				rw.Write(msg)
			}
		} else {
			msg, _ := json.Marshal(global.NewResult(&global.Result{Code: 0, Msg: "缺少参数id"}))
			rw.Write(msg)
		}
	})
}
