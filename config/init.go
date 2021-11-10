package config

import "log"

func init() {
	log.Println("init mysql ip")
	if Online {
		DB["ip"] = "127.0.0.1"
	} else {
		DB["ip"] = "/* 远端ip */"
	}
}

// true		打包后，连接数据库
// false	本地启项目，连接远端数据库
var Online = false

// 数据库连接信息
var DB = map[string]string{
	"title":    "mysql",
	"user":     "",
	"pw":       "",
	"port":     "",
	"database": "",
	"ip":       "", // 在init中配置
}

// 百度智能云 应用 token
var Baidu_Access_token string

// 百度智能云 应用
var OCR_URL = []map[string]string{
	{
		"name": "通用文字（标准版）",
		"url":  "/ocr/v1/general_basic",
		"type": "1",
	},
	{
		"name": "身份证",
		"url":  "/ocr/v1/idcard",
		"type": "2",
	},
	{
		"name": "营业执照",
		"url":  "/ocr/v1/business_license",
		"type": "4",
	},
	{
		"name": "火车票",
		"url":  "/ocr/v1/train_ticket",
		"type": "8",
	},
	{
		"name": "通用票据",
		"url":  "/ocr/v1/receipt",
		"type": "16",
	},
	{
		"name": "驾驶证",
		"url":  "/ocr/v1/driving_license",
		"type": "32",
	},
	{
		"name": "车牌",
		"url":  "/ocr/v1/license_plate",
		"type": "64",
	},
}
var FACE_URL = []map[string]string{
	{
		"name": "人脸检测",
		"url":  "/face/v3/detect",
		"type": "1",
	},
	{
		"name": "在线活体检测",
		"url":  "/face/v3/faceverify",
		"type": "2",
	},
	{
		"name": "人脸属性编辑",
		"url":  "/face/v1/editattr",
		"type": "4",
	},
}
var IMAGE_URL = []map[string]string{
	{
		"name": "通用物体和场景识别",
		"url":  "/image-classify/v2/advanced_general",
		"type": "1",
	},
	{
		"name": "植物识别",
		"url":  "/image-classify/v1/plant",
		"type": "2",
	},
	{
		"name": "logo识别",
		"url":  "/image-classify/v2/logo",
		"type": "4",
	},
	{
		"name": "相似图片",
		"url":  "/image-classify/v1/realtime_search/similar/search",
		"type": "8",
	},
}
