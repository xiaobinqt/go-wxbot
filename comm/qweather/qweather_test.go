package qweather

import (
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

func TestGetLocationID(t *testing.T) {
	initAction(t)
	id, err := GetLocationID("北京")
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func TestGetQWeatherDetail(t *testing.T) {
	initAction(t)
	detail, err := GetQWeatherDetail("101010100", "北京")
	if err != nil {
		t.Error(err)
	}
	t.Log(detail)
}
