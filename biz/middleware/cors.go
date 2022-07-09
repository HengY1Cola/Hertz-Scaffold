package middleware

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
)

func DefaultCorsMiddleware() app.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                           // 默认所有的都可以跨域
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"}, // 先支持这4中方法
		AllowHeaders:     []string{"Origin"},                      // 支持的 non simple headers 字段
		ExposeHeaders:    []string{"Content-Length"},              // 指示向CORS的API公开哪些头是安全的
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 默认给12个小时
	})
}
