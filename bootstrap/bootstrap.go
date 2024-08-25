package bootstrap

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
)

// 版本信息version information
var (
	BuildDate string
	GitCommit string
	Version   string
)

type Config struct {
	// ocr
	BaiduRequestUrl      string `mapstructure:"baidu_request_url"`
	BaiduAccesstoken     string `mapstructure:"baidu_accessToken"`
	JfymRequestUrl       string `mapstructure:"jfym_request_url"`
	JfymToken            string `mapstructure:"jfym_token"`
	YescaptchaRequestUrl string `mapstructure:"yescaptcha_request_url"`
	YesCaptchaToken      string `mapstructure:"yescaptcha_token"`

	// extra
	PgConfig   string `mapstructure:"pg_config"`
	JwtKey     string `mapstructure:"jwt_key"`
	ListenPort string `mapstructure:"listen_port"`
	UserAgent  string `mapstructure:"user_agent"`
}

var C Config

func Loadconfig() {
	parseFlag()

	log.Info("\033[1;34m**********Initing env/config **********\033[0m")
	// Default config
	viper.SetDefault("user_agent", "Mozilla/5.0 (Macintosh Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	viper.SetDefault("jwt_key", "9385g0x98n347tx980y34g9sfgsldkjvilr")

	// 环境变量前缀EDU_
	viper.SetEnvPrefix("edu")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	viper.AddConfigPath(filepath.Dir(pathAbs))

	viper.AddConfigPath("/Users/hello/Projects/eduData-GoBack")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")

	viper.SetConfigFile(viper.GetString("configFile"))

	if err = viper.ReadInConfig(); err != nil {
		log.Errorf("read config file error: %+v\n", err)
	}
	log.Infof("Using config file: %s\n", viper.ConfigFileUsed())

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	if err = viper.Unmarshal(&C); err != nil {
		log.Printf("unmarshal config file error: %+v\n", err)
		return
	}

	// 遍历结构体
	t := reflect.TypeOf(C)
	v := reflect.ValueOf(C)
	for i := 0; i < t.NumField(); i++ {
		log.Infof("%s: %s", t.Field(i).Name, v.Field(i).Interface())
	}

	log.Info("\033[1;34m**********Init flag/env/config success!**********\033[0m")
	return
}

func parseFlag() {
	pflag.StringP("configFile", "c", "", "config file")
	pflag.StringP("listen_port", "l", "8080", "listen address")
	pflag.BoolP("version", "v", false, "version information")
	pflag.BoolP("help", "h", false, "display help information")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}
	if viper.GetBool("version") {
		fmt.Println("version:", Version)
		fmt.Println("buildDate:", BuildDate)
		fmt.Println("gitCommit:", GitCommit)
		os.Exit(0)
	}
}

func InitLog() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
