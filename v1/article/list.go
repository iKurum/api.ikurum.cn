package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/article/list", func(rw http.ResponseWriter, r *http.Request) {
		global.SetHeader(rw)
		DB := global.OpenDB()

		var count int64
		DB.QueryRow("select count(*) from essay").Scan(&count)
		fmt.Println("count", count)

		query := r.URL.Query()
		page := 1
		size := 10

		var err error
		if query.Get("page") != "" {
			page, err = strconv.Atoi(query.Get("page"))
			if err != nil {
				msg, _ := json.Marshal(global.NewResult(&global.Result{
					Code: 0,
					Msg:  "参数page错误",
				}))
				rw.Write(msg)
				return
			}
			fmt.Println("query page:", page)
		}
		if query.Get("size") != "" {
			size, err = strconv.Atoi(query.Get("size"))
			if err != nil {
				msg, _ := json.Marshal(global.NewResult(&global.Result{
					Code: 0,
					Msg:  "参数size错误",
				}))
				rw.Write(msg)
				return
			}
			fmt.Println("query size:", size)
		}

		// 超出数据
		if int64((page-1)*size) >= count {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "没有更多数据",
				Page: page,
				Size: size,
			}))
			rw.Write(msg)
			return
		}

		d := make([]interface{}, size)

		result, err := DB.Query("select aid,size,title,uptime,addtime,note from essay where aid > ? and aid < ?", (page-1)*size, page*size+1)
		global.CheckErr(err, "")

		index := 0
		for result.Next() {
			var data global.Essay_list
			err = result.Scan(&data.Id, &data.Size, &data.Title, &data.Uptime, &data.Addtime, &data.Note)
			if err != nil {
				break
			}
			d[index] = map[string]interface{}{
				"id":      data.Id,
				"size":    data.Size,
				"title":   data.Title,
				"addTime": data.Addtime,
				"note":    data.Note,
			}
			index++
		}
		result.Close()

		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  fmt.Sprint(err),
				Page: page,
				Size: size,
			}))
			rw.Write(msg)
			return
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: d,
			Page: page,
			Size: size,
			More: int64(page*size) < count,
		}))
		rw.Write(msg)
	})
}
