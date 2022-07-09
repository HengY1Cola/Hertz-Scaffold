package dal

import (
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/repository"
	"sync"
)

type DemoRuleDal struct {
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

func (ins *DemoRuleDal) Find(id int) (*model.Demo, error) {
	temp := &model.Demo{}
	db, err := repository.GetGormDb()
	if err != nil {
		return nil, err
	}
	res := db.Table(model.Demo{}.TableName()).Where("id = ?", id).First(&temp)
	if res.Error != nil {
		return temp, res.Error
	}
	return temp, nil
}
