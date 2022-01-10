package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util/ocr"
)

func init() {
	route.POST("/v1/ocr/text", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg, _ := json.Marshal(global.NewResult(&global.Result{
					Code: 0,
					Msg:  "recover:" + fmt.Sprint(err),
				}))
				rw.Write(msg)
			}
		}()

		r.ParseMultipartForm(32 << 20)

		file, _, err := r.FormFile("image")
		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "image:" + fmt.Sprint(err),
			}))
			rw.Write(msg)
			return
		}
		first, err := strconv.Atoi(r.FormValue("first"))
		if err != nil || first == 0 {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "参数first错误",
			}))
			rw.Write(msg)
			return
		}

		second, err := strconv.Atoi(r.FormValue("second"))
		if err != nil || second == 0 {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  "参数second错误",
			}))
			rw.Write(msg)
			return
		}

		b, _ := ioutil.ReadAll(file)

		ocr, err := ocr.Read_file(b, first, second)
		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: -1,
				Msg:  ocr,
			}))
			rw.Write(msg)
			return
		}

		var (
			count  int64
			ocr_le []interface{}
		)
		if o, ok := ocr.(map[string]interface{}); ok {
			if ocr_le, ok = o["result"].([]interface{}); ok {
				count = int64(len(ocr_le))
			}
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code:  200,
			Data:  ocr,
			Count: count,
		}))
		rw.Write(msg)
	})
}
