package env

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils/common"
	"Hertz-Scaffold/biz/validate"
	"Hertz-Scaffold/conf"
	"database/sql"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

type Module string

const (
	Conf            Module = "Conf"
	Mysql           Module = "Mysql"
	Logger          Module = "Logger"
	Validate        Module = "Validate"
	UserMapJwtToken Module = "UserMapJwtToken"
)

// 基础部分需要初始化
func getBaseModel() []Module {
	// 需要注意顺序
	return []Module{
		Conf,
		Mysql,
		Logger,
		Validate,
	}
}

func InitModules(modules []Module) {
	for _, base := range getBaseModel() {
		InitModule(base)
	}
	for _, module := range modules {
		InitGlobalModels(module)
	}
}

func InitTestModules(modules []Module) {
	for _, base := range getBaseModel() {
		InitTestModule(base)
	}
	for _, module := range modules {
		InitGlobalModels(module)
	}
}

func InitModule(module Module) {
	switch module {
	case Conf:
		initLoadConf()
	case Mysql:
		initMysqlDb()
	case Logger:
		initGlobalLogger()
	case Validate:
		initValidator()
	}
}

func InitGlobalModels(module Module) {
	switch module {
	case UserMapJwtToken:
		initUserIdMapJtwToken()
	}
}

func InitTestModule(module Module) {
	switch module {
	case Conf:
		initLoadTestConf()
	case Mysql:
		initMysqlDb()
	case Logger:
		initGlobalLogger()
	case Validate:
		initValidator()
	case UserMapJwtToken:
		initUserIdMapJtwToken()
	}
}

func initLoadTestConf() {
	conf.ParseTest()
	conf.ParseConf()
	fmt.Println("########  init test conf over")
}

func initLoadConf() {
	conf.ParseFlags()
	conf.ParseConf()
	fmt.Println("########  init conf over")
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
	fmt.Println("########  init mysql over")
}

func initGlobalLogger() {
	common.GlobalLogger = logrus.New()
	common.GlobalLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	common.GlobalLogger.SetLevel(common.GetInfoLevel()) // 设置日志级别
	common.GlobalLogger.SetReportCaller(false)          // 设置在输出日志中添加文件名和方法信息 默认关闭
	writer, err := common.DivisionWriter(common.GlobalLoggerName)
	if err != nil {
		panic(err)
	}
	common.GlobalLogger.SetOutput(io.MultiWriter(writer))
	fmt.Println("########  init global logger over")
}

func initValidator() {
	binding.SetLooseZeroMode(true)
	for _, value := range validate.GetFuncArray() {
		binding.MustRegValidateFunc(value.Name, value.Func)
	}
	fmt.Println("########  init validate over")
}

func initUserIdMapJtwToken() {
	bo.TempUserIdMapJwtTokenHandle = &bo.UserIdMapJtwToken{
		UserIdMapJtwTokenMap:   map[string]string{},
		UserIdMapJtwTokenSlice: []string{},
		Locker:                 sync.RWMutex{},
	}
	fmt.Println("########  init userIdMap JtwToken over")
}
