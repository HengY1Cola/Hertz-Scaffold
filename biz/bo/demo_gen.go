package bo

type BaseMsgAndFlagResponse struct {
	Msg  string `json:"msg"`
	Flag int    `json:"flag"`
}

type DetailInfo2GetRequest struct {
	Info int `path:"info, required" json:"info"`
}

type DemoRequest struct {
	Type string `json:"type,required" vd:"test($)"`
}

type DemoQueryRequest struct {
	Id int `query:"id, required"`
}
