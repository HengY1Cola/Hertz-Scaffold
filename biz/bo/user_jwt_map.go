package bo

import (
	"errors"
	"sync"
)

var TempUserIdMapJwtTokenHandle *UserIdMapJtwToken

type UserIdMapJtwToken struct {
	UserIdMapJtwTokenMap   map[string]string
	UserIdMapJtwTokenSlice []string
	Locker                 sync.RWMutex
}

func (handle *UserIdMapJtwToken) InHandler(tempUserId, JwtToken string) error {
	for _, item := range handle.UserIdMapJtwTokenSlice {
		if item == tempUserId {
			return errors.New("写入失败")
		}
	}

	handle.UserIdMapJtwTokenSlice = append(handle.UserIdMapJtwTokenSlice, tempUserId)
	handle.Locker.Lock()
	defer handle.Locker.Unlock()
	handle.UserIdMapJtwTokenMap[tempUserId] = JwtToken
	return nil
}

func (handle *UserIdMapJtwToken) OutHandler(tempUserId string) string {
	for index, item := range handle.UserIdMapJtwTokenSlice {
		if item == tempUserId {
			handle.Locker.Lock()
			defer handle.Locker.Unlock() // 这个场景不会存在for defer
			temp := handle.UserIdMapJtwTokenMap[tempUserId]
			delete(handle.UserIdMapJtwTokenMap, tempUserId) // 删除map的
			handle.UserIdMapJtwTokenSlice = append(handle.UserIdMapJtwTokenSlice[:index], handle.UserIdMapJtwTokenSlice[index+1:]...)
			return temp
		}
	}
	return ""
}
