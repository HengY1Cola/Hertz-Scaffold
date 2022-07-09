package main

import (
	"Hertz-Scaffold/biz/middleware"
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils"
	"Hertz-Scaffold/biz/validate"
	"Hertz-Scaffold/conf"
	"Hertz-Scaffold/cron_job"
)

func init() {
	conf.InitLoadConf()      // 加载配置文件
	repository.InitMysqlDb() // Mysql链接池子
	utils.InitGlobalLogger() // 全局Logger
	validate.InitValidator() // 验证器自定义函数注册
}

func main() {
	engine := InitRouter(
		middleware.RequestDoTracerId(),     // 全局链路中间件
		middleware.RecoveryMiddleware(),    // 最后捕获panic错误
		middleware.DefaultCorsMiddleware(), // 默认跨域
	)
	go cron_job.InitCronJob() // 开启定时任务
	defer repository.SqlDbPool.Close()
	engine.Spin() // 开启Http服务
}
