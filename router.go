package main

import (
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/utils/common"
	"Hertz-Scaffold/conf"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"strconv"
)

func InitRouter(middlewares ...app.HandlerFunc) *server.Hertz {
	engine := server.New(server.WithHostPorts("[::]:" + strconv.Itoa(conf.AppConf.GetBaseInfo().ServicePort)))
	engine.Use(middlewares...) // 通用中间件

	defaultRouter := engine.Engine
	defaultRouter.Use()
	adminRouter := engine.Engine
	adminRouter.Use()
	devopsRouter := engine.Engine
	devopsRouter.Use()

	for _, v := range common.GlobalApiList {
		switch v.Model {
		case constant.DefaultAPIModule:
			common.EngineRegister(defaultRouter, v)
		case constant.DevOpsAPIModule:
			common.EngineRegister(devopsRouter, v)
		case constant.AdminApIModel:
			common.EngineRegister(adminRouter, v)
		}
	}
	return engine
}
