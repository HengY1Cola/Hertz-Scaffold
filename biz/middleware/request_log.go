package middleware

import (
	"Hertz-Scaffold/biz/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
)

func RequestInLog(c *app.RequestContext) {
	traceContext := utils.NewTrace()
	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)
}

func RequestDoTracerId() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		RequestInLog(c)
		info := getRequestInfo(c)
		utils.GlobalLogger.WithFields(logrus.Fields{
			"tracer_id": utils.GetTracerId(c),
		}).Info(info)
		c.Next(ctx)
	}
}

func getRequestInfo(c *app.RequestContext) string {
	requestMap := map[string]interface{}{
		"uri":        string(c.Request.RequestURI()),
		"method":     string(c.Request.Method()),
		"args":       string(c.Request.PostArgString()),
		"body":       string(c.Request.Body()),
		"from":       c.ClientIP(),
		"user-agent": string(c.UserAgent()),
		"status":     c.GetResponse().StatusCode(),
	}
	marshal, err := json.Marshal(requestMap)
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}
