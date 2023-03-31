package common

import (
	"Hertz-Scaffold/conf"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"time"
)

type Logger struct {
	TempLogger *logrus.Logger
	Ctx        *app.RequestContext
}

var GlobalLogger *logrus.Logger

const (
	CtxLoggerName    = "CtxLogger"
	GlobalLoggerName = "GlobalLogger"
)

var LoggerList = []string{CtxLoggerName, GlobalLoggerName}

func GetCtxLogger(c *app.RequestContext) Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(GetInfoLevel()) // 设置日志级别
	logger.SetReportCaller(false)   // 设置在输出日志中添加文件名和方法信息 默认关闭
	writer, _ := DivisionWriter(CtxLoggerName)
	logger.SetOutput(io.MultiWriter(writer))
	return Logger{TempLogger: logger, Ctx: c}
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.TempLogger.WithFields(logrus.Fields{
		"tracer_id": GetTracerId(l.Ctx),
	}).Info(fmt.Sprintf(format, a))
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.TempLogger.WithFields(logrus.Fields{
		"tracer_id": GetTracerId(l.Ctx),
	}).Error(fmt.Errorf(format, a))
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

func GetLogFileName(name string) string {
	return filepath.Join(conf.AppConf.LogAbsoluteDir, name+".log")
}

func GetInfoLevel() logrus.Level {
	if v := conf.AppConf.FlagConfig.Type; v == "dev" {
		return logrus.DebugLevel
	} else {
		return logrus.InfoLevel
	}
}

func DivisionWriter(name string) (*rotatelogs.RotateLogs, error) {
	writer, err := rotatelogs.New(
		GetLogFileName(name)+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(GetLogFileName(name)),
		rotatelogs.WithMaxAge(time.Duration(72)*time.Hour),       // 保留最近3天的日志文件
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour), // 每隔1天轮转一个新文件
	)
	return writer, err
}
