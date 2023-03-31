package dal

import (
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CommonDAL struct {
}

const TransactionDbInstance string = "TransactionDbInstance"

func (ins *CommonDAL) GetBaseTransaction(c *app.RequestContext) (*gorm.DB, error) {
	logger := utils.GetCtxLogger(c)
	db, err := repository.GetGormDb()
	if err != nil {
		logger.Error("[GetBaseTransaction] GetGormDb err: %v", err)
		return nil, err
	}
	return db, nil
}

func (ins *CommonDAL) GetTransaction(c *app.RequestContext) (*gorm.DB, error) {
	tx, err := ins.GetBaseTransaction(c)
	if err != nil {
		return nil, err
	}
	return tx.Where("is_delete = ?", model.DeleteFalse), nil
}

func (ins *CommonDAL) GetBaseTransactionWithCtx(c *app.RequestContext) (*app.RequestContext, error) {
	logger := utils.GetCtxLogger(c)
	_, ok := c.Get(TransactionDbInstance)
	if !ok {
		db, err := repository.GetGormDb()
		if err != nil {
			logger.Error("[GetBaseTransactionWithCtx] GetGormDb err: %v", err)
			return c, err
		}
		c.Set(TransactionDbInstance, db)
		return c, nil
	}
	return c, nil
}

func (ins *CommonDAL) GetTransactionWithCtx(c *app.RequestContext) (*gorm.DB, error) {
	value, ok := c.Get(TransactionDbInstance)
	if !ok {
		return nil, fmt.Errorf("no ctx TransactionDbInstance")
	}
	return value.(*gorm.DB), nil
}

func (ins *CommonDAL) BatchDeleteWithCtx(c *app.RequestContext, ids []int, tableName schema.Tabler) error {
	tx, err := ins.GetTransactionWithCtx(c)
	if err != nil {
		return err
	}

	if err := tx.Table(tableName.TableName()).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"is_delete": model.DeleteTrue}).Error; err != nil {
		return err
	}
	return nil
}

func (ins *CommonDAL) BatchDelete(c *app.RequestContext, ids []int, tableName schema.Tabler) error {
	tx, err := ins.GetBaseTransaction(c)
	if err != nil {
		return err
	}

	if err := tx.Table(tableName.TableName()).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"is_delete": model.DeleteTrue}).Error; err != nil {
		return err
	}
	return nil
}
