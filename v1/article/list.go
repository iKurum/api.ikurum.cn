package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/article/list", func(rw http.ResponseWriter, r *http.Request) {
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

		query := r.URL.Query()
		var page int64 = 1
		var size int64 = 10
		var archive string = query.Get("archive")

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

		var count int64
		if archive != "" {
			DB.QueryRow("select count(*) from essay where archive like ?", "%"+archive+"%").Scan(&count)
			fmt.Printf("select essay like %s: count %d\n", archive, count)
		} else {
			DB.QueryRow("select count(*) from essay").Scan(&count)
			fmt.Println("select essay count", count)
		}

		var csize int64 = size
		if (page-1)*size >= count {
			// 超出数据
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
		d := make([]interface{}, csize)

		var (
			result *sql.Rows
			stmt   *sql.Stmt
		)
		if page > 1 {
			stmt, err = DB.Prepare("select aid,size,title,addtime,note,archive from essay  where archive like ? order by addtime desc limit ?,?")
		} else {
			stmt, err = DB.Prepare("select aid,size,title,addtime,note,archive from essay where archive like ? order by addtime desc limit ?")
		}
		global.CheckErr(err, "")

		if page > 1 {
			result, err = stmt.Query("%"+archive+"%", csize, (page-1)*size)
		} else {
			result, err = stmt.Query("%"+archive+"%", csize)
		}
		global.CheckErr(err, "")

		index := 0
		for result.Next() {
			var data global.Essay_list
			err = result.Scan(
				&data.Id,
				&data.Size,
				&data.Title,
				&data.Addtime,
				&data.Note,
				&data.Archive,
			)
			if err != nil {
				break
			}
			d[index] = map[string]interface{}{
				"id":      data.Id,
				"size":    data.Size,
				"title":   data.Title,
				"addTime": data.Addtime,
				"note":    data.Note,
				"archive": data.Archive,
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
			Code:  200,
			Data:  d,
			Page:  page,
			Size:  size,
			More:  int64(page*size) < count,
			Count: csize,
		}))
		rw.Write(msg)
	})
}
