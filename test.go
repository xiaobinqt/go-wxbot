package main

import (
	"fmt"

	"go-wxbot/openwechat/comm/global"
)

func Test() {
	groups, _ := global.WxSelf.Groups(true)
	for _, each := range groups {
		fmt.Println(each.NickName, each.UserName)
	}
}
