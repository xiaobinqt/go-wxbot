package tian

import (
	"fmt"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/sirupsen/logrus"
	conf2 "go-wxbot/openwechat/comm/conf"
	"go-wxbot/openwechat/comm/global"
)

func initAction(t *testing.T) (conf *conf2.Conf) {
	extra.RegisterFuzzyDecoders()
	logrus.SetLevel(logrus.DebugLevel)
	var (
		err error
	)
	conf, err = conf2.GetConf("../../config/prod.yaml")
	if err != nil {
		t.Logf("get conf err:%s ", err.Error())
		return
	}

	global.Conf = conf

	return conf
}

func TestGetMessage(t *testing.T) {
	_ = initAction(t)
	ret, err := GetMessage(C_godreply)
	fmt.Println(ret, err)
	ret, err = GetMessage(C_mingyan)
	fmt.Println(ret, err)
	ret, err = GetMessage(C_caipu, "红烧肉")
	fmt.Println(ret, err)
	ret, err = GetMessage(C_caipu, "西红柿鸡蛋汤")
	fmt.Println(ret, err)
	ret, err = GetMessage(C_englishSentence)
	fmt.Println(ret, err)
}

func TestGetMessageV1(t *testing.T) {
	_ = initAction(t)
	ret, err := GetMessageV1(C_lizhiguyan)
	fmt.Println(ret, err)
}
