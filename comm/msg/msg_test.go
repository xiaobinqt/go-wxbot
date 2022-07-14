package msg

import (
	"fmt"
	"os"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/sirupsen/logrus"
	conf2 "go-wxbot/openwechat/comm/conf"
	"go-wxbot/openwechat/comm/funcs"
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

func Test_1(t *testing.T) {
	initAction(t)
	os.Chdir("../../")
	base64, err := funcs.Img2base64("D:\\go\\src\\go-wxbot\\avatar\\test.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hatBase64, err := AvatarAddChristmasHat(base64)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	filename, err := SaveImageToDisk("ddd", hatBase64)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("success ===================", filename)
}
