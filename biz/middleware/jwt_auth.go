package middleware

import (
	con "Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/utils"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// JwtAuthMiddleware 示例
func JwtAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logger := utils.GetCtxLogger(c)
		headerAuth := string(c.GetHeader("Authorization"))
		authArray := strings.Split(headerAuth, " ")
		if len(authArray) != 2 {
			logger.DoError("[JwtAuthMiddleware] authArray length is not 2")
			utils.ResponseError(c, con.MiddleWareErrCode, errors.New("非法操作"))
			c.Abort()
			return
		}

		jwtToken, jwtStruct := authArray[1], utils.JwtStruct{}
		DecodeMap, err := utils.JwtDecode(jwtToken)
		if err != nil {
			logger.DoError(fmt.Sprintf("[JwtAuthMiddleware] JwtDecode err %v", err.Error()))
			utils.ResponseError(c, con.MiddleWareErrCode, errors.New("非法操作"))
			c.Abort()
			return
		}
		if dueTime, ok := DecodeMap["due_time"].(float64); !ok || dueTime <= float64(time.Now().Unix()) {
			logger.DoError("[JwtAuthMiddleware] due_time err")
			utils.ResponseError(c, con.MiddleWareErrCode, errors.New("非法操作"))
			c.Abort()
			return
		}
		if userId, ok := DecodeMap["user_id"].(float64); !ok || userId == 0 {
			logger.DoError("[JwtAuthMiddleware] user_id err")
			utils.ResponseError(c, con.MiddleWareErrCode, errors.New("非法操作"))
			c.Abort()
			return
		} else {
			jwtStruct.UserId = int(userId)
		}
		if nickName, ok := DecodeMap["nick_name"].(string); !ok || nickName == "" {
			logger.DoError("[JwtAuthMiddleware] nick_name err")
			utils.ResponseError(c, con.MiddleWareErrCode, errors.New("非法操作"))
			c.Abort()
			return
		} else {
			jwtStruct.NickName = nickName
		}

		c.Set("jwtStruct", jwtStruct)
		c.Next(ctx)
	}
}
