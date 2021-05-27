package global

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"api.ikurum.cn/initDB"
	"github.com/boltdb/bolt"
)

// 接口返回
type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	More bool        `json:"more"`
	Page int         `json:"page"`
	Size int         `json:"size"`
}

// 初始化数据库
func BoltInit() {
	db, err := bolt.Open("ikurum.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(t *bolt.Tx) error {
		var c = make(map[string]string)
		if initDB.DB["clientID"] != "" {
			c = initDB.DB
		} else {
			c = initdb
		}

		b, err := t.CreateBucketIfNotExists([]byte("global"))
		if err != nil {
			return fmt.Errorf("create global bucket: %s", err)
		}

		fmt.Println("set global ...")

		for k, v := range c {
			err = b.Put([]byte(k), []byte(v))
			if err != nil {
				return fmt.Errorf("create bucket global %s: %s", k, err)
			}
			fmt.Printf("create bucket global: %s\n", k)
		}

		return err
	})
}

// 是否已有bucket
func HasBucket(bucket string) error {
	db, err := bolt.Open("ikurum.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(t *bolt.Tx) error {
		b, e := t.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			b.Delete([]byte(bucket))
		}
		return e
	})
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

// 获取bucket数据
func GetByBucket(bucket string) map[string]string {
	db, err := bolt.Open("ikurum.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m := make(map[string]string)

	if err := db.Update(func(t *bolt.Tx) error {
		_, e := t.CreateBucketIfNotExists([]byte(bucket))
		return e
	}); err == nil {
		fmt.Printf("获取%s bucket key-value\n", bucket)
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(bucket))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				m[string(k)] = string(v)
			}

			return nil
		})
	} else {
		fmt.Printf("没有%s bucket\n", bucket)
	}

	return m
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
	baseURL := GetByDB("global", "baseURL")
	accessToken := GetByDB("global", "accessToken")

	req, _ := http.NewRequest("GET", baseURL+url, nil)
	req.Header.Set("Authorization", accessToken)

	if t == "" {
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
		err := UpdateByDB("photo", "str", string(jsonTxt))
		if err == nil {
			fmt.Println("更新头像完成")
		}
	}

	return jsonTxt
}

// func getPhoto(body io.Reader) error {
// 	f, err := os.Create("v1/user/photo.jpg")
// 	if err != nil {
// 		return err
// 	}
// 	io.Copy(f, body)
// 	return nil
// }

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
