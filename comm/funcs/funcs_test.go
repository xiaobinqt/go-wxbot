package funcs

import (
	"fmt"
	"testing"
	"time"

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

func TestImg2base64(t *testing.T) {
	ret, err := Img2base64("D:\\go\\src\\go-wxbot\\avatar\\test.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(ret)
}

func TestGetDiffDays(t *testing.T) {

	x := GetDiffDaysSolar(getCurrentDate(), "08-22")
	fmt.Println(x)
}

func TestGetqiXi(t *testing.T) {
	x := getLunar2SolarDate(int64(time.Now().Year()), 7, 10)
	fmt.Println("xxxxx1111", x)
	xx := GetDiffDaysLunar(getCurrentDate(), x, 7, 10)
	fmt.Println(xx)
}

func TestImportDateFormatMsg(t *testing.T) {
	initAction(t)
	x := ImportDateFormatMsg()
	t.Log(x)
}

func TestRemainingDays(t *testing.T) {
	t.Log(RemainingDays())
}
