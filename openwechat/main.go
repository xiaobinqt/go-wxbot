package main

import (
	"flag"
	"fmt"

	"github.com/eatmoreapple/openwechat"
	"github.com/json-iterator/go/extra"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	conf2 "go-wxbot/openwechat/comm/conf"
	"go-wxbot/openwechat/comm/global"
	msg2 "go-wxbot/openwechat/comm/msg"
	"go-wxbot/openwechat/comm/ticker"
)

var (
	cfgPath = flag.String("c", "config/prod.yaml", "*.yaml config path")
	err     error
)

func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func main() {
	extra.RegisterFuzzyDecoders()
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("config: %s", *cfgPath)

	// 加载配置文件
	global.Conf, err = conf2.GetConf(*cfgPath)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	bot := openwechat.DefaultBot()
	//bot := openwechat.DefaultBot(openwechat.Normal) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			_, err = msg.ReplyText("pong")
			if err != nil {
				err = errors.Wrapf(err, "ping msg replyText err")
				logrus.Error(err.Error())
			}
		}

		// 处理消息
		msg2.HandleMsg(msg)
	}

	// 注册登陆二维码回调
	//bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	bot.UUIDCallback = ConsoleQrCode

	// 登陆
	if err = bot.Login(); err != nil {
		logrus.Fatalf("bot.Login err %s", err.Error())
	}

	// 获取登陆的用户
	global.WxSelf, err = bot.GetCurrentUser()
	if err != nil {
		logrus.Fatalf("GetCurrentUser err: %s ", err.Error())
	}

	// 获取所有的好友
	global.WxFriends, err = global.WxSelf.Friends(true)
	if err != nil {
		logrus.Fatalf("wx self get friends err: %s ", err.Error())
	}

	// 获取所有的群组
	global.WxGroups, err = global.WxSelf.Groups(true)
	if err != nil {
		logrus.Fatalf("wx self get groups err: %s ", err.Error())
	}

	ticker.Ticker()

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
