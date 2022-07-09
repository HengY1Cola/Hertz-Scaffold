package handler

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/utils"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type DemoController struct {
}

func DemoRegister(group *route.RouterGroup) {
	demoController := &DemoController{}
	group.GET("/ping", demoController.Ping)
	group.GET("/get_with_bind/:info", demoController.GetBinding)
	group.POST("/post_with_bind", demoController.PostJsonBinding)
}

// Ping 最基础的Get方法
func (demo *DemoController) Ping(ctx context.Context, c *app.RequestContext) {
	logger := utils.GetCtxLogger(c)
	logger.DoInfo("[Ping] Return Pong")
	utils.ResponseSuccess(c, bo.BaseMsgAndFlagResponse{
		Msg:  "pong",
		Flag: constant.SuccessCode,
	})
}

// GetBinding 使用GET方法在URL上进行绑定
func (demo *DemoController) GetBinding(ctx context.Context, c *app.RequestContext) {
	request := bo.DetailInfo2GetRequest{}
	logger := utils.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.DoError(fmt.Sprintf("[GetBinding] Bind Err: %v", err.Error()))
		utils.ResponseError(c, constant.BindErrCode, err)
		return
	}
	logger.DoInfo("[GetBinding] Success")
	utils.ResponseSuccess(c, request)
}

// PostJsonBinding 使用Post进行绑定
func (demo *DemoController) PostJsonBinding(ctx context.Context, c *app.RequestContext) {
	request := bo.DemoRequest{}
	logger := utils.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.DoError(fmt.Sprintf("[PostJsonBinding] Bind Err: %v", err.Error()))
		utils.ResponseError(c, constant.BindErrCode, err)
		return
	}
	logger.DoInfo("[PostJsonBinding] Success")
	utils.ResponseSuccess(c, request)
}
