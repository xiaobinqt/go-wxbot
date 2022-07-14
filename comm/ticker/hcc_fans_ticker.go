package ticker

import (
	"fmt"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/tian"
)

// 后厂村吴彦祖粉丝团

func SendMessageToFans(prefix, stype string) {
	var (
		err     error
		message string
		groups  openwechat.Groups
	)
	message, err = tian.GetMessage(stype)
	if err != nil {
		err = errors.Wrapf(err, "SendMessageToFans get message err")
		logrus.Error(err.Error())
		return
	}

	message = fmt.Sprintf("%s%s", prefix, message)

	groups, err = global.WxSelf.Groups(true)
	if err != nil {
		err = errors.Wrapf(err, "SendMessageToFans get groups err")
		logrus.Error(err.Error())
		return
	}

	err = groups.SearchByNickName(1, global.Conf.Keys.HouchangcunFans).SendText(message)
	if err != nil {
		err = errors.Wrapf(err, "SendMessageToFans to groups err, group nickname: %s",
			global.Conf.Keys.HouchangcunFans)
		logrus.Error(err.Error())
	}

	err = groups.SearchByNickName(1, global.Conf.Keys.BanzhuanGroup).SendText(message)
	if err != nil {
		err = errors.Wrapf(err, "SendMessageToFans to groups err, group nickname: %s",
			global.Conf.Keys.HouchangcunFans)
		logrus.Error(err.Error())
	}
}

func FansTicker() {
	for {
		select {
		case t := <-time.After(5 * time.Second):
			nowTime := t.Format("15:04")
			if nowTime == "09:30" {
				pp := ""
				if t.Weekday() >= 1 && t.Weekday() <= 5 {
					pp = fmt.Sprintf(`星期%s快乐，不要忘记上班签到哦~`, weekdayCn(int(t.Weekday())))
				} else {
					pp = fmt.Sprintf(`星期%s快乐，如果今天你得了福报要加班的话，不要忘记签到哦~`,
						weekdayCn(int(t.Weekday())))
				}
				prefix := fmt.Sprintf("%s\n新的一天从一碗毒鸡汤开始：", pp)
				SendMessageToFans(prefix, tian.C_dujitang)
			}
		}
	}
}

// 不背单词打卡群
func BubeiGroupTicker() {
	for {
		select {
		case t := <-time.After(5 * time.Second):
			nowTime := t.Format("15:04")
			if nowTime == "22:30" {
				// 获取群列表
				groups, err := global.WxSelf.Groups(true)
				if err != nil {
					err = errors.Wrapf(err, "BubeiGroupTicker get groups err")
					logrus.Error(err.Error())
					continue
				}
				// 搜索群
				for _, group := range groups {
					if strings.Contains(group.NickName, global.Conf.Keys.BubeiGroup) {
						group.SendText("22:30 了，没打卡的小伙伴，赶紧去打卡吧！")
					}
				}
			}
		}
	}
}

func weekdayCn(i int) string {
	var m = map[int]string{
		0: "日",
		1: "一",
		2: "二",
		3: "三",
		4: "四",
		5: "五",
		6: "六",
	}

	return m[i]
}
