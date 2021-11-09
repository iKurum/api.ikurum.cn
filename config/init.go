package config

var DB = map[string]string{
	"title":    "", // 数据库名称，如mysql
	"user":     "",
	"pw":       "",
	"ip":       "",
	"port":     "",
	"database": "",
}

var Baidu_Access_token string

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
