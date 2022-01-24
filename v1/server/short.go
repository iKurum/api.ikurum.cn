package v1

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"api.ikurum.cn/global"
	"api.ikurum.cn/route"
	"api.ikurum.cn/util/logs"
)

func init() {
	route.POST("/short", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg, _ := json.Marshal(global.NewResult(&global.Result{
					Code: 0,
					Msg:  "recover:" + fmt.Sprint(err),
				}))
				rw.Write(msg)
			}
		}()

		r.ParseMultipartForm(32 << 20)
		// // 根据请求body创建一个json解析器实例
		// decoder := json.NewDecoder(r.Body)
		// // 用于存放参数key=value数据
		// var params map[string]string
		// // 解析参数 存入map
		// decoder.Decode(&params)
		fmt.Println("params:", r.Form)

		DB := global.OpenDB()

		var (
			msg     []byte
			errData string
			s       string
			err     error
		)

		URL := r.Form.Get("url")
		slength, e := strconv.Atoi(r.Form.Get("length"))
		if e != nil || slength > 4 || slength < 1 {
			slength = 4
		}

		if URL == "" {
			errData = "url 参数错误"
		} else {
			s, err = Transform(URL, slength, DB)
			if err != nil {
				logs.Warning("短链转换错误: ", err)
				errData = "短链转换错误"
			}
		}

		if errData != "" {
			msg, _ = json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  errData,
			}))
		} else if err != nil {
			msg, _ = json.Marshal(global.NewResult(&global.Result{
				Code: 0,
				Msg:  err,
			}))
		} else {
			msg, _ = json.Marshal(global.NewResult(&global.Result{
				Code:  200,
				Count: 1,
				Data:  "https://s-url.ikurum.cn/" + s,
			}))
		}

		rw.Write(msg)
	})
}

const (
	VAL   = 0x3FFFFFFF
	INDEX = 0x0000003D
)

var (
	alphabet = []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ~_+-*")
)

func Transform(longURL string, length int, DB *sql.DB) (string, error) {
	var sr string
	err := DB.QueryRow("select surl from short where url=?", longURL).Scan(&sr)
	global.CheckErr(err)
	if err == nil && sr != "" {
		return sr, nil
	}

	r, err := getShortURL(longURL, length, DB)
	if err == nil {
		setDBShort(DB, longURL, r)
	}
	return r, nil
}

func getShortURL(longURL string, length int, DB *sql.DB) (string, error) {
	md5Str := getMd5Str(longURL)
	var (
		tempVal int64
		result  []string
		tempUri []byte
		err     error
		sr      string
		r       string
	)

	for i := 0; i < 4; i++ {
		tempSubStr := md5Str[i*8 : (i+1)*8]
		hexVal, err := strconv.ParseInt(tempSubStr, 16, 64)
		if err != nil {
			return "", err
		}

		tempVal = int64(VAL) & hexVal
		var index int64
		tempUri = []byte{}
		for i := 0; i < 6; i++ {
			index = INDEX & tempVal
			tempUri = append(tempUri, alphabet[index])
			tempVal = tempVal >> 5
		}
		result = append(result, string(tempUri))
	}

	for i := 0; i < length; i++ {
		if sr == "" {
			sr += result[i]
		} else {
			sr += "-" + result[i]
		}
	}

	err = DB.QueryRow("select url from short where surl=?", sr).Scan(&r)
	global.CheckErr(err)

	if r != "" {
		getShortURL(longURL+longURL, length, DB)
	}

	return sr, nil
}

func setDBShort(DB *sql.DB, longURL string, shortURL string) {
	var (
		sql *sql.Stmt
		row int64
		err error
	)
	sql, err = DB.Prepare("insert into short(surl, url) values(?, ?)")
	global.CheckErr(err)

	res, errDB := sql.Exec(
		shortURL,
		longURL,
	)
	global.CheckErr(errDB, "insert short failed")
	err = errDB

	//查询影响的行数，判断修改插入成功
	row, err = res.RowsAffected()
	global.CheckErr(err, "insert rows failed")

	if row != 0 {
		logs.Info("insert short succ: ", row)
	} else {
		logs.Info("无修改")
	}
}

func getMd5Str(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	c := m.Sum(nil)
	return hex.EncodeToString(c)
}
