package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.Mux.GET("/v1/article/list", func(rw http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		page := "1"
		size := "10"

		if query.Get("page") != "" {
			page = query.Get("page")
			fmt.Println("query page:", page)
		}
		if query.Get("size") != "" {
			size = query.Get("size")
			fmt.Println("query size:", size)
		}

		d := make(map[int]interface{})

		str := global.GetByDB("detailList", "id")
		arr := strings.Split(str, ",")

		if len(arr) != 0 && arr[len(arr)-1] == "" {
			arr = arr[0 : len(arr)-1]
		}
		fmt.Println("detail arr:", arr)

		p, err := strconv.Atoi(page)
		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "page错误",
			}))
			rw.Write(msg)
		}
		s, err := strconv.Atoi(size)
		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "size错误",
			}))
			rw.Write(msg)
		}

		fmt.Println("detail:", arr)
		if len(arr) == 0 {
			d[0] = map[string]string{
				"error": "没有数据",
			}
		} else {
			fmt.Println("开始获取详情 ...", p, s)
			tep := 0
			for i := (p - 1) * s; i < p*s; i++ {
				fmt.Printf("detail len: %d\t%t\n", len(arr), len(arr) > i)

				if len(arr) > i {
					m := global.GetByBucket(arr[i])
					if m["id"] != "" {
						fmt.Println("get detail data:", m["name"])

						d[tep] = map[string]string{
							"id":    m["id"],
							"cTime": m["createdDateTime"],
							"mTime": m["lastModifiedDateTime"],
							"name":  m["name"],
							"note":  m["note"],
							"size":  m["size"],
						}
						tep++
					}
				}
			}
		}

		a := make([]interface{}, len(d))
		for k, v := range d {
			a[k] = v
		}

		// fmt.Printf("more: %v;len: %d, %d * %d = %d\n", len(arr) > p*s, len(arr), p, s, p*s)
		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: a,
			Page: p,
			Size: s,
			More: len(arr) > p*s,
		}))
		rw.Write(msg)
	})
}
