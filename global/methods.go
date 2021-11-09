package global

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"database/sql"

	"api.ikurum.cn/config"
	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
)

// token刷新时间  \小时
var SetTokenTime int = 1

// 接口返回
type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	More bool        `json:"more"`
	Page int         `json:"page"`
	Size int         `json:"size"`
}

// 文章列表
type Essay_list struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Size   string `json:"size"`
	Note   string `json:"note"`
	Uptime int    `json:"upTime"`
}

// 文章详情
type Essay struct {
	Essay_list
	Content string `json:"content"`
	Err     string `json:"err"`
}

// 链接数据库
func OpenDB() *sql.DB {
	var c = make(map[string]string)
	if config.DB["ip"] != "" {
		c = config.DB
	} else {
		c = initdb
	}
	path := strings.Join([]string{c["user"], ":", c["pw"], "@tcp(", c["ip"], ":", c["port"], ")/", c["database"], "?charset=utf8"}, "")
	DB, _ := sql.Open(c["title"], path)
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

// 获取essay数据
func GetByEssay(essayId string) Essay {
	DB := OpenDB()

	var essay Essay
	err := DB.QueryRow("select aid,title,note,content,size,uptime from essay where aid=?", essayId).Scan(
		&essay.Id,
		&essay.Title,
		&essay.Note,
		&essay.Content,
		&essay.Size,
		&essay.Uptime,
	)
	t := CheckErr(err, "")
	if t == 1 {
		essay.Err = "参数id错误"
	}
	return essay
}

// 删除bucket
func DelBucket(bucket string) error {
	db, err := bolt.Open("ikurum.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(t *bolt.Tx) error {
		b, e := t.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			e = b.Delete([]byte(bucket))
		}
		return e
	})
}

// db 获取单个数据
func GetByDB(bucket string, k string) string {
	db, err := bolt.Open("ikurum.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var s string
	db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(bucket))
		s = string(b.Get([]byte(k))[:])
		return nil
	})

	return s
}

// db 更新单个数据
func UpdateByDB(bucket string, k string, v string) error {
	db, err := bolt.Open("ikurum.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create global bucket: %s", err)
		}

		err = b.Put([]byte(k), []byte(v))
		fmt.Printf("更新%s bucket数据: %s\n", bucket, k)
		return err
	})

	return err
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
		fmt.Println("更新头像完成")
	}

	return jsonTxt
}

// err 检查
func CheckErr(err error, str string) int {
	if err != nil {
		if err == sql.ErrNoRows {
			return 1
		} else {
			log.Fatalln(str, err)
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
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

// 设置 请求头
func SetHeader(rw http.ResponseWriter) {
	rw.Header().Add("x-content-type-options", "nosniff")
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json;utf-8")
}
