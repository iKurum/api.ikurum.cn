package ocr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func getImg(f string, name string, t string, image_url string) []byte {
	var params = make(url.Values)
	params.Add("image", f)
	post_data := params.Encode()

	// 调用文字识别服务
	fmt.Println("连接 IMG 服务 ...", name)

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
		if k == "result" {
			var a = v.([]interface{})
			for i := 0; i < len(a); i++ {
				switch t {
				case "8":
					var b = a[i].(map[string]interface{})
					if b["score"].(float64) > 0.5 {
						txt += "brief: " + b["brief"].(string) + "\n"
						txt += "score: " + fmt.Sprintf("%v", b["score"].(float64)) + "\n"
						txt += "\n"
					}
				default:
					var b = a[i].(map[string]interface{})
					txt += "keyword: " + b["keyword"].(string) + "\n"
					txt += "root: " + b["root"].(string) + "\n"
					txt += "score: " + fmt.Sprintf("%v", b["score"].(float64)) + "\n"
					txt += "\n"
				}
			}
		}
	}

	if txt != "" {
		fmt.Println("图片 IMG 内容:\n", txt)
	} else {
		fmt.Printf("!!! %v IMG 内容无\n", name)
	}

	return []byte("")
}
