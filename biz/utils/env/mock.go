package env

import (
	"Hertz-Scaffold/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetMockCtx() *app.RequestContext {
	ctx := app.NewContext(10)
	value := utils.JwtStruct{
		UserId:    1,
		NickName:  "Test用户",
		AvatarUrl: "https://www.baidu.com",
		DueTime:   -1,
	}
	ctx.Set("jwtStruct", value)
	return ctx
}

func GetMockCtxWithoutUser() *app.RequestContext {
	ctx := app.NewContext(10)
	return ctx
}

func GetMockCtxUserId(c *app.RequestContext) int {
	res, _ := c.Get("jwtStruct")
	info := res.(utils.JwtStruct)
	return info.UserId
}

func GetMockCtxWithAimId(userId int) *app.RequestContext {
	ctx := app.NewContext(10)
	value := utils.JwtStruct{
		UserId:    userId,
		NickName:  "Test用户",
		AvatarUrl: "https://www.baidu.com",
		DueTime:   -1,
	}
	ctx.Set("jwtStruct", value)
	return ctx
}
