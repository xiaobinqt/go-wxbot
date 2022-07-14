package global

import (
	"github.com/eatmoreapple/openwechat"
	"go-wxbot/openwechat/comm/conf"
)

var (
	Conf      *conf.Conf
	WxSelf    *openwechat.Self
	WxFriends openwechat.Friends // 可能有缓存
	WxGroups  openwechat.Groups  // 可能有缓存
)
