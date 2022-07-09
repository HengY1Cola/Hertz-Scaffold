package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GetTimeTick64() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetFormatTime(time time.Time) string {
	return time.Format("0102")
}

// GenerateCode 用于单机房生成唯一Id
func GenerateCode() int {
	date := GetFormatTime(time.Now())
	r := rand.Intn(1000)
	code := fmt.Sprintf("%s%d%03d", date, GetTimeTick64(), r)
	randomId, _ := strconv.Atoi(code)
	return randomId
}
