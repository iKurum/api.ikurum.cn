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
	global.CheckErr(err)

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
func Read_file(file []byte, first int, second int) (interface{}, error) {
	DB := global.OpenDB()
	var url string
	err := DB.QueryRow("select url from bdocr where pid=? and ocrid=?", first, second).Scan(&url)
	t := global.CheckErr(err)
	if t == global.NoRows {
		return "", fmt.Errorf("类别错误")
	}

	var quantity int64 = 0
	err = DB.QueryRow("select quantity from bdocr where pid=? and ocrid=?", first, second).Scan(&quantity)
	global.CheckErr(err)
	if quantity <= 0 {
		return "", fmt.Errorf("今日次数已耗尽")
	}

	if err := fetch_token(); err == nil {
		rt, e := getTxt(base64.StdEncoding.EncodeToString(file), fmt.Sprint(first, ".", second), bai.API_URL+url+"?access_token="+config.Baidu_Access_token)

		sql, err := DB.Prepare("UPDATE bdocr SET quantity=? WHERE pid=? and ocrid=?")
		global.CheckErr(err)
		res, err := sql.Exec(quantity-1, first, second)
		global.CheckErr(err, "exec failed")

		//查询影响的行数，判断修改插入成功
		row, err := res.RowsAffected()
		global.CheckErr(err, "rows failed")
		fmt.Println("update bdocr quantity succ:", row)

		return rt, e
	}
	return "", fmt.Errorf("fetch token错误")
}
