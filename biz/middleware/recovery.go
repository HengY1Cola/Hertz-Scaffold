package middleware

import (
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/utils/common"
	"context"
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
				logger := common.GetCtxLogger(c)
				logger.Error("error: %v; stack: %v", fmt.Sprint(err), string(debug.Stack()))
				common.ResponseError(c, constant.ErrServerError, fmt.Errorf("%v", err))
				return
			}
		}()
		c.Next(ctx)
	}
}
