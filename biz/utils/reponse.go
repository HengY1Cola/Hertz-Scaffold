package utils

import (
	con "Hertz-Scaffold/biz/constant"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
)

type Response struct {
	ErrorCode con.ResponseCode `json:"errno"`
	ErrorMsg  string           `json:"errmsg"`
	Data      interface{}      `json:"data"`
}

func ResponseError(c *app.RequestContext, code con.ResponseCode, err error) {
	c.Header("tracer_id", GetTracerId(c))
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	resp := &Response{ErrorCode: code, ErrorMsg: errStr, Data: ""}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(200, err)
}

func ResponseSuccess(c *app.RequestContext, data interface{}) {
	c.Header("tracer_id", GetTracerId(c))
	resp := &Response{ErrorCode: con.ResponseCode(con.SuccessCode), ErrorMsg: "", Data: data}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
