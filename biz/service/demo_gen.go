package service

import (
	"sync"
)

type DemoRuleService struct {
}

var (
	demoRuleService     *DemoRuleService
	demoRuleServiceOnce sync.Once
)

func GetDemoService() *DemoRuleService {
	demoRuleServiceOnce.Do(func() {
		demoRuleService = &DemoRuleService{}
	})
	return demoRuleService
}

// 主要在dal上面进行复杂的业务组合
