package repository

import (
	"Hertz-Scaffold/conf"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

// RedisConfPip 流水线模式
func RedisConfPip(pip ...func(c redis.Conn)) error {
	c, err := RedisAloneConnWithoutPasswd()
	if err != nil {
		return err
	}
	defer c.Close()
	for _, f := range pip {
		f(c)
	}
	c.Flush()
	return nil
}

// RedisConfDo 单独操作
func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	c, err := RedisAloneConnWithoutPasswd()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.Do(commandName, args...)
}

// TODO 后面在丰富吧，现在仅仅单机无密码模式

type RedisConf struct {
	ProxyList    []string
	Password     string
	Db           int
	ConnTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

func RedisAloneConnWithoutPasswd() (redis.Conn, error) {
	if len(conf.AppConf.RedisInfo().ProxyList) >= 1 {
		aloneHost := conf.AppConf.RedisInfo().ProxyList[0]
		c, err := redis.Dial(
			"tcp",
			aloneHost,
			redis.DialConnectTimeout(time.Duration(50)*time.Millisecond),
			redis.DialReadTimeout(time.Duration(100)*time.Millisecond),
			redis.DialWriteTimeout(time.Duration(100)*time.Millisecond))
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, errors.New("create redis conn fail")
}
