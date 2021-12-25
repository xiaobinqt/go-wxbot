package msg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/comm/funcs"
	"go-wxbot/comm/global"
	"go-wxbot/comm/tian"
	"go-wxbot/comm/weather"
	"go-wxbot/comm/web"
)

func HandleMsg(msg *openwechat.Message) {
	if msg.IsSendBySelf() { // 自己的消息不处理
		return
	}

	var (
		contentText = ""
		err         error
	)

	if msg.IsText() { // 处理文本消息
		contentText = trimMsgContent(msg.Content)
		if contentText != "打赏" && contentText != "圣诞帽" {
			reply := contextTextBypass(contentText)
			reply = strings.TrimLeft(reply, "\n")
			reply = strings.TrimRight(reply, "\n")
			_, err = msg.ReplyText(reply)
			if err != nil {
				err = errors.Wrapf(err, "reply text msg err,contentText: %s", contentText)
				logrus.Error(err.Error())
			}
		}

		handleImageReply(msg, contentText)

	}

}

func handleImageReply(msg *openwechat.Message, txt string) {
	if txt == "打赏" {
		img, err := os.Open("reword.png")
		defer img.Close()
		if err != nil {
			err = errors.Wrapf(err, "reword open file err")
			logrus.Error(err.Error())
			_, err = msg.ReplyText("学雷锋，视钱财如粪土，不用打赏。")
			handleErr(err)
			return
		}

		_, err = msg.ReplyImage(img)
		handleErr(err)
		return
	}

	if txt == "圣诞帽" {
		handleChristmasHatMsg(msg)
	}

}

func handleChristmasHatMsg(msg *openwechat.Message) {
	var (
		sender *openwechat.User
		err    error
		avatarPath,
		avatarBase64, base64Hat, base64SaveName string
	)

	// 保存用户的头像
	sender, err = msg.SenderInGroup()
	if err != nil {
		err = errors.Wrapf(err, "SenderInGroup err")
		logrus.Error(err.Error())
		msg.ReplyText(fmt.Sprintf("%s处理不过来了，过会儿再来生成圣诞帽吧！", global.Conf.Keys.BotName))
		return
	}

	avatarPath = avatarSavePath(sender.NickName)
	logrus.Debugf("avatarPath: %s ", avatarPath)
	err = sender.SaveAvatar(avatarPath)
	if err != nil {
		err = errors.Wrapf(err, "SaveAvatar err:%s", sender.NickName)
		logrus.Error(err.Error())
		msg.ReplyText(fmt.Sprintf("%s处理不过来了，过会儿再来生成圣诞帽吧！", global.Conf.Keys.BotName))
		return
	}

	// 头像转成 base64
	avatarBase64, err = funcs.Img2base64(avatarPath)
	if err != nil {
		err = errors.Wrapf(err, "Img2base64 err:%s,path: %s", sender.NickName, avatarPath)
		logrus.Error(err.Error())
		msg.ReplyText(fmt.Sprintf("%s处理不过来了，过会儿再来生成圣诞帽吧！", global.Conf.Keys.BotName))
		return
	}

	logrus.Debugf("Img2base64 avatarBase64: %s \n", avatarBase64)

	// 调用接口把 base64 头像加个圣诞帽
	base64Hat, err = AvatarAddChristmasHat(avatarBase64)
	if err != nil {
		msg.ReplyText(fmt.Sprintf("%s处理不过来了，过会儿再来生成圣诞帽吧！", global.Conf.Keys.BotName))
		return
	}

	// 将返回的图片存到本地
	base64SaveName = fmt.Sprintf("%s_base64", url.QueryEscape(sender.NickName))
	filename, err := SaveImageToDisk(base64SaveName, base64Hat)
	if err != nil {
		err = errors.Wrapf(err, "saveImageToDisk err:%s,path: %s", sender.NickName, avatarPath)
		logrus.Error(err.Error())
		msg.ReplyText(fmt.Sprintf("%s处理不过来了，过会儿再来生成圣诞帽吧！", global.Conf.Keys.BotName))
		return
	}

	// 发送加了圣诞帽的头像
	avatarBase64Path := fmt.Sprintf("%s/avatar/%s", funcs.Wd(), filename)
	img, err := os.Open(avatarBase64Path)
	defer img.Close()
	if err != nil {
		err = errors.Wrapf(err, "avatarBase64Path open err:%s,path: %s", sender.NickName, avatarPath)
		logrus.Error(err.Error())
		msg.ReplyText(fmt.Sprintf("%s处理不过来了，过会儿再来生成圣诞帽吧！", global.Conf.Keys.BotName))
		return
	}

	msg.ReplyImage(img)
}

type ChristmasHatReq struct {
	Base64 string `json:"base64"`
}

type AvatarResp struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Msg     string `json:"msg"`
}

func AvatarAddChristmasHat(avatarBase64 string) (hatBase64 string, err error) {
	var (
		req        ChristmasHatReq
		surl       string
		reqBytes   = make([]byte, 0)
		respBody   string
		statusCode int
		hatRet     AvatarResp
	)
	req.Base64 = avatarBase64
	surl = fmt.Sprintf("%s", global.Conf.Keys.ChristmasHatURL)
	reqBytes, _ = jsoniter.Marshal(req)

	respBody, statusCode, err = web.HTTP(surl, http.MethodGet, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}, 30*time.Second, string(reqBytes))
	if err != nil {
		err = errors.Wrapf(err, "avatarAddChristmasHat http err")
		logrus.Error(err.Error())
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("avatarAddChristmasHat statusCode not 200 is %d", statusCode)
	}

	err = jsoniter.Unmarshal([]byte(respBody), &hatRet)
	if err != nil {
		err = errors.Wrapf(err, "avatarAddChristmasHat Unmarshal err")
		logrus.Error(err.Error())
		return "", err
	}

	if hatRet.Success == false {
		return "", fmt.Errorf("avatarAddChristmasHat success not true is %t", hatRet.Success)
	}

	return hatRet.Data, nil
}

