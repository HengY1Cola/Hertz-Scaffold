package main

import (
	"Hertz-Scaffold/biz/middleware"
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils"
	"Hertz-Scaffold/cron_job"
	"time"
)

func main() {
	start := time.Now()
	modules := []Module{
		Conf,
		Mysql,
		Logger,
		Validate,
	}
	InitModules(modules)
	utils.GlobalLogger.Infof("########time: %s ", time.Since(start))
	engine := InitRouter(
		middleware.RequestDoTracerId(),     // 全局链路中间件
		middleware.RecoveryMiddleware(),    // 最后捕获panic错误
		middleware.DefaultCorsMiddleware(), // 默认跨域
	)
	go cron_job.InitCronJob() // 开启定时任务
	defer repository.SqlDbPool.Close()
	engine.Spin() // 开启Http服务
}
