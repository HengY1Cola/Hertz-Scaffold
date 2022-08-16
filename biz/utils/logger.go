package utils

import (
	"Hertz-Scaffold/conf"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Logger struct {
	TempLogger *logrus.Logger
	Ctx        *app.RequestContext
}

var GlobalLogger *logrus.Logger

func GetCtxLogger(c *app.RequestContext) Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(GetInfoLevel()) // 设置日志级别
	logger.SetReportCaller(false)   // 设置在输出日志中添加文件名和方法信息 默认关闭
	logfile, _ := os.OpenFile(conf.AppConf.BaseInfo.LogAbsoluteDir+"/"+GetLogFileName(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	logger.Out = logfile
	return Logger{TempLogger: logger, Ctx: c}
}

func (l *Logger) DoInfo(info string) {
	l.TempLogger.WithFields(logrus.Fields{
		"tracer_id": GetTracerId(l.Ctx),
	}).Info(info)
}

func (l *Logger) DoError(err string) {
	l.TempLogger.WithFields(logrus.Fields{
		"tracer_id": GetTracerId(l.Ctx),
	}).Error(err)
}

func (l *Logger) DoDebug(err string) {
	l.TempLogger.WithFields(logrus.Fields{
		"tracer_id": GetTracerId(l.Ctx),
	}).Debug(err)
}

func GetTracerId(c *app.RequestContext) string {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}
	return traceId
}

func GetLogFileName() string {
	timeStr := time.Now().Format("01_02")
	return timeStr + ".log"
}

func GetInfoLevel() logrus.Level {
	if v := conf.AppConf.FlagConfig.Type; v == "dev" {
		return logrus.DebugLevel
	} else {
		return logrus.InfoLevel
	}
}
