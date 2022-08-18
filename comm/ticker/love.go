package ticker

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/funcs"
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
			if nowTime == "10:00" {
				lz, err := tian.GetMessageV1(tian.C_lizhiguyan)
				message := ""
				if err != nil {
					message = fmt.Sprintf("盛年不重来，一日难再晨。及时当勉励，岁月不待人。\n今年还剩 %d 天。", funcs.RemainingDays())
				} else {
					message = fmt.Sprintf("%s\n今年还剩 %d 天。", lz, funcs.RemainingDays())
				}

				err = global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.HoneyLove).
					SendText(message)
				if err != nil {
					err = errors.Wrapf(err, "SendMessageToHoneyLove err")
					logrus.Error(err.Error())
				}
			}
			if nowTime == "23:00" {
				SendMessageToLover("亲爱的，11 点了，该洗漱睡觉了！\n临睡之际送你一句土味情话：", tian.C_saylove)
			}
		}
	}
}
