package dal

import (
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/utils"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type DemoRuleDal struct {
	*CommonDAL
}

var (
	demoRuleDal     *DemoRuleDal
	demoRuleDalOnce sync.Once
)

func GetDemoDal() *DemoRuleDal {
	demoRuleDalOnce.Do(func() {
		demoRuleDal = &DemoRuleDal{}
	})
	return demoRuleDal
}

// 主要定义一些基础的方法,给Service提供，当然也可以直接使用

func (ins *DemoRuleDal) Find(c *app.RequestContext, id int) (*model.Demo, error) {
	logger := utils.GetCtxLogger(c)
	temp := &model.Demo{}
	db, err := ins.GetTransaction(c)
	if err != nil {
		return nil, err
	}
	res := db.Table(model.Demo{}.TableName()).Where("id = ?", id).First(&temp)
	if res.Error != nil {
		logger.DoError(fmt.Sprintf("[DemoRuleDal] Find err: %v", res.Error))
		return temp, res.Error
	}
	return temp, nil
}
