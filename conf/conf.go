package conf

import (
	"flag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var AppConf *AppConfig

type AppConfig struct {
	FlagConfig
	YamlConfig
	BaseInfo
}

type FlagConfig struct {
	RunDir string
	Type   string
}

type BaseInfo struct {
	ConfAbsoluteDir string
	LogAbsoluteDir  string
}

type YamlConfig struct {
	Develop      HertzBase
	Product      HertzBase
	MysqlDevelop MysqlBase
	MysqlProduct MysqlBase
	RedisDevelop RedisBase
	RedisProduct RedisBase
}

type HertzBase struct {
	ServicePort int
	Level       string
}

type MysqlBase struct {
	MysqlUrl        string
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifeTime int
}

type RedisBase struct {
	ProxyList []string
	MaxActive int
	MaxIdle   int
	DownGrade bool
}

const (
	HertzConfigFile = "hertz.config"
)

func InitLoadConf() {
	parseFlags()
	parseConf()
}

func parseFlags() {
	AppConf = &AppConfig{}
	confType, confRunDir := "", ""
	flag.StringVar(&confType, "env", "", "pro or dev")
	flag.StringVar(&confRunDir, "run_dir", "", "where to run")
	if !flag.Parsed() {
		flag.Parse()
	}
	AppConf.FlagConfig = FlagConfig{
		RunDir: confRunDir,
		Type:   confType,
	}
	if AppConf.FlagConfig.Type == "" {
		usage()
	}
	if AppConf.FlagConfig.RunDir == "" {
		usage()
	}
	if AppConf.FlagConfig.Type != "dev" {
		AppConf.FlagConfig.Type = "pro"
	}
	AppConf.BaseInfo = BaseInfo{
		ConfAbsoluteDir: filepath.Join(AppConf.FlagConfig.RunDir, "conf"),
		LogAbsoluteDir:  filepath.Join(AppConf.FlagConfig.RunDir, "log"),
	}
}

func usage() {
	flag.Usage()
	os.Exit(-1)
}

func parseConf() {
	v := viper.New()
	confFile := getConfigFile(HertzConfigFile)
	v.SetConfigFile(confFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	var res *YamlConfig
	err := v.Unmarshal(&res)
	if err != nil {
		panic(err)
	}
	AppConf.YamlConfig = *res
}

func getConfigFile(filenames string) string {
	absolutePath := filepath.Join(AppConf.BaseInfo.ConfAbsoluteDir, filenames+".yaml")
	if _, err := os.Stat(absolutePath); err != nil {
		msg := "Failed to load app config"
		panic(msg)
	}
	return absolutePath
}

func (a *AppConfig) GetBaseInfo() HertzBase {
	if a.FlagConfig.Type == "dev" {
		return a.YamlConfig.Develop
	} else {
		return a.YamlConfig.Product
	}
}

func (a *AppConfig) GetMysqlInfo() MysqlBase {
	if a.FlagConfig.Type == "dev" {
		return a.YamlConfig.MysqlDevelop
	} else {
		return a.YamlConfig.MysqlProduct
	}
}

func (a *AppConfig) RedisInfo() RedisBase {
	if a.FlagConfig.Type == "dev" {
		return a.YamlConfig.RedisDevelop
	} else {
		return a.YamlConfig.RedisProduct
	}
}