func avatarSavePath(nickname string) (path string) {
	return fmt.Sprintf("%s/avatar/%s.png", funcs.Wd(), url.QueryEscape(nickname))
}

func handleErr(err error, grep ...string) {
	prefix := ""
	if len(grep) > 0 {
		for _, each := range grep {
			prefix = fmt.Sprintf("%s %s", prefix, each)
		}
	}
	if err != nil {
		logrus.Errorf("%s [%s]", prefix, err.Error())
	}
}

func trimMsgContent(content string) string {
	content = strings.TrimLeft(content, " ")
	content = strings.TrimRight(content, " ")
	return content
}

func SaveImageToDisk(saveName, data string) (filename string, err error) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return "", fmt.Errorf("InvalidImage")
	}
	ImageType := data[11:idx]
	log.Println(ImageType)

	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		return "", fmt.Errorf("Cannot decode b64")
	}
	r := bytes.NewReader(unbased)
	switch ImageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			return "", fmt.Errorf("Bad png")
		}

		filename = fmt.Sprintf("%s.png", saveName)
		f, err := os.OpenFile(fmt.Sprintf("%s/avatar/%s", funcs.Wd(), filename), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", fmt.Errorf("Cannot open file")
		}

		png.Encode(f, im)
	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			return "", fmt.Errorf("Bad jpeg")
		}

		filename = fmt.Sprintf("%s.jpeg", saveName)
		f, err := os.OpenFile(fmt.Sprintf("%s/avatar/%s", funcs.Wd(), filename), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", fmt.Errorf("Cannot open file")
		}

		jpeg.Encode(f, im, nil)
	case "gif":
		im, err := gif.Decode(r)
		if err != nil {
			return "", fmt.Errorf("Bad gif")
		}

		filename = fmt.Sprintf("%s.gif", saveName)
		f, err := os.OpenFile(fmt.Sprintf("%s/avatar/%s", funcs.Wd(), filename), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", fmt.Errorf("Cannot open file")
		}

		gif.Encode(f, im, nil)
	}

	return filename, nil
}

func contextTextBypass(txt string) (retMsg string) {
	var (
		err error
	)
	if txt == "菜单" {
		return `
天气查询，如：泾县天气。
菜谱查询，如: 红烧肉菜谱，红烧肉做法。
输入【打赏】打赏卫小兵。
输入【毒鸡汤】关键字回复毒鸡汤。
输入【圣诞帽】关键字回复简单处理后的圣诞帽头像，个别用户获取不到头像信息。
输入【英语一句话】关键字回复一句学习英语。
`
	}

	if txt == "天气" {
		return "支持天气查询，如: 泾县天气。"
	}

	if txt == "菜谱" || txt == "做法" {
		return "支持菜谱查询，如: 红烧肉菜谱，红烧肉做法。"
	}

	// 天气处理
	if strings.HasSuffix(txt, "天气") {
		return handleWeatherMsg(txt)
	}

	// 毒鸡汤处理
	if txt == "毒鸡汤" {
		retMsg, err = tian.GetMessage(tian.C_dujitang)
		if err != nil {
			logrus.Error(err.Error())
			return ""
		}
		return retMsg
	}

	// 菜谱处理
	if strings.HasSuffix(txt, "菜谱") || strings.HasSuffix(txt, "做法") {
		return handleCookbookMsg(txt)
	}

	// 英语一句话
	if txt == "英语一句话" {
		retMsg, err = tian.GetMessage(tian.C_englishSentence)
		if err != nil {
			logrus.Error(err.Error())
			return ""
		}
		return retMsg
	}

	// todo 其他的一些

	return ""
}

func reword() (img *os.File, err error) {
	img, err = os.Open("reword.png")
	defer img.Close()
	return img, err
}

func handleCookbookMsg(txt string) (cookbook string) {
	var (
		err error
	)
	originTxt := txt

	txt = strings.ReplaceAll(txt, "做法", "")
	txt = strings.ReplaceAll(txt, "菜谱", "")
	cookbook, err = tian.GetMessage(tian.C_caipu, txt)
	if err != nil && err != tian.ErrNotfoundCaiPu {
		logrus.Error(err.Error())
		return ""
	}

	if err == tian.ErrNotfoundCaiPu {
		return fmt.Sprintf("暂未找到%s，请重新输入关键字查询", originTxt)
	}

	return cookbook
}

func handleWeatherMsg(txt string) (formatWeatherMsg string) {
	var (
		err error
	)
	originTxt := txt
	txt = strings.ReplaceAll(txt, "天气", "")
	formatWeatherMsg, err = weather.GetFormatWeatherMessage(txt)
	if err != nil {
		err = errors.Wrapf(err, "handleWeatherMsg GetFormatWeatherMessage err")
		logrus.Error(err.Error())
		return fmt.Sprintf("未查询到%s，请检查城市输入是否正确，当前只支持到区/县一级的城市查询，如：泾县天气，新市区天气。", originTxt)
	}

	return formatWeatherMsg
}
