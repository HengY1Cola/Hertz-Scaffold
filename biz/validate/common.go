package validate

import (
	"fmt"
	"regexp"
)

type ValidatorFun struct {
	Name string
	Func func(args ...interface{}) error
}

type ValidateFuns []ValidatorFun

// GetFuncArray 进行验证方法注册
func GetFuncArray() ValidateFuns {
	return []ValidatorFun{
		ValidatorEmail(),
		ValidatorNum(),
	}
}

func ValidatorEmail() ValidatorFun {
	return ValidatorFun{
		Name: "email",
		Func: func(args ...interface{}) error {
			s, _ := args[0].(string)
			pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
			reg := regexp.MustCompile(pattern)
			if reg.MatchString(s) {
				return nil
			} else {
				return fmt.Errorf("email format error")
			}
		},
	}
}

func ValidatorNum() ValidatorFun {
	return ValidatorFun{
		Name: "mustPositiveNum",
		Func: func(args ...interface{}) error {
			n, _ := args[0].(float64)
			if n <= 0 {
				return fmt.Errorf("num format error")
			} else {
				return nil
			}
		},
	}
}
