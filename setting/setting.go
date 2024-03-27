package setting

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	// 百度ocr配置
	RequestUrl       string
	BaiduAccessToken string
	// jfbym的ocr配置
	CustomUrl string
	Token     string

	// 数据库配置
	Postgres string

	JwtKey   string
	HttpPort string

	UserAgent string
)

func init() {
	//初始化配置文件
	CofigIni, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件加载失败，请检查config.ini文件路径：", err)
	}
	loadServer(CofigIni)
	loadDBConfig(CofigIni)
	loadBaidu(CofigIni)
	loadUG(CofigIni)
	loadJfbym(CofigIni)
}

func loadBaidu(file *ini.File) {
	selectIni := file.Section("baidu")
	BaiduAccessToken = selectIni.Key("BaiduAccessToken").String()
	RequestUrl = selectIni.Key("RequestUrl").String()
}

func loadJfbym(file *ini.File) {
	selectIni := file.Section("jfbym")
	CustomUrl = selectIni.Key("CustomUrl").String()
	Token = selectIni.Key("Token").String()
}

func loadDBConfig(file *ini.File) {
	selectIni := file.Section("database")
	Postgres = selectIni.Key("Postgres").String()
}

func loadServer(file *ini.File) {
	selectIni := file.Section("server")
	HttpPort = selectIni.Key("HttpPort").MustString("8080")
	JwtKey = selectIni.Key("HmacSampleSecret").MustString("89js82js72")
}

func loadUG(file *ini.File) {
	selectIni := file.Section("personal")
	UserAgent = selectIni.Key("UserAgent").String()
}
