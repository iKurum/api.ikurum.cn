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

	"api.ikurum.cn/config"
	"api.ikurum.cn/global"
)

type baidu struct {
	API_KEY    string
	SECRET_KEY string
	TOKEN_URL  string
	API_URL    string
}

var bai baidu

// 获取 token
func fetch_token() error {
	fmt.Println("获取更新 baidu token ...")

	DB := global.OpenDB()
	err := DB.QueryRow("select API_KEY,SECRET_KEY,TOKEN_URL,API_URL from global where gid=1").Scan(
		&bai.API_KEY,
		&bai.SECRET_KEY,
		&bai.TOKEN_URL,
		&bai.API_URL,
	)
	global.CheckErr(err, "")

	var params = make(url.Values)
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", bai.API_KEY)
	params.Add("client_secret", bai.SECRET_KEY)

	post_data := params.Encode()

	// 获取token
	req, _ := http.NewRequest("POST", bai.TOKEN_URL, strings.NewReader(post_data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	jsonTxt, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(jsonTxt, &j)

	config.Baidu_Access_token = j["access_token"].(string)

	if config.Baidu_Access_token != "" {
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
			url = config.OCR_URL
		case "2":
			url = config.FACE_URL
		case "3":
			url = config.IMAGE_URL
		}

		for i := 0; i < len(url); i++ {
			if url[i]["type"] == s[1] {
				switch s[0] {
				case "1":
					rt = getTxt(base64.StdEncoding.EncodeToString(file), url[i]["name"], url[i]["type"], bai.API_URL+url[i]["url"]+"?access_token="+config.Baidu_Access_token)
				case "2":
					rt = getFace(base64.StdEncoding.EncodeToString(file), url[i]["name"], url[i]["type"], bai.API_URL+url[i]["url"]+"?access_token="+config.Baidu_Access_token)
				case "3":
					rt = getImg(base64.StdEncoding.EncodeToString(file), url[i]["name"], url[i]["type"], bai.API_URL+url[i]["url"]+"?access_token="+config.Baidu_Access_token)
				}
			}
		}
	}
	return rt
}
