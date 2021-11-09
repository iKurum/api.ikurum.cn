package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/article/list", func(rw http.ResponseWriter, r *http.Request) {
		global.SetHeader(rw)
		DB := global.OpenDB()

		defer func() {
			if err := recover(); err != nil {
				msg, _ := json.Marshal(global.NewResult(&global.Result{
					Code: 0,
					Msg:  fmt.Sprint(err),
				}))
				rw.Write(msg)
			}
		}()

		var count int64
		DB.QueryRow("select count(*) from essay").Scan(&count)
		fmt.Println("count", count)

		query := r.URL.Query()
		var page int64 = 1
		var size int64 = 10

		var err error
		if query.Get("page") != "" {
			page, err = strconv.ParseInt(query.Get("page"), 10, 64)
			if err != nil || page < 1 {
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
			size, err = strconv.ParseInt(query.Get("size"), 10, 64)
			if err != nil || size < 1 || size > 20 {
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
		var csize int64 = size
		if (page-1)*size >= count {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "没有更多数据",
				Page: page,
				Size: size,
			}))
			rw.Write(msg)
			return
		} else {
			if page*size >= count {
				csize = count - (page-1)*size
			}
		}

		log.Println("d size:", csize)
		d := make([]interface{}, csize)

		var result *sql.Rows
		if page > 1 {
			stmt, err := DB.Prepare("select aid,size,title,addtime,note from essay order by addtime desc limit ?,?")
			global.CheckErr(err, "")
			result, err = stmt.Query(csize, (page-1)*size)
			global.CheckErr(err, "")
		} else {
			stmt, err := DB.Prepare("select aid,size,title,addtime,note from essay order by addtime desc limit ?")
			global.CheckErr(err, "")
			result, err = stmt.Query(csize)
			global.CheckErr(err, "")
		}

		index := 0
		for result.Next() {
			var data global.Essay_list
			err = result.Scan(&data.Id, &data.Size, &data.Title, &data.Addtime, &data.Note)
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
			if index >= int(csize) {
				break
			}
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
