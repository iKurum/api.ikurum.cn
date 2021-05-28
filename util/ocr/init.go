package ocr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"api.ikurum.cn/global"
)

// 获取 token
func fetch_token() error {
	fmt.Println("获取更新 token ...")

	var params = make(url.Values)
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", global.API_KEY)
	params.Add("client_secret", global.SECRET_KEY)

	post_data := params.Encode()

	// 获取token
	req, _ := http.NewRequest("POST", global.TOKEN_URL, strings.NewReader(post_data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	jsonTxt, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(jsonTxt, &j)

	global.Access_token = j["access_token"].(string)

	if global.Access_token != "" {
		return nil
	} else {
		log.Fatal("获取 token 错误:", j)
		return fmt.Errorf("获取 token 错误")
	}
}

// 读取图片
func Read_file(file []byte, t float64) []byte {
	var rt []byte
	fmt.Println("正在读取图片 ...")

	err := fetch_token()
	if err == nil {
		var s = make([]string, 2)
		s = strings.Split(fmt.Sprintf("%v", t), ".")
		var url []map[string]string

		fmt.Println("type:", s[0], s[1])
		switch s[0] {
		case "1":
			url = global.OCR_URL
		case "2":
			url = global.FACE_URL
		case "3":
			url = global.IMAGE_URL
		}

		for i := 0; i < len(url); i++ {
			if url[i]["type"] == s[1] {
				switch s[0] {
				case "1":
					rt = getTxt(base64.StdEncoding.EncodeToString(file), url[i]["name"], url[i]["type"], global.API_URL+url[i]["url"]+"?access_token="+global.Access_token)
				case "2":
					rt = getFace(base64.StdEncoding.EncodeToString(file), url[i]["name"], url[i]["type"], global.API_URL+url[i]["url"]+"?access_token="+global.Access_token)
				case "3":
					rt = getImg(base64.StdEncoding.EncodeToString(file), url[i]["name"], url[i]["type"], global.API_URL+url[i]["url"]+"?access_token="+global.Access_token)
				}
			}
		}
	}
	return rt
}
