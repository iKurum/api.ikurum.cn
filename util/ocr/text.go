package ocr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"api.ikurum.cn/config"
)

func getTxt(f string, t string, image_url string) (map[string]interface{}, error) {
	var params = make(url.Values)
	params.Add("image", f)
	if t == "2.1" {
		params.Add("id_card_side", "back")
		params.Add("detect_risk", "true")
		params.Add("detect_photo", "true")
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
	if t == "2.1" && j["image_status"] != "normal" && j["image_status"] != "reversed_side" {
		delete(j, "log_id")
		delete(j, "words_result")
		delete(j, "words_result_num")
		delete(j, "edit_tool")
		delete(j, "photo")
		e = fmt.Errorf(j["image_status"].(string))
	}

	if j["words_result"] != nil {
		data["result"] = j["words_result"]
		data["resultNum"] = j["words_result_num"]
		if t == "2.1" {
			// 身份证识别
			data["type"] = config.Baidu_idcard_number_type[j["idcard_number_type"].(float64)]
			data["edit"] = j["edit_tool"]
			data["photo"] = j["photo"]
			data["riskType"] = config.Baidu_risk_type[j["risk_type"].(string)]
		}
	} else if j["result"] != nil {
		if t == "2.2" {
			// 银行卡识别
			data["result"] = make(map[string]string)
			data["result"].(map[string]string)["type"] = config.Baidu_bank_card_type[j["result"].(map[string]interface{})["bank_card_type"].(float64)]
			data["result"].(map[string]string)["number"] = j["result"].(map[string]interface{})["bank_card_number"].(string)
			data["result"].(map[string]string)["bankName"] = j["result"].(map[string]interface{})["bank_name"].(string)
			data["result"].(map[string]string)["holderName"] = j["result"].(map[string]interface{})["holder_name"].(string)
			data["result"].(map[string]string)["date"] = j["result"].(map[string]interface{})["valid_date"].(string)
		} else {
			data["result"] = j["result"]
		}
	}

	if e != nil {
		return j, e
	}
	return data, e
}
