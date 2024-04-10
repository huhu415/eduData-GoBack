package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 版本信息version information
var (
	version   = ""
	buildDate = ""
	gitCommit = ""
)

type Config struct {
	// ocr
	BaiduRequestUrl  string `mapstructure:"baidu_request_url"`
	BaiduAccesstoken string `mapstructure:"baidu_accessToken"`
	JfymRequestUrl   string `mapstructure:"jfym_request_url"`
	JfymToken        string `mapstructure:"jfym_token"`

	// extra
	PgConfig   string `mapstructure:"pg_config"`
	JwtKey     string `mapstructure:"jwt_key"`
	ListenPort string `mapstructure:"listen_port"`
	UserAgent  string `mapstructure:"user_agent"`
}

var C Config

func Loadconfig() {
	log.Println("Init from config file")
	// Default config
	viper.SetDefault("baidu_request_url", "https://aip.baidubce.com/rest/2.0/ocr/v1/numbers")
	viper.SetDefault("jfym_request_url", "http://www.jfbym.com/api/YmServer/customApi")
	viper.SetDefault("user_agent", "Mozilla/5.0 (Macintosh Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	viper.SetDefault("jwt_key", "9385g0x98n347tx980y34g9sfgsldkjvilr")

	viper.SetEnvPrefix("edu")
	viper.AutomaticEnv()

	parseFlag()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	viper.AddConfigPath(pathAbs)
	viper.AddConfigPath("/Users/hello/Library/Mobile Documents/com~apple~CloudDocs/代码项目/eduData-GoBack")
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("read config file error: %+v\n", err)
		return
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Printf("unmarshal config file error: %+v\n", err)
		return
	}
	log.Println(viper.GetString("listen_port"))
	log.Println("read config file success")
	return
}

func parseFlag() {
	pflag.StringP("configFile", "c", "config.yaml", "config file")
	pflag.StringP("listen_port", "l", "8080", "listen address")
	pflag.BoolP("version", "v", false, "version information")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	if viper.GetBool("version") {
		fmt.Println("version:", version)
		fmt.Println("buildDate:", buildDate)
		fmt.Println("gitCommit:", gitCommit)
		os.Exit(0)
	}
}
