package utils

import (
	"Hertz-Scaffold/biz/constant"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt"
)

type JwtStruct struct {
	UserId    int
	NickName  string
	AvatarUrl string
	DueTime   int
}

func (j *JwtStruct) JwtEncode() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    j.UserId,
		"nick_name":  j.NickName,
		"avatar_url": j.AvatarUrl,
		"due_time":   time.Now().Unix() + 3600*2, // 2个小时时间有效
	})

	tokenString, err := token.SignedString([]byte(constant.BackEndJwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JwtStruct) JwtDecode(tokenString string) (JwtStruct, error) {
	resStruct := JwtStruct{}
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{
		"user_id":    "user_id",
		"nick_name":  "nick_name",
		"avatar_url": "avatar_url",
		"due_time":   "due_time",
	}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.BackEndJwtKey), nil
	})
	if err != nil {
		return resStruct, err
	}

	tokenMap := token.Claims.(jwt.MapClaims)
	if dueTime, ok := tokenMap["due_time"].(float64); !ok || dueTime <= float64(time.Now().Unix()) {
		return resStruct, errors.New("due_time err")
	} else {
		resStruct.DueTime = int(dueTime)
	}
	if userId, ok := tokenMap["user_id"].(float64); !ok || userId == 0 {
		return resStruct, errors.New("user_id err")
	} else {
		resStruct.UserId = int(userId)
	}
	if nickName, ok := tokenMap["nick_name"].(string); !ok || nickName == "" {
		return resStruct, errors.New("nick_name err")
	} else {
		resStruct.NickName = nickName
	}
	if avatar, ok := tokenMap["avatar_url"].(string); !ok || avatar == "" {
		return resStruct, errors.New("avatar err")
	} else {
		resStruct.AvatarUrl = avatar
	}

	return resStruct, nil
}

func HandleCtxJwtStruct(c *app.RequestContext) JwtStruct {
	tempStruct, _ := c.Get("jwtStruct")
	jwtStruct := tempStruct.(JwtStruct)
	return jwtStruct
}
