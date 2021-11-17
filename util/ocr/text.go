package ocr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func getTxt(f string, t string, image_url string) (map[string]interface{}, error) {
	var params = make(url.Values)
	params.Add("image", f)
	if t == "2.1" {
		params.Add("id_card_side", "back")
	}
	post_data := params.Encode()

	req, _ := http.NewRequest("POST", image_url, strings.NewReader(post_data))
	req.Header.Set("Content-Type", "application/application/x-www-form-urlencoded")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var (
		j    map[string]interface{}
		data       = make(map[string]interface{})
		e    error = nil
	)
	jsonTxt, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(jsonTxt, &j)

	if j["error_code"] != nil {
		delete(j, "log_id")
		delete(j, "words_result")
		delete(j, "words_result_num")
		e = fmt.Errorf(j["error_message"].(string))
	}
	if t == "2.1" && j["image_status"] != "normal" {
		delete(j, "log_id")
		delete(j, "words_result")
		delete(j, "words_result_num")
		e = fmt.Errorf(j["image_status"].(string))
	}

	if j["words_result"] != nil {
		data["result"] = j["words_result"]
		data["resultNum"] = j["words_result_num"]
		if t == "2.1" {
			data["type"] = j["idcard_number_type"]
		}
	} else if j["result"] != nil {
		data["result"] = j["result"]
	}

	if e != nil {
		return j, e
	}
	return data, e
}
