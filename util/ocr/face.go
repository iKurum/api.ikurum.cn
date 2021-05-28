package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getFace(f string, name string, t string, image_url string) []byte {
	var jsonStr string
	switch t {
	case "1":
		jsonStr = fmt.Sprintf(`{"image":"%v","image_type": "BASE64"}`, f)
	case "2":
		jsonStr = fmt.Sprintf(`[{"image":"%v","image_type": "BASE64"}]`, f)
	case "4":
		jsonStr = fmt.Sprintf(`{"image":"%v","image_type": "BASE64", "action_type":"TO_KID","quality_control":"HIGH"}`, f)
	}

	// 调用人脸识别服务
	fmt.Println("连接 FACE 服务 ...", name)

	req, _ := http.NewRequest("POST", image_url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	jsonTxt, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(jsonTxt, &j)

	if j["error_code"].(float64) == 0 {
		if t == "4" {
			var d = fmt.Sprintf("img/output/%v.jpg", "TO_KID")
			fmt.Println("写入图片:", d)
			write_for_base64(d, j["result"].(map[string]interface{})["image"].(string))
		} else {
			fmt.Printf("%v\n", j["result"])
		}
	} else {
		log.Fatal(j["error_msg"])
	}

	return []byte("")
}

// 写出图片OCR
func write_for_base64(filename string, data string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666) //打开文件
	if err != nil {
		panic(err)
	}

	decodeBytes, err1 := base64.StdEncoding.DecodeString(data)
	if err1 != nil {
		panic(err)
	}

	n, err2 := io.WriteString(f, string(decodeBytes)) //写入文件(字符串)
	if err2 != nil {
		panic(err)
	}
	fmt.Printf("写入 %d 个字节\n", n)
}
