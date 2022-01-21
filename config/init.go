package config

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"api.ikurum.cn/util/logs"
	"gopkg.in/yaml.v3"
)

var (
	Online       bool
	SetTokenTime time.Duration
	DB           = map[string]string{}
)

type ConfigYaml struct {
	DB        DBYaml `yaml:"db"`
	Online    bool   `yaml:"online"`
	TokenTime string `yaml:"tokenTime"`
}

type DBYaml struct {
	Title    string `yaml:"title"`
	Ip       string `yaml:"ip"`
	User     string `yaml:"user"`
	Pw       string `yaml:"pw"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func GetApiConfig(file ...string) {
	logs.Warning("读取配置文件")

	url := strings.Join(file, "")
	if url == "" {
		url = "ApiConfig.yaml"
	}

	yamlFile, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Warning(err.Error())
	}

	var _config ConfigYaml
	err = yaml.Unmarshal(yamlFile, &_config)
	if err != nil {
		logs.Warning(err.Error())
	}

	Online = _config.Online

	s := strings.Split(_config.TokenTime, "_")
	if len(s) == 2 {
		t, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			logs.Warning("token 刷新时间设置错误")
		} else {
			var p time.Duration = 0
			switch s[1] {
			case "d":
				p = time.Hour * 24
			case "h":
				p = time.Hour
			case "m":
				p = time.Hour / 60
			}

			if p == 0 {
				logs.Warning("token 刷新时间单位设置错误")
			} else {
				SetTokenTime = time.Duration(t) * p
			}
		}
	}

	DB["title"] = _config.DB.Title
	DB["user"] = _config.DB.User
	DB["pw"] = _config.DB.Pw
	DB["port"] = _config.DB.Port
	DB["database"] = _config.DB.Database

	if _config.Online {
		DB["ip"] = "172.17.0.1"
	} else {
		DB["ip"] = _config.DB.Ip
	}
}
