package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

func init() {
	route.GET("/getshort/", func(rw http.ResponseWriter, r *http.Request) {
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
		var (
			u   = query.Get("short")
			msg []byte
			sr  string
		)
		fmt.Println("short url:", u)

		if u == "" {
			msg, _ = json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "short 参数错误",
			}))
		} else {
			DB := global.OpenDB()
			err := DB.QueryRow("select url from short where surl=?", u).Scan(&sr)
			if err = global.CheckErr(err); err == nil {
				msg, _ = json.Marshal(global.NewResult(&global.Result{
					Code:  200,
					Data:  sr,
					Count: 1,
				}))
			} else {
				msg, _ = json.Marshal(global.NewResult(&global.Result{
					Code: 0,
					Msg:  err,
				}))
			}
		}

		if sr == "" {
			rw.Write(msg)
		} else {
			if strings.HasPrefix(sr, "https://") || strings.HasPrefix(sr, "http://") {
				http.Redirect(rw, r, sr, http.StatusFound)
			} else {
				sDec, err := base64.StdEncoding.DecodeString(sr)
				if err != nil {
					msg, _ = json.Marshal(global.NewResult(&global.Result{
						Code: 0,
						Msg:  err,
					}))
					rw.Write(msg)
				} else {
					http.Redirect(rw, r, string(sDec), http.StatusFound)
				}
			}
		}
	})
}
