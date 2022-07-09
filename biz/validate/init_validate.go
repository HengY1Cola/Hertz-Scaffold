package validate

import "github.com/cloudwego/hertz/pkg/app/server/binding"

type ValidatorFun struct {
	Name string
	Func func(args ...interface{}) error
}

type GlobalValidateFunc []ValidatorFun

func InitValidator() {
	binding.SetLooseZeroMode(true)
	funcArray := GlobalValidateFunc{ // 进行函数的注册
		ValidatorTest(),
	}
	for _, value := range funcArray {
		binding.MustRegValidateFunc(value.Name, value.Func)
	}
}
