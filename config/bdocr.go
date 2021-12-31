package config

// 百度智能云 应用 token
var Baidu_Access_token string

var Baidu_risk_type = map[string]string{
	"normal":    "正常身份证",
	"copy":      "复印件",
	"temporary": "临时身份证",
	"screen":    "翻拍",
	"unknown":   "其他未知情况",
}

var Baidu_idcard_number_type = map[float64]string{
	-1: "身份证正面所有字段全为空",
	0:  "身份证证号不合法",
	1:  "身份证证号和性别、出生信息一致",
	2:  "身份证证号和性别、出生信息都不一致",
	3:  "身份证证号和出生信息不一致",
	4:  "身份证证号和性别信息不一致",
}

var Baidu_bank_card_type = map[float64]string{
	0: "无法识别卡片类型",
	1: "借记卡",
	2: "贷记卡（原信用卡大部分为贷记卡）",
	3: "准贷记卡",
	4: "预付费卡",
}

// 百度智能云 应用
var OCR_URL = []map[string]interface{}{
	{
		"pid":      1,
		"ocrid":    1,
		"title":    "通用文字识别（标准版）",
		"quantity": 50000,
		"url":      "/general_basic",
	},
	{
		"pid":      1,
		"ocrid":    2,
		"title":    "通用文字识别（标准含位置版）",
		"quantity": 500,
		"url":      "/general",
	},
	{
		"pid":      1,
		"ocrid":    3,
		"title":    "通用文字识别（高精度版）",
		"quantity": 500,
		"url":      "/accurate_basic",
	},
	{
		"pid":      1,
		"ocrid":    4,
		"title":    "通用文字识别（高精度含位置版）",
		"quantity": 50,
		"url":      "/accurate",
	},
	{
		"pid":      1,
		"ocrid":    5,
		"title":    "网络图片文字识别",
		"quantity": 500,
		"url":      "/webimage",
	},
	{
		"pid":      1,
		"ocrid":    6,
		"title":    "数字识别",
		"quantity": 200,
		"url":      "/numbers",
	},
	{
		"pid":      1,
		"ocrid":    7,
		"title":    "手写文字识别",
		"quantity": 50,
		"url":      "/handwriting",
	},
	{
		"pid":      2,
		"ocrid":    1,
		"title":    "身份证识别",
		"quantity": 500,
		"url":      "/idcard",
	},
	{
		"pid":      2,
		"ocrid":    2,
		"title":    "银行卡识别",
		"quantity": 500,
		"url":      "/bankcard",
	},
	{
		"pid":      2,
		"ocrid":    3,
		"title":    "营业执照识别",
		"quantity": 200,
		"url":      "/business_license",
	},
	{
		"pid":      2,
		"ocrid":    4,
		"title":    "名片识别",
		"quantity": 500,
		"url":      "/business_card",
	},
	{
		"pid":      3,
		"ocrid":    1,
		"title":    "驾驶证识别",
		"quantity": 200,
		"url":      "/driving_license",
	},
	{
		"pid":      3,
		"ocrid":    2,
		"title":    "行驶证识别",
		"quantity": 200,
		"url":      "/vehicle_license",
	},
	{
		"pid":      3,
		"ocrid":    3,
		"title":    "车牌识别",
		"quantity": 200,
		"url":      "/license_plate",
	},
	{
		"pid":      4,
		"ocrid":    1,
		"title":    "通用票据识别",
		"quantity": 200,
		"url":      "/receipt",
	},
	{
		"pid":      4,
		"ocrid":    2,
		"title":    "增值税发票识别",
		"quantity": 500,
		"url":      "/vat_invoice",
	},
	{
		"pid":      4,
		"ocrid":    3,
		"title":    "火车票识别",
		"quantity": 50,
		"url":      "/train_ticket",
	},
	{
		"pid":      4,
		"ocrid":    4,
		"title":    "出租车票识别",
		"quantity": 50,
		"url":      "/taxi_receipt",
	},
	{
		"pid":      4,
		"ocrid":    5,
		"title":    "定额发票识别",
		"quantity": 500,
		"url":      "/quota_invoice",
	},
	{
		"pid":      5,
		"ocrid":    1,
		"title":    "印章识别",
		"quantity": 100,
		"url":      "/seal",
	},
	{
		"pid":      5,
		"ocrid":    2,
		"title":    "通信行程卡识别",
		"quantity": -1,
		"url":      "/travel_card",
	},
}
