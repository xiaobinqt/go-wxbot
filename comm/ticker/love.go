package ticker

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/tian"
)

func SendMessageToLover(prefix, stype string) {
	var (
		err     error
		message string
	)
	message, err = tian.GetMessage(stype)
	if err != nil {
		err = errors.Wrapf(err, "SendMessageToLover get message err")
		logrus.Error(err.Error())
		return
	}

	message = fmt.Sprintf("%s%s", prefix, message)
	err = global.WxFriends.SearchByRemarkName(1, global.Conf.Keys.HoneyLove).SendText(message)
	if err != nil {
		err = errors.Wrapf(err, "SendMessageToLover err")
		logrus.Error(err.Error())
	}
}

func LoveTicker() {
	for {
		select {
		case t := <-time.After(1 * time.Minute):
			nowTime := t.Format("15:04")
			if nowTime == "09:30" {
				SendMessageToLover("亲爱的，早上好！爱你每一天！\n新的一天从一句土味情话开始：", tian.C_saylove)
			}
			if nowTime == "23:00" {
				SendMessageToLover("亲爱的，11 点了，该洗漱睡觉了！\n临睡之际送你一句土味情话：", tian.C_saylove)
			}
		}
	}
}
