package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"api.ikurum.cn/global"
)

var clientID []string
var clientSecret []string

func getToken(refresh string) string {
	urlStr := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refresh},
		"client_id":     clientID,
		"client_secret": clientSecret,
		"scope":         {"user.read"},
		"redirect_uri":  {"http://localhost:53682/"},
	}

	resp, err := http.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	jsonTxt, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(jsonTxt, &j)

	refreshToken := j["refresh_token"]
	accessToken := j["access_token"]

	fmt.Println("refresh token:", refreshToken)

	if global.Online {
		global.UpdateByDB("global", "refreshTokenOnline", refreshToken.(string))
	} else {
		global.UpdateByDB("global", "refreshTokenTest", refreshToken.(string))
	}

	return accessToken.(string)
}

func getAccessToken() {
	var refresh string
	if global.Online {
		refresh = global.GetByDB("global", "refreshTokenOnline")
	} else {
		refresh = global.GetByDB("global", "refreshTokenTest")
	}

	err := global.UpdateByDB("global", "accessToken", getToken(refresh))
	if err == nil {
		// 更新detail
		getDetail()

		// 更新photo
		global.GetBody("/photos/120x120/$value", "img")

		// 更新info
		jsonTxt := global.GetBody("/", "")

		var j map[string]interface{}
		json.Unmarshal(jsonTxt, &j)

		for name, value := range j {
			if name == "mail" || name == "surname" {
				global.UpdateByDB("userInfo", name, value.(string))
			}
		}
	}
}

type intervalTime struct {
	interval time.Duration
	job      func()
	enabled  bool
	wg       sync.WaitGroup
}

func (it *intervalTime) isr() {
	if it.enabled {
		it.job()
		time.AfterFunc(it.interval, it.isr)
	} else {
		it.wg.Done()
	}
}

func (it *intervalTime) start() {
	if it.enabled {
		it.wg.Add(1)
		time.AfterFunc(it.interval, it.isr)
	}
}

// 更新 token
func StartToken() {
	clientID = []string{global.GetByDB("global", "clientID")}
	if global.Online {
		clientSecret = []string{global.GetByDB("global", "clientSecretOnline")}
	} else {
		clientSecret = []string{global.GetByDB("global", "clientSecretTest")}
		fmt.Println("测试环境")
	}

	getAccessToken()

	it := &intervalTime{
		interval: time.Duration(global.SetTokenTime) * time.Hour,
		job:      getAccessToken,
		enabled:  true,
	}

	it.start()
	it.wg.Wait()
}

type detailID struct {
	id              string
	createdDateTime int64
}

// 文章id map
var arr []detailID

// 更新文章
func getDetail() {
	var jsonTxt []byte

	fmt.Println("开始检查文章更新 ...")
	arr = arr[0:0]
	global.DelBucket("detailList")
	jsonTxt = global.GetBody("/drive/root:/article:/children?$top=100000", "")

	ch := make(chan string, 10)

	var j map[string]interface{}
	json.Unmarshal(jsonTxt, &j)

	for n, v := range j {
		if n == "value" {
			data := v.([]interface{})
			for i := 0; i < len(data); i++ {
				da := data[i].(map[string]interface{})
				for name, value := range da {
					if name == "@microsoft.graph.downloadUrl" {
						reg := regexp.MustCompile(`[A-Za-z]`)
						c := strings.TrimSpace(reg.ReplaceAllString(da["createdDateTime"].(string), " "))
						time, _ := time.ParseInLocation("2006-01-02 15:04:05", c, time.Local)
						arr = append(arr, detailID{da["id"].(string), time.Unix()})

						ch <- value.(string)
						go setDetail(da, ch)
					}
				}
			}
		}
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].createdDateTime > arr[j].createdDateTime // 降序
		// return arr[i].createdDateTime < arr[j].createdDateTime // 升序
	})

	var a = ""
	for i := 0; i < len(arr); i++ {
		a += arr[i].id + ","
	}

	global.UpdateByDB("detailList", "id", a)
	fmt.Println("文章 id:", arr)
}

// 检查文章状态
func setDetail(da map[string]interface{}, ch chan string) {
	var e error = nil

	if e = global.HasBucket(da["id"].(string)); e == nil {
		if da["lastModifiedDateTime"].(string) != global.GetByDB(da["id"].(string), "lastModifiedDateTime") {
			fmt.Printf("lastModifiedDateTime: %s\t%s\n", da["lastModifiedDateTime"].(string), global.GetByDB(da["id"].(string), "lastModifiedDateTime"))
			e = fmt.Errorf("detail has new")
			fmt.Println("更新bucket detail:", da["id"].(string))
		}
	} else {
		fmt.Println("创建bucket detail:", da["id"].(string))
	}

	if e != nil {
		// 首次创建detail bucket 或更新detail
		fmt.Println("detail nil ...", da["id"].(string))
		if md := getMD(da["@microsoft.graph.downloadUrl"].(string), da["id"].(string)); md {
			f, err := ioutil.ReadFile("md/" + da["id"].(string) + ".md")
			if err != nil {
				log.Fatal(err)
			}

			toSetDetail(string(f), da, ch)
		}
	}
}

// 存储文章详情
func toSetDetail(f string, data map[string]interface{}, ch chan string) {
	a := strings.Split(f, "<!-- more -->")

	for k, v := range data {
		if k == "lastModifiedDateTime" || k == "name" || k == "size" || k == "createdDateTime" || k == "id" {
			switch v.(type) {
			case float64:
				global.UpdateByDB(data["id"].(string), k, strconv.FormatFloat(v.(float64), 'f', -1, 64))

			case string:
				global.UpdateByDB(data["id"].(string), k, v.(string))
			}
		}
	}

	global.UpdateByDB(data["id"].(string), "article", f)
	global.UpdateByDB(data["id"].(string), "note", a[0])
	<-ch
}

// 获取文章详情
func getMD(url string, id string) bool {
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	f, err := os.Create("md/" + id + ".md")
	if err != nil {
		return false
	}
	fmt.Println("resp body:", resp.Body)
	io.Copy(f, resp.Body)

	return true
}
