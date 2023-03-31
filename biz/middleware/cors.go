package middleware

import (
	"Hertz-Scaffold/conf"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
)

func DefaultCorsMiddleware() app.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     getOrigin(),                                         // 默认所有的都可以跨域
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},   // 先支持这4中方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 支持的 non simple headers 字段
		ExposeHeaders:    []string{"Content-Length", "tracer_id"},             // 指示向CORS的API公开哪些头是安全的
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 默认给12个小时
	})
}

func getOrigin() []string {
	if conf.AppConf.FlagConfig.Type == "dev" {
		return []string{"*"}
	} else {
		var res []string
		res = append(res, "http://127.0.0.1")
		res = append(res, "localhost")
		res = append(res, "http://"+conf.AppConf.GetDomainInfo())
		res = append(res, "https://"+conf.AppConf.GetDomainInfo())
		return res
	}
}
