package constant

type ResponseCode int

// BusinessCode 定义业务返回码
type BusinessCode ResponseCode

const (
	RefuseCode  BusinessCode = 0 // 业务逻辑错误
	SuccessCode BusinessCode = 1 // 业务操作成功
)

// SystemCode 定义系统错误
type SystemCode ResponseCode

const (
	SysErrCode        ResponseCode = 500  // panic错误
	MiddleWareErrCode ResponseCode = 998  // 中间件错误
	DataBaseErrCode   ResponseCode = 999  // 数据库发生错误
	BindErrCode       ResponseCode = 1000 // 绑定对应Request错误
)

// LogicCode 定义逻辑上错误
type LogicCode ResponseCode

const (
	NoFindUserErrCode   ResponseCode = 1001 // 没有找到用户错误
	NoCorrectPwdErrCode ResponseCode = 1002 // 密码错误
)
