package weather

import (
	"fmt"
	"os"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/sirupsen/logrus"
	conf2 "go-wxbot/openwechat/comm/conf"
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

	return conf
}

func TestGetWeatherInfo(t *testing.T) {
	conf := initAction(t)
	os.Chdir("../../")
	format, err := GetFormatWeatherMessage(conf, "泾县")
	fmt.Println(format, err)
}
