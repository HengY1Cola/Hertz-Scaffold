package handler

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/service"
	"Hertz-Scaffold/biz/utils"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"strconv"
	"time"
)

type DemoController struct {
}

func DemoRegister(group *route.RouterGroup) {
	demoController := &DemoController{}
	group.GET("/ping", demoController.Ping)
	group.GET("/get_with_bind/:info", demoController.GetBinding)
	group.GET("/get_with_query", demoController.GetQuery)
	group.POST("/post_with_bind", demoController.PostJsonBinding)
}

// Ping 最基础的Get方法
func (demo *DemoController) Ping(ctx context.Context, c *app.RequestContext) {
	logger := utils.GetCtxLogger(c)
	logger.Info("[Ping] Return Pong Success %d", time.Now().Unix())
	utils.ResponseSuccess(c, bo.BaseMsgAndFlagResponse{
		Msg:  "pong",
		Flag: 200,
	})
}

// GetBinding 使用GET方法在URL上进行绑定
func (demo *DemoController) GetBinding(ctx context.Context, c *app.RequestContext) {
	request := bo.DetailInfo2GetRequest{}
	logger := utils.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.Error("[GetBinding] Bind Err: %v", err.Error())
		utils.ResponseError(c, constant.ErrServerError, err)
		return
	}
	logger.Info("[GetBinding] Success")
	utils.ResponseSuccess(c, request)
}

// GetQuery 使用GET方法在URL上?id=
func (demo *DemoController) GetQuery(ctx context.Context, c *app.RequestContext) {
	request := bo.DemoQueryRequest{}
	logger := utils.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.Error("[GetQuery] Bind Err: %v", err.Error())
		utils.ResponseError(c, constant.ErrServerError, err)
		return
	}

	res, err := service.GetDemoService().GetString(c, strconv.Itoa(request.Id))
	if err != nil {
		logger.Error("[GetQuery] GetString err %v", err)
		utils.ResponseError(c, constant.ErrNoPermission, err)
		return
	}
	utils.ResponseSuccess(c, res)
}

// PostJsonBinding 使用Post进行绑定
func (demo *DemoController) PostJsonBinding(ctx context.Context, c *app.RequestContext) {
	request := bo.DemoRequest{}
	logger := utils.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.Error("[PostJsonBinding] Bind Err: %v", err.Error())
		utils.ResponseError(c, constant.ErrServerError, err)
		return
	}
	logger.Info("[PostJsonBinding] Success")
	utils.ResponseSuccess(c, request)
}
