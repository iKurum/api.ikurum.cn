package global

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"database/sql"

	"api.ikurum.cn/config"
	"api.ikurum.cn/util/logs"
	_ "github.com/go-sql-driver/mysql"
)

// 接口返回
type Result struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   interface{} `json:"msg"`
	More  bool        `json:"more"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Count int64       `json:"count"`
}

// 文章列表
type Essay_list struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Size    string `json:"size"`
	Note    string `json:"note"`
	Uptime  int64  `json:"upTime"`
	Addtime int64  `json:"addTime"`
	Archive string `json:"archive"`
}

// 文章详情
type Essay struct {
	Essay_list
	Content       string `json:"content"`
	Next          int    `json:"next"`
	Nexttitle     string `json:"nextTitle"`
	Previous      int    `json:"previous"`
	Previoustitle string `json:"previousTitle"`
	Err           string `json:"err"`
}

// 链接数据库
func OpenDB() *sql.DB {
	path := strings.Join([]string{config.DB["user"], ":", config.DB["pw"], "@tcp(", config.DB["ip"], ":", config.DB["port"], ")/", config.DB["database"], "?charset=utf8"}, "")
	DB, _ := sql.Open(config.DB["title"], path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)

	//验证连接
	err := DB.Ping()
	CheckErr(err, "open database fail")
	return DB
}

// 是否已有essay
func HasEssay(essayId string) error {
	DB := OpenDB()

	var time int
	err := DB.QueryRow("select uptime from essay where essayId=?", essayId).Scan(&time)
	t := CheckErr(err, "")
	if t == 1 {
		return err
	}
	return nil
}

// 获取essay详情
func GetByEssay(essayId string) Essay {
	DB := OpenDB()

	var essay Essay
	err := DB.QueryRow("select aid,title,content,size,uptime,addtime,archive from essay where aid=?", essayId).Scan(
		&essay.Id,
		&essay.Title,
		&essay.Content,
		&essay.Size,
		&essay.Uptime,
		&essay.Addtime,
		&essay.Archive,
	)
	t := CheckErr(err, "")
	if t == 1 {
		essay.Err = "参数id错误"
	}

	//获取上一条 id
	err = DB.QueryRow("select aid,title from essay where addtime<? order by addtime desc", essay.Addtime).Scan(
		&essay.Previous,
		&essay.Previoustitle,
	)
	CheckErr(err, "")

	//获取下一条 id
	err = DB.QueryRow("select aid,title from essay where addtime>? order by addtime asc", essay.Addtime).Scan(
		&essay.Next,
		&essay.Nexttitle,
	)
	CheckErr(err, "")

	return essay
}

// 请求数据
func GetBody(url string, t string) []byte {
	var (
		baseURL string
		access  string
	)

	DB := OpenDB()
	err := DB.QueryRow("select access from user where uid=1").Scan(&access)
	CheckErr(err, "")

	err = DB.QueryRow("select BASE_URL from global where gid=1").Scan(&baseURL)
	CheckErr(err, "")
	// baseURL := GetByDB("global", "baseURL")
	// accessToken := GetByDB("global", "accessToken")

	req, _ := http.NewRequest("GET", baseURL+url, nil)
	req.Header.Set("Authorization", access)

	if t == "" || t == "json" {
		req.Header.Set("Content-Type", "application/json")
	} else if t == "img" {
		req.Header.Set("Content-Type", t)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var jsonTxt []byte
	jsonTxt, _ = ioutil.ReadAll(resp.Body)

	if t == "img" {
		sql, err := DB.Prepare("update user set photo=? where uid=1")
		CheckErr(err, "")

		res, err := sql.Exec(base64.StdEncoding.EncodeToString(jsonTxt))
		CheckErr(err, "exec failed")

		//查询影响的行数，判断修改插入成功
		_, err = res.RowsAffected()
		CheckErr(err, "rows failed")
		logs.Info("更新头像完成")
	}

	return jsonTxt
}

// err 检查
func CheckErr(err error, str string) int {
	if err != nil {
		if err == sql.ErrNoRows {
			return 1
		} else {
			logs.Exit(str, err)
		}
	}

	return 0
}

// 返回格式
func NewResult(res *Result) *Result {
	if res == nil {
		res = &Result{}
	}

	if res.Page == 0 {
		res.Page = 1
	}

	if res.Size == 0 {
		res.Size = 10
	}

	return res
}

// 数据库初始化 设置一言
func SetOne() {
	fin, err := os.OpenFile("./one", os.O_RDONLY, 0)
	if err != nil {
		logs.Exit(err)
	}
	defer fin.Close()

	DB := OpenDB()
	sql := "insert into one(md) values(?)"
	stmt, err := DB.Prepare(sql)
	CheckErr(err, "")
	defer stmt.Close()

	sc := bufio.NewScanner(fin)
	for sc.Scan() {
		t := sc.Text()
		stmt.Exec(t)
	}

	if err = sc.Err(); err != nil {
		logs.Exit(err)
	}
}

// 数据库初始化 设置百度智能云接口
func SetBd() {
	DB := OpenDB()
	sql := "insert into bdocr(pid,ocrid,title,quantity,url) values(?,?,?,?,?)"
	stmt, err := DB.Prepare(sql)
	CheckErr(err, "")
	defer stmt.Close()

	for i := 0; i < len(config.OCR_URL); i++ {
		stmt.Exec(config.OCR_URL[i]["pid"], config.OCR_URL[i]["ocrid"], config.OCR_URL[i]["title"], config.OCR_URL[i]["quantity"], config.OCR_URL[i]["url"])
	}
}
