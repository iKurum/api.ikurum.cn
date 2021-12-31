package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
)

type bdocr struct {
	pid      []uint8
	pidName  []uint8
	ocrid    uint8
	title    string
	quantity int64
}

func init() {
	route.Mux.GET("/v1/ocr/list", func(rw http.ResponseWriter, r *http.Request) {
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

		var count int64 = 0
		DB.QueryRow("select count(*) from bdocr").Scan(&count)

		result, err := DB.Query("select pid+0,pid,ocrid,title,quantity from bdocr order by pid")
		global.CheckErr(err, "")

		data := make(map[string][]interface{})
		for result.Next() {
			var d bdocr
			err = result.Scan(
				&d.pid,
				&d.pidName,
				&d.ocrid,
				&d.title,
				&d.quantity,
			)
			if err != nil {
				log.Fatal(err)
				break
			}
			data[B2S(d.pidName).(string)] = append(data[B2S(d.pidName).(string)], map[string]interface{}{
				"pid":      B2S(d.pid),
				"pName":    B2S(d.pidName),
				"id":       d.ocrid,
				"title":    d.title,
				"quantity": d.quantity,
			})
		}

		if err != nil {
			msg, _ := json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  fmt.Sprint(err),
			}))
			rw.Write(msg)
			return
		}

		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code:  200,
			Count: count,
			Data:  data,
		}))
		rw.Write(msg)
	})
}

func B2S(bs []uint8) interface{} {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}

	b, err := strconv.Atoi(string(ba))
	if err != nil {
		return string(ba)
	}
	return b
}
