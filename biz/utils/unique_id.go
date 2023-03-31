package utils

import (
	"github.com/google/uuid"
	"strconv"
)

func generateUniqueString(n int) string {
	if n > 36 {
		n = 36
	}
	uniqueString := uuid.New().String()
	uniqueString = uniqueString[:n]
	return uniqueString
}

func generateUniqueInt(n int) int {
	uniqueString := generateUniqueString(n)
	uniqueInt, _ := strconv.Atoi(uniqueString)
	return uniqueInt
}
