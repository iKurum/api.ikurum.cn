package ocr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func getTxt(f string, name string, t string, image_url string) []byte {
	var params = make(url.Values)
	params.Add("image", f)
	if t == "2" {
		params.Add("id_card_side", "back")
	}
	post_data := params.Encode()

	// 调用文字识别服务
	fmt.Println("连接 OCR 服务 ...", name)

	req, _ := http.NewRequest("POST", image_url, strings.NewReader(post_data))
	req.Header.Set("Content-Type", "application/application/x-www-form-urlencoded")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	jsonTxt, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(jsonTxt, &j)

	var txt string = ""
	for k, v := range j {
		if k == "words_result" {
			switch t {
			case "4":
			case "2":
				if j["image_status"] == "normal" {
					var m = v.(map[string]interface{})
					for key, value := range m {
						txt += key + ":" + value.(map[string]interface{})["words"].(string) + "\n"
					}
				} else {
					fmt.Printf("%v\n", j["image_status"])
				}
			case "8":
				var m = v.(map[string]interface{})
				for key, value := range m {
					txt += key + ":" + value.(string) + "\n"
				}
			default:
				var m = v.([]interface{})
				for i := 0; i < len(m); i++ {
					txt += m[i].(map[string]interface{})["words"].(string)
				}
			}
		}
	}

	return []byte(txt)
}
