package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"api.ikurum.cn/global"
)

type globalConfig struct {
	CLIENT_ID     string
	CLIENT_SECRET string
	BASE_URL      string
	refresh       string
}

var row globalConfig

// 初始化token，设置定时任务
func StartToken() {
	DB := global.OpenDB()

	err := DB.QueryRow("select CLIENT_ID,CLIENT_SECRET from global where gid=1").Scan(&row.CLIENT_ID, &row.CLIENT_SECRET)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("There is not row")
		} else {
			log.Fatalln(err)
		}
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

func getAccessToken() {
	DB := global.OpenDB()

	err := DB.QueryRow("select refresh from user where uid=1").Scan(&row.refresh)
	global.CheckErr(err, "")

	accessToken, refreshToken := getToken()

	// 执行SQL语句
	sql, err := DB.Prepare("UPDATE user SET access=?, refresh=?, uptime=? WHERE uid=1")
	global.CheckErr(err, "")
	res, err := sql.Exec(accessToken, refreshToken, time.Now().Unix()*1000)
	global.CheckErr(err, "exec failed")

	//查询影响的行数，判断修改插入成功
	row, err := res.RowsAffected()
	global.CheckErr(err, "rows failed")
	fmt.Println("update refresh succ:", row)

	// 更新photo
	global.GetBody("/photos/120x120/$value", "img")

	// 更新detail
	getDetail()

	// 更新info
	jsonTxt := global.GetBody("/", "")
	var j map[string]interface{}
	json.Unmarshal(jsonTxt, &j)

	for name, value := range j {
		// 执行SQL语句
		if name == "mail" {
			sql, err := DB.Prepare("UPDATE user SET email=? WHERE uid=1")
			global.CheckErr(err, "")
			res, err := sql.Exec(value.(string))
			global.CheckErr(err, "exec failed")

			//查询影响的行数，判断修改插入成功
			row, err := res.RowsAffected()
			global.CheckErr(err, "rows failed")
			fmt.Println("update user email succ:", row)
		}
		if name == "surname" {
			sql, err := DB.Prepare("UPDATE user SET name=? WHERE uid=1")
			global.CheckErr(err, "")
			res, err := sql.Exec(value.(string))
			global.CheckErr(err, "exec failed")

			//查询影响的行数，判断修改插入成功
			row, err := res.RowsAffected()
			global.CheckErr(err, "rows failed")
			fmt.Println("update user name succ:", row)
		}
	}
}

// 刷新token
func getToken() (string, string) {
	urlStr := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": []string{row.refresh},
		"client_id":     []string{row.CLIENT_ID},
		"client_secret": []string{row.CLIENT_SECRET},
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

	return accessToken.(string), refreshToken.(string)
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

// 更新文章
func getDetail() {
	var jsonTxt []byte

	fmt.Println("开始检查文章更新 ...")
	jsonTxt = global.GetBody("/drive/root:/article:/children?$top=100000", "")

	ch := make(chan map[string]interface{}, 6)

	var j map[string]interface{}
	json.Unmarshal(jsonTxt, &j)

	for n, v := range j {
		if n == "value" {
			data := v.([]interface{})
			for i := 0; i < len(data); i++ {
				da := data[i].(map[string]interface{})
				for name := range da {
					if name == "@microsoft.graph.downloadUrl" {
						reg := regexp.MustCompile(`[A-Za-z]`)
						l := strings.TrimSpace(reg.ReplaceAllString(da["lastModifiedDateTime"].(string), " "))
						c := strings.TrimSpace(reg.ReplaceAllString(da["createdDateTime"].(string), " "))
						time_last, _ := time.ParseInLocation("2006-01-02 15:04:05", l, time.Local)
						time_create, _ := time.ParseInLocation("2006-01-02 15:04:05", c, time.Local)
						da["lastModifiedDateTime"] = time_last.Unix() * 1000
						da["createdDateTime"] = time_create.Unix() * 1000
						ch <- da
						go setDetail(ch)
					}
				}
			}
		}
	}
	log.Println("文章更新完成")
}

// 检查文章状态
func setDetail(ch chan map[string]interface{}) {
	DB := global.OpenDB()
	da := <-ch

	var e error = nil
	if e = global.HasEssay(da["id"].(string)); e == nil {
		var (
			time_create int64
			time_last   int64
		)
		e = DB.QueryRow("select uptime, addtime from essay where essayId=?", da["id"].(string)).Scan(
			&time_last,
			&time_create,
		)
		global.CheckErr(e, "")

		if da["lastModifiedDateTime"].(int64) != time_last ||
			da["createdDateTime"].(int64) != time_create {
			e = fmt.Errorf("essay detail has new")
		}
	}

	if e != nil {
		// 创建或更新essay detail
		if md := getMD(da["@microsoft.graph.downloadUrl"].(string), da["id"].(string)); md {
			f, err := ioutil.ReadFile("md/" + da["id"].(string) + ".md")
			if err != nil {
				log.Fatal(err)
			}

			toSetDetail(e, string(f), da)
		}
	}
}

// 存储文章详情
func toSetDetail(e error, f string, data map[string]interface{}) {
	DB := global.OpenDB()

	a := strings.Split(f, "<!-- more -->")

	if e == sql.ErrNoRows {
		// insert
		log.Println("insert essay", data["name"])
		sql, err := DB.Prepare("insert into essay(essayId, title, size, content, note, uptime, addtime)values(?,?,?,?,?,?,?)")
		global.CheckErr(err, "")
		res, err := sql.Exec(
			data["id"],
			data["name"],
			data["size"],
			f,
			a[0],
			data["lastModifiedDateTime"],
			data["createdDateTime"],
		)
		global.CheckErr(err, "insert exec failed")

		//查询影响的行数，判断修改插入成功
		row, err := res.RowsAffected()
		global.CheckErr(err, "insert rows failed")
		fmt.Println("insert essay succ:", row, data["id"])
	} else {
		// update
		log.Println("update essay", data["name"])
		sql, err := DB.Prepare("update essay set title=?, size=?, content=?, note=?, uptime=?, addtime=? where essayId=?")
		global.CheckErr(err, "")
		res, err := sql.Exec(
			data["name"],
			data["size"],
			f,
			a[0],
			data["lastModifiedDateTime"],
			data["createdDateTime"],
			data["id"],
		)
		global.CheckErr(err, "update exec failed")

		//查询影响的行数，判断修改插入成功
		row, err := res.RowsAffected()
		global.CheckErr(err, "update rows failed")
		fmt.Println("update essay succ:", row, data["id"])
	}
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
	io.Copy(f, resp.Body)

	return true
}
