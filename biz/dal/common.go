package dal

import (
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

type CommonDAL struct {
}

func (ins *CommonDAL) GetTransaction(c *app.RequestContext) (*gorm.DB, error) {
	logger := utils.GetCtxLogger(c)
	db, err := repository.GetGormDb()
	if err != nil {
		logger.DoError(fmt.Sprintf("[FindById] GetGormDb err: %v", err))
		return nil, err
	}
	return db, nil
}
