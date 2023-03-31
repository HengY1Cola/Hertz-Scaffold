package middleware

import (
	con "Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/utils"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/pkg/errors"
)

// JwtAuthMiddleware 示例
func JwtAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logger := utils.GetCtxLogger(c)
		headerAuth := string(c.GetHeader("Authorization"))
		authArray := strings.Split(headerAuth, " ")
		if len(authArray) != 2 {
			logger.Error("[JwtAuthMiddleware] authArray length is not 2")
			utils.ResponseError(c, con.ErrJwtError, errors.New("非法操作"))
			c.Abort()
			return
		}

		var jwtStruct utils.JwtStruct
		jwtStruct, err := jwtStruct.JwtDecode(authArray[1])
		if err != nil {
			logger.Error(fmt.Sprintf("[JwtAuthMiddleware] JwtDecode err %v", err.Error()))
			utils.ResponseError(c, con.ErrJwtError, errors.New("非法操作"))
			c.Abort()
			return
		}

		c.Set("jwtStruct", jwtStruct)
		c.Next(ctx)
	}
}
