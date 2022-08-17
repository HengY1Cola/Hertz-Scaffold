package main

import (
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils"
	"Hertz-Scaffold/biz/validate"
	"Hertz-Scaffold/conf"
	"database/sql"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type Module string

const (
	Conf     Module = "Conf"     // 配置文件
	Mysql    Module = "Mysql"    // Mysql数据库
	Logger   Module = "Logger"   // 全局日志
	Validate Module = "Validate" // 验证器
)

func InitModules(modules []Module) {
	for _, module := range modules {
		InitModule(module)
	}
}

func InitModule(module Module) {
	switch module {
	case Conf:
		initLoadConf()
	case Mysql:
		initMysqlDb()
	case Logger:
		//initBaseLog() // 手动关闭 会自己创建
		initGlobalLogger()
	case Validate:
		initValidator()
	}
}

func initLoadConf() {
	conf.ParseFlags()
	conf.ParseConf()
}

func initMysqlDb() {
	sqlDB, err := sql.Open("mysql", conf.AppConf.GetMysqlInfo().MysqlUrl)
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(conf.AppConf.GetMysqlInfo().MaxIdleConn)                       // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(conf.AppConf.GetMysqlInfo().MaxOpenConn)                       // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Duration(conf.AppConf.GetMysqlInfo().MaxConnLifeTime)) //设置了连接可复用的最大时间
	repository.SqlDbPool = sqlDB
}

func initBaseLog() {
	for _, v := range utils.LoggerList {
		_, err := os.Open(utils.GetLogFileName(v))
		if err != nil && os.IsNotExist(err) {
			file, _ := os.Create(utils.GetLogFileName(v))
			defer file.Close()
		}
	}
}

func initGlobalLogger() {
	utils.GlobalLogger = logrus.New()
	utils.GlobalLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	utils.GlobalLogger.SetLevel(utils.GetInfoLevel()) // 设置日志级别
	utils.GlobalLogger.SetReportCaller(false)         // 设置在输出日志中添加文件名和方法信息 默认关闭
	writer, err := utils.DivisionWriter(utils.GlobalLoggerName)
	if err != nil {
		panic(err)
	}
	utils.GlobalLogger.SetOutput(io.MultiWriter(writer))
}

func initValidator() {
	binding.SetLooseZeroMode(true)
	funcArray := validate.ValidateFuncs{ // 进行函数的注册
		validate.ValidatorTest(),
	}
	for _, value := range funcArray {
		binding.MustRegValidateFunc(value.Name, value.Func)
	}
}
