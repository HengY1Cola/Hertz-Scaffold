package utils

import (
	"Hertz-Scaffold/biz/constant"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt"
)

type JwtStruct struct {
	UserId   int
	NickName string
	DueTime  int
}

func JwtEncode(jwtDemo JwtStruct) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   jwtDemo.UserId,
		"nick_name": jwtDemo.NickName,
		"due_time":  jwtDemo.DueTime,
	})

	tokenString, err := token.SignedString([]byte(constant.BackEndJwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JwtDecode(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{
		"user_id":   "user_id",
		"nick_name": "nick_name",
		"due_time":  "due_time",
	}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.BackEndJwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

func HandleCtxJwtStruct(c *app.RequestContext) JwtStruct {
	tempStruct, _ := c.Get("jwtStruct")
	jwtStruct := tempStruct.(JwtStruct)
	return jwtStruct
}
