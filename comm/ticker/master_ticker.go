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

// 每天提醒自己一些事
func MasterTicker() {
	for {
		select {
		case t := <-time.After(1 * time.Minute):
			nowTime := t.Format("15:04")

			yasiEnd, _ := time.ParseInLocation("2006-01-02", "2023-10-07", time.Local)
			yasiRemaindays := int(yasiEnd.Sub(t).Hours() / 24)

			if nowTime == "10:00" {
				lz, err := tian.GetMessageV1(tian.C_lizhiguyan)
				message := ""
				if err != nil {
					message = fmt.Sprintf("盛年不重来，一日难再晨。及时当勉励，岁月不待人。\n今年还剩 %d 天。", funcs.RemainingDays())
				} else {
					message = fmt.Sprintf("今年还剩 %d 天。\n\n%s", funcs.RemainingDays(), lz)
				}

				err = global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.MasterAccount).
					SendText(message)
				if err != nil {
					err = errors.Wrapf(err, "SendMessageToMasterAccout err")
					logrus.Error(err.Error())
				}

				_ = global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.MasterAccount).
					SendText(fmt.Sprintf(`离雅思过期时间还有 %d 天，兄弟，留给你的时间不多了！`, yasiRemaindays))
			}

			if nowTime == "22:00" {
				message := "记得背单词兄弟，别一天天的想偷懒！"
				err := global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.MasterAccount).
					SendText(message)
				if err != nil {
					err = errors.Wrapf(err, "SendMessageToMasterAccout err")
					logrus.Error(err.Error())
				}
			}

			if nowTime == "23:00" {
				message := "休息一下，整理一下今天的账单吧！记日记的时间也到了，不要忘记了哦！"
				err := global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.MasterAccount).
					SendText(message)
				if err != nil {
					err = errors.Wrapf(err, "SendMessageToMasterAccout err")
					logrus.Error(err.Error())
				}

				_ = global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.MasterAccount).
					SendText(fmt.Sprintf(`离雅思过期时间还有 %d 天，兄弟，留给你的时间不多了！`, yasiRemaindays))
			}

			if nowTime == "23:30" {
				message := funcs.ImportDateFormatMsg()
				logrus.Infof("send remind msg: %s", message)
				err := global.WxFriends.
					SearchByRemarkName(1, global.Conf.Keys.MasterAccount).
					SendText(message)
				if err != nil {
					err = errors.Wrapf(err, "SendMessageToMasterAccout err")
					logrus.Error(err.Error())
				}
			}

		}
	}
}

func KeepLive() {
	for {
		select {
		case <-time.After(10 * time.Minute):
			global.WxSelf.FileHelper().SendText(fmt.Sprintf(`保活: %s`,
				time.Now().Format("2006-01-02 15:04:05")))
		}
	}
}
