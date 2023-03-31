package service

import (
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type DemoRuleService interface {
	GetString(c *app.RequestContext, str string) (string, error)
}

type DemoRuleServiceProxy struct {
	common *CommonService
}

var (
	demoRuleService     DemoRuleService
	demoRuleServiceOnce sync.Once
)

func GetDemoService() DemoRuleService {
	demoRuleServiceOnce.Do(func() {
		demoRuleService = &DemoRuleServiceProxy{
			common: &CommonService{},
		}
	})
	return demoRuleService
}

func (d DemoRuleServiceProxy) GetString(c *app.RequestContext, str string) (string, error) {
	return str, nil
}
