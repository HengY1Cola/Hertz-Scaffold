package common

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

func WithAdminCheck(fn func(ctx context.Context, c *app.RequestContext)) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("该操作方法前触发 可作为检测等")
		fn(ctx, c)
	}
}
