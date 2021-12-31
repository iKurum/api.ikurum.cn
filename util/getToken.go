package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"api.ikurum.cn/config"
	"api.ikurum.cn/global"
	"api.ikurum.cn/util/logs"
)

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
			logs.Info("There is not row")
		} else {
			logs.Exit(err)
		}
	}
	getAccessToken()

	it := &intervalTime{
		interval: config.SetTokenTime,
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
	logs.Info("update refresh succ: ", row)

	// 更新photo
	global.GetBody("/photos/120x120/$value", "img")

	var one_count int64 = 0
	DB.QueryRow("select count(*) from one").Scan(&one_count)
	if one_count < 1 {
		global.SetOne()
		logs.Info("初始化一言")
	}

	var bd_count int64 = 0
	DB.QueryRow("select count(*) from bdocr").Scan(&bd_count)
	if bd_count < 1 {
		global.SetBd()
		logs.Info("初始化百度智能云接口")
	}

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
			logs.Info("update user email succ: ", row)
		}
		if name == "surname" {
			sql, err := DB.Prepare("UPDATE user SET name=? WHERE uid=1")
			global.CheckErr(err, "")
			res, err := sql.Exec(value.(string))
			global.CheckErr(err, "exec failed")

			//查询影响的行数，判断修改插入成功
			row, err := res.RowsAffected()
			global.CheckErr(err, "rows failed")
			logs.Info("update user name succ: ", row)
		}
	}

	// 更新detail
	getDetail()
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

// 更新文章
func getDetail() {
	var jsonTxt []byte

	logs.Warning("开始检查文章更新 ...")
	jsonTxt = global.GetBody("/drive/root:/article:/children?$top=100000", "")

	ch := make(chan map[string]interface{}, 1)

	var j map[string]interface{}
	json.Unmarshal(jsonTxt, &j)
	logs.Info("获取列表")
	for n, v := range j {
		if n == "value" {
			// 开始文章更新
			data := v.([]interface{})
			logs.Info("解析数据")

			for i := 0; i < len(data); i++ {
				da := data[i].(map[string]interface{})

				for name := range da {
					if name == "@microsoft.graph.downloadUrl" {
						logs.Info("检查文章状态 ", da["name"])
						reg := regexp.MustCompile(`[A-Za-z]`)
						l := strings.TrimSpace(reg.ReplaceAllString(da["lastModifiedDateTime"].(string), " "))
						c := strings.TrimSpace(reg.ReplaceAllString(da["createdDateTime"].(string), " "))
						time_last, _ := time.ParseInLocation("2006-01-02 15:04:05", l, time.Local)
						time_create, _ := time.ParseInLocation("2006-01-02 15:04:05", c, time.Local)
						da["lastModifiedDateTime"] = time_last.Unix() * 1000
						da["createdDateTime"] = time_create.Unix() * 1000
						ch <- da
						setDetail(ch)
					}
				}
			}
		}
	}
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
		// 获取数据库 文章更新时间
		e = DB.QueryRow("select uptime, addtime from essay where essayId=?", da["id"].(string)).Scan(
			&time_last,
			&time_create,
		)
		global.CheckErr(e, "")

		if da["lastModifiedDateTime"].(int64) != time_last {
			e = fmt.Errorf("essay detail has change")
		}
	}

	// 创建或更新essay detail
	if md := getMD(da["@microsoft.graph.downloadUrl"].(string), da["id"].(string)); md {
		f, err := ioutil.ReadFile("md/" + da["id"].(string) + ".md")
		if err != nil {
			logs.Exit(err)
		}

		toSetDetail(e, string(f), da)
	}
}

// 存储文章详情
func toSetDetail(e error, f string, data map[string]interface{}) {
	DB := global.OpenDB()
	reg := regexp.MustCompile(`^<!-- config {[\s\S]*} -->`)
	resf := reg.FindAllStringSubmatch(f, -1)

	var (
		archive  string = ""
		sql_arch string
	)
	if len(resf) != 0 && len(resf[0]) != 0 {
		f = strings.Replace(f, resf[0][0], "", 1)

		archive = strings.Replace(resf[0][0], "<!-- config {", "", 1)
		archive = strings.Replace(archive, "} -->", "", 1)
		archive = strings.ReplaceAll(archive, " ", "")
	}

	if archive != "" {
		archa := strings.Split(archive, "\r\n")
		for i := 0; i < len(archa); i++ {
			a := strings.Split(archa[i], ":")

			if a[0] == "archive" && a[1] != "" {
				sql_arch = strings.ReplaceAll(a[1], "\"\"", ",")
				sql_arch = strings.ReplaceAll(sql_arch, "\"", "")
				sql_arch = strings.ReplaceAll(sql_arch, "_", " ")
			}
		}
	}

	a := strings.Split(f, "<!-- more -->")

	data["name"] = strings.Split(data["name"].(string), ".md")[0]

	if e == sql.ErrNoRows {
		// insert
		logs.Warning("insert essay ", data["name"])
		sql, err := DB.Prepare("insert into essay(essayId, title, size, content, note, archive, uptime, addtime)values(?,?,?,?,?,?,?)")
		global.CheckErr(err, "")
		res, err := sql.Exec(
			data["id"],
			data["name"],
			data["size"],
			f,
			a[0],
			sql_arch,
			data["lastModifiedDateTime"],
			data["createdDateTime"],
		)
		global.CheckErr(err, "insert exec failed")

		//查询影响的行数，判断修改插入成功
		row, err := res.RowsAffected()
		global.CheckErr(err, "insert rows failed")
		logs.Info("insert essay succ:", row, data["name"])
	} else {
		// update
		if e != nil {
			logs.Warning("update essay ", data["name"])
		}
		sql, err := DB.Prepare("update essay set title=?, size=?, content=?, note=?, archive=?, uptime=?, addtime=? where essayId=?")
		global.CheckErr(err, "")
		res, err := sql.Exec(
			data["name"],
			data["size"],
			f,
			a[0],
			sql_arch,
			data["lastModifiedDateTime"],
			data["createdDateTime"],
			data["id"],
		)
		global.CheckErr(err, "update exec failed")

		//查询影响的行数，判断修改插入成功
		row, err := res.RowsAffected()
		global.CheckErr(err, "update rows failed")
		if row != 0 {
			logs.Info("update essay succ: ", row, " ", data["name"])
		}
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
