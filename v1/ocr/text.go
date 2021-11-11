package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util/ocr"
)

func init() {
	route.Mux.POST("/v1/ocr/text", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		file, fileHandle, err := r.FormFile("image")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(fileHandle.Header, fileHandle.Filename, fileHandle.Size)

		b, _ := ioutil.ReadAll(file)

		msg, _ := json.Marshal(global.NewResult(&global.Result{
			Code: 200,
			Data: string(ocr.Read_file(b, 1.1)),
		}))
		rw.Write(msg)
	})
}
