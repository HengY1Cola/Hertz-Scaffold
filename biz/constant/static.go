package constant

// 用于定义静态数据

const (
	InitTrade = iota
	Trading
	TradeOver
)

const (
	DefaultAPIModule = "defaultAPI"
	DevOpsAPIModule  = "devOpsAPI"
	AdminApIModel    = "adminAPI"

	DefaultURLPrefix = "/api"
	DevOpsURLPrefix  = "/api-dev"
	AdminURLPrefix   = "/api-admin"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodOptions = "OPTIONS"
)
