package common

import (
	con "Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/conf"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

type ErrorResponse struct {
	ErrorCode con.ResponseCodeType `json:"errno"`
	ErrorMsg  string               `json:"err_msg"`
	Data      interface{}          `json:"data"`
	Success   bool                 `json:"success"`
	DebugMsg  string               `json:"debug_msg"`
}

type SuccessResponse struct {
	ErrorCode con.ResponseCodeType `json:"errno"`
	ErrorMsg  string               `json:"err_msg"`
	Data      interface{}          `json:"data"`
	Success   bool                 `json:"success"`
}

func ResponseError(c *app.RequestContext, response con.ErrorResponse, debugErr error) {
	c.Header("tracer_id", GetTracerId(c))
	resp := &ErrorResponse{
		ErrorCode: response.ErrCode,
		ErrorMsg:  response.ErrServerError.CN,
		Data:      nil,
		Success:   response.Success,
	}
	if conf.AppConf.FlagConfig.Type == "dev" {
		resp.DebugMsg = fmt.Sprintf("%v", debugErr)
	}
	c.JSON(response.HttpCode, resp)
	res, _ := json.Marshal(resp)
	c.Set("response", string(res))
	c.AbortWithError(int(response.ErrCode), debugErr)
}

func ResponseSuccess(c *app.RequestContext, data interface{}) {
	c.Header("tracer_id", GetTracerId(c))
	resp := &SuccessResponse{
		ErrorCode: con.ErrSuccess.ErrCode,
		ErrorMsg:  "",
		Data:      data,
		Success:   con.ErrSuccess.Success,
	}
	c.JSON(con.ErrSuccess.HttpCode, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
