package constant

import "net/http"

type ResponseCodeType uint64

type Locale struct {
	EN string
	CN string
}

type ErrorResponse struct {
	HttpCode       int
	ErrCode        ResponseCodeType
	Success        bool
	ErrServerError Locale
}

var (
	ErrSuccess = ErrorResponse{
		HttpCode:       http.StatusOK,
		ErrCode:        0,
		Success:        true,
		ErrServerError: Locale{},
	}

	ErrServerError = ErrorResponse{
		HttpCode: http.StatusBadRequest,
		ErrCode:  100001,
		Success:  false,
		ErrServerError: Locale{
			EN: "System error",
			CN: "系统异常",
		},
	}

	ErrNoPermission = ErrorResponse{
		HttpCode: http.StatusBadRequest,
		ErrCode:  100002,
		Success:  false,
		ErrServerError: Locale{
			EN: "Permission error",
			CN: "权限错误",
		},
	}

	ErrJwtError = ErrorResponse{
		HttpCode: http.StatusBadRequest,
		ErrCode:  100003,
		Success:  false,
		ErrServerError: Locale{
			EN: "JWT parsing error",
			CN: "JWT解析错误",
		},
	}
)
