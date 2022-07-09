package main

import (
	"Hertz-Scaffold/biz/handler"
	"Hertz-Scaffold/conf"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"strconv"
)

func InitRouter(middlewares ...app.HandlerFunc) *server.Hertz {
	engine := server.New(server.WithHostPorts("[::]:" + strconv.Itoa(conf.AppConf.GetBaseInfo().ServicePort)))
	engine.Use(middlewares...) // 通用中间件

	// 进行分组路由
	demoRouter := engine.Group("/api/demo")
	demoRouter.Use() // 配置该路由的中间件
	{
		handler.DemoRegister(demoRouter)
	}

	return engine
}
