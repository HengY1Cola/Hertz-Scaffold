package middleware

import (
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/utils"
	"Hertz-Scaffold/conf"
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/cloudwego/hertz/pkg/app"
)

func RecoveryMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				fmt.Println(string(debug.Stack()))
				logger := utils.GetCtxLogger(c)
				logger.DoError(fmt.Sprintf("error: %v; stack: %v", fmt.Sprint(err), string(debug.Stack())))
				if conf.AppConf.FlagConfig.Type != "dev" {
					utils.ResponseError(c, constant.SysErrCode, errors.New("内部错误"))
					return
				} else {
					utils.ResponseError(c, constant.SysErrCode, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next(ctx)
	}
}
