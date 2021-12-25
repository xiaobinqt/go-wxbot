package funcs

import (
	"fmt"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/sirupsen/logrus"
	conf2 "go-wxbot/comm/conf"
	"go-wxbot/comm/global"
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

func TestImg2base64(t *testing.T) {
	ret, err := Img2base64("D:\\go\\src\\go-wxbot\\avatar\\test.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(ret)
}
