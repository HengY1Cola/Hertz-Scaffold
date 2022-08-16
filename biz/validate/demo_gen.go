package validate

import (
	"fmt"
)

type ValidatorFun struct {
	Name string
	Func func(args ...interface{}) error
}

type ValidateFuncs []ValidatorFun

func ValidatorTest() ValidatorFun {
	return ValidatorFun{
		Name: "test",
		Func: func(args ...interface{}) error {
			s, _ := args[0].(string)
			if s == "123" {
				return fmt.Errorf("the args can not be 123")
			}
			return nil
		},
	}
}
