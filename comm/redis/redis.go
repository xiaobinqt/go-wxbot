package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
)

var (
	gClient *redis.Client
)

// RedisInit .
func RedisInit() {
	gClient = redis.NewClient(&redis.Options{
		Addr:     getAddr(),
		Password: global.Conf.RedisConf.Passwd,
	})

	return
}

func getAddr() string {
	return fmt.Sprintf("%s:%s", global.Conf.RedisConf.IP, global.Conf.RedisConf.Port)
}

// GetRedis 单例连接池
func GetRedis() (c *redis.Client) {
	if gClient != nil {
		return gClient
	}

	return redis.NewClient(&redis.Options{
		Addr:     getAddr(),
		Password: global.Conf.RedisConf.Passwd,
	})
}

func CloseRedis() {
	err := gClient.Close()
	if err != nil {
		logrus.Errorf("CloseRedis redis err:%s ", err.Error())
	}
	return
}
