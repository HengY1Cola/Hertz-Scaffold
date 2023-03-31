package handler

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/service"
	"Hertz-Scaffold/biz/utils/common"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"time"
)

func init() {
	common.Register(constant.DefaultAPIModule, constant.MethodGet, "/ping", Ping)
	common.Register(constant.DefaultAPIModule, constant.MethodGet, "/get_with_bind/:info", GetBinding)
	common.Register(constant.DefaultAPIModule, constant.MethodGet, "/get_with_query", GetQuery)
	common.Register(constant.DefaultAPIModule, constant.MethodPost, "/post_with_bind", PostJsonBinding)
}

// Ping 最基础的Get方法
func Ping(ctx context.Context, c *app.RequestContext) {
	logger := common.GetCtxLogger(c)
	logger.Info("[Ping] Return Pong Success %d", time.Now().Unix())
	common.ResponseSuccess(c, bo.BaseMsgAndFlagResponse{
		Msg:  "pong",
		Flag: 200,
	})
}

// GetBinding 使用GET方法在URL上进行绑定
func GetBinding(ctx context.Context, c *app.RequestContext) {
	request := bo.DetailInfo2GetRequest{}
	logger := common.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.Error("[GetBinding] Bind Err: %v", err.Error())
		common.ResponseError(c, constant.ErrServerError, err)
		return
	}
	logger.Info("[GetBinding] Success")
	common.ResponseSuccess(c, request)
}

// GetQuery 使用GET方法在URL上?id=
func GetQuery(ctx context.Context, c *app.RequestContext) {
	request := bo.DemoQueryRequest{}
	logger := common.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.Error("[GetQuery] Bind Err: %v", err.Error())
		common.ResponseError(c, constant.ErrServerError, err)
		return
	}

	res, err := service.GetDemoService().GetString(c, strconv.Itoa(request.Id))
	if err != nil {
		logger.Error("[GetQuery] GetString err %v", err)
		common.ResponseError(c, constant.ErrNoPermission, err)
		return
	}
	common.ResponseSuccess(c, res)
}

// PostJsonBinding 使用Post进行绑定
func PostJsonBinding(ctx context.Context, c *app.RequestContext) {
	request := bo.DemoRequest{}
	logger := common.GetCtxLogger(c)
	err := c.BindAndValidate(&request)
	if err != nil {
		logger.Error("[PostJsonBinding] Bind Err: %v", err.Error())
		common.ResponseError(c, constant.ErrServerError, err)
		return
	}
	logger.Info("[PostJsonBinding] Success")
	common.ResponseSuccess(c, request)
}
