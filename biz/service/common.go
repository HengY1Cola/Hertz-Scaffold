package service

import (
	"Hertz-Scaffold/biz/dal"
	"github.com/cloudwego/hertz/pkg/app"
)

type CommonService struct {
}

func (ins *CommonService) ExecTransaction(c *app.RequestContext, fn func(c *app.RequestContext) error) error {
	// 给上下文注入tx
	c, err := (&dal.CommonDAL{}).GetBaseTransactionWithCtx(c)
	if err != nil {
		return err
	}
	// 拿到上下文中的tx
	tx, err := (&dal.CommonDAL{}).GetTransactionWithCtx(c)
	if err != nil {
		return err
	}
	// 执行对应的事务
	tx = tx.Begin()
	err = fn(c)
	if err != nil {
		tx.Callback()
		return err
	}
	tx.Commit()
	return nil
}
