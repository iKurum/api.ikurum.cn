package v1

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.GET("/v1/one", func(rw http.ResponseWriter, r *http.Request) {
		DB := global.OpenDB()

		var count int64
		DB.QueryRow("select count(*) from one").Scan(&count)

		//将时间戳设置成种子数
		rand.Seed(time.Now().UnixNano())

		var data string
		err := DB.QueryRow("select md from one where oid=?", rand.Int63n(count-1)+1).Scan(&data)
		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  fmt.Sprint("查询错误", err),
			}))
			rw.Write(msg)
			return
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: data,
		}))
		rw.Write(msg)
	})
}
