package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/sirupsen/logrus"
	conf2 "go-wxbot/openwechat/comm/conf"
	"go-wxbot/openwechat/comm/global"
)

func initAction(t *testing.T) {
	extra.RegisterFuzzyDecoders()
	logrus.SetLevel(logrus.DebugLevel)
	var (
		err error
	)
	conf, err := conf2.GetConf("../../config/prod.yaml")
	if err != nil {
		t.Logf("get conf err:%s ", err.Error())
		return
	}

	global.Conf = conf
}

func TestGetRedis(t *testing.T) {
	initAction(t)

	client := GetRedis()
	err := client.Set(context.Background(), "test", 1111, 0).Err()
	fmt.Println("err ===", err)
}
