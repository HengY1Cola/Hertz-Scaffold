package bo

import (
	c "Hertz-Scaffold/biz/constant"
)

type BaseMsgAndFlagResponse struct {
	Msg  string         `json:"msg"`
	Flag c.BusinessCode `json:"flag"`
}

type DetailInfo2GetRequest struct {
	Info int `path:"info, required" json:"info"` // 对应的申请ID
}

type DemoRequest struct {
	Type string `json:"type,required" vd:"test($)"` // 申请类型
}
