package ticker

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/redis"
)

/**
提醒消息
格式1：+s15:32,消息内容,3,60 // 今天 15:31 提醒我「消息内容」,提醒 3 次每次间隔 60s
格式2: +st20221227 15:35,,消息内容,3,60 // 20221227日 15:35 提醒我「消息内容」,提醒 3 次每次间隔 60s


*/

var (
	ctx = context.Background()
)

const (
	messageType1 = "+s"
	messageType2 = "+st"
	memberKey    = "weixin:schedule.notice"
)

func IsScheduleNotice(message string) bool {
	if strings.HasPrefix(message, messageType1) || strings.HasPrefix(message, messageType2) {
		return true
	}
	return false
}

func parseNoticeMessage(tf string) (count, interval int, startTimestamp int64, message string, err error) {
	if strings.HasPrefix(tf, messageType1) == false && strings.HasPrefix(tf, messageType2) == false { // 格式不正确
		return 0, 0, 0, "", fmt.Errorf("提醒事件格式错误，可以输入「事件提醒」关键字获取帮助。")
	}

	var (
		msgType string
	)

	if tf[0:2] == messageType1 && tf[0:3] != messageType2 {
		tf = tf[2:]
		msgType = messageType1
	} else if tf[0:3] == messageType2 {
		tf = tf[3:]
		msgType = messageType2
	}

	tf = strings.ReplaceAll(tf, "，", ",")
	tf = strings.ReplaceAll(tf, "：", ":")
	tfArr := strings.Split(tf, ",")
	if len(tfArr) < 2 { // 格式不正确
		return 0, 0, 0, "", fmt.Errorf("格式错误，可以输入「事件提醒」关键字获取帮助.")
	}

	stime, message, counts := tfArr[0], tfArr[1], "1"
	intervals := "" // 间隔
	if len(tfArr) >= 3 {
		counts = tfArr[2]
	}
	if len(tfArr) > 3 {
		intervals = tfArr[3]
	}

	interval, _ = strconv.Atoi(intervals)
	count, _ = strconv.Atoi(counts)

	if message == "" {
		return 0, 0, 0, "", fmt.Errorf("提醒消息不能为空")
	}

	if count <= 0 {
		return 0, 0, 0, "", fmt.Errorf("提醒次数最小为 1")
	}

	if count >= 1 && interval < 60 {
		count, interval = 1, 0
	}

	if msgType == messageType1 {
		t, err := time.ParseInLocation("2006-01-02 15:04",
			fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), stime), time.Local)
		if err != nil {
			return 0, 0, 0, "", fmt.Errorf("时间格式错误:%s，可以输入「事件提醒」关键字获取帮助。", tf)
		}
		startTimestamp = t.Unix()
	} else {
		t, err := time.ParseInLocation("20060102 15:04", stime, time.Local)
		if err != nil {
			return 0, 0, 0, "", fmt.Errorf("时间格式错误:%s，可以输入「事件提醒」关键字获取帮助。", tf)
		}
		startTimestamp = t.Unix()
	}

	if startTimestamp < time.Now().Unix()-60 {
		return 0, 0, 0, "", fmt.Errorf("时间错误，提醒时间至少要大与当前时间一分钟")
	}

	return count, interval, startTimestamp, message, nil
}

func formatMember(count, interval int, startTimestamp int64, message, userID string) []*redis2.Z {
	var (
		members = make([]*redis2.Z, 0)
	)

	for i := 0; i < count; i++ {
		tmp := startTimestamp + int64(i*interval)
		members = append(members, &redis2.Z{
			Score:  float64(tmp),
			Member: fmt.Sprintf("%s.placeholder.%s.placeholder.%s", message, userID, uuid.NewV4().String()),
		})
	}

	return members
}

func set(tf, userID string) (replyMsg string) {
	var (
		count, interval int
		startTimestamp  int64
		message         string
		members         = make([]*redis2.Z, 0)
		err             error
	)

	count, interval, startTimestamp, message, err = parseNoticeMessage(tf)
	if err != nil {
		return err.Error()
	}

	members = formatMember(count, interval, startTimestamp, message, userID)

	redisClient := redis.GetRedis()
	if redisClient == nil {
		err = fmt.Errorf("get redis client err")
		logrus.Error(err.Error())
		return err.Error()
	}

	err = redisClient.ZAdd(ctx, memberKey, members...).Err()
	if err != nil {
		err = errors.Wrapf(err, "redis zadd err")
		return err.Error()
	}

	replyMsg = fmt.Sprintf("提醒事件添加成功，会在 %s 提醒：%s",
		time.Unix(startTimestamp, 0).Format("2006-01-02 15:04:05"), message)
	if count > 1 && interval > 0 {
		replyMsg = fmt.Sprintf("%s。一共提醒 %d 次，每次间隔 %d 秒。", replyMsg, count, interval)
	}

	return replyMsg
}

func get(timestamp int64) (msg []string, err error) {
	redisClient := redis.GetRedis()
	if redisClient == nil {
		err = fmt.Errorf("get redis client err")
		logrus.Error(err.Error())
		return nil, err
	}

	// zrangebyscore weixin:schedule.notice -inf (1672300344
	zRangeByScore := redisClient.ZRangeByScore(ctx, memberKey, &redis2.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("(%d", timestamp),
	})
	if zRangeByScore.Err() != nil {
		err = errors.Wrapf(err, "get ZRangeByScore err")
		logrus.Error(err.Error())
		return nil, err
	}

	return zRangeByScore.Val(), nil
}

func del() {
	var (
		err error
	)
	redisClient := redis.GetRedis()
	if redisClient == nil {
		err = fmt.Errorf("get redis client err")
		logrus.Error(err.Error())
		return
	}

	zRemRangeByScore := redisClient.ZRemRangeByScore(ctx, memberKey,
		"0", fmt.Sprintf("(%d", time.Now().Unix()))
	if zRemRangeByScore.Err() != nil {
		err = errors.Wrapf(err, "schedule notice del err")
		logrus.Error(err.Error())
	}

	return
}

func AddScheduleNotice(msg, userID string) (replyMsg string) {
	return set(msg, userID)
}

type Msg struct {
	Message string
	UID     string
}

func ScheduleNoticeTicker() {
	var (
		msg   = make([]string, 0)
		err   error
		doing bool
	)
	for {
		select {
		case t := <-time.After(20 * time.Second):
			if doing {
				continue
			}
			msg, err = get(t.Unix())
			if err != nil {
				doing = false
				continue
			}

			if len(msg) == 0 {
				doing = false
				continue
			}

			msgf := make([]Msg, 0)
			for _, each := range msg {
				tmpArr := strings.Split(each, ".placeholder.")
				if len(tmpArr) < 2 {
					continue
				}
				msgf = append(msgf, Msg{
					Message: tmpArr[0],
					UID:     tmpArr[1],
				})
			}
			if len(msgf) == 0 {
				doing = false
				continue
			}

			// 先删除再发消息，发送时有网络请求会慢
			del()
			go sendNotice(msgf)

			doing = false
		}
	}
}

func sendNotice(msgf []Msg) {
	var (
		err     error
		idMap   = make(map[string]string, 0)
		gMsgMap = make(map[string]*openwechat.Group) // TODO 暂时不考虑同一个群消息重复问题
		fMsgMap = make(map[string]*openwechat.Friend)
	)

	for _, each := range msgf {
		idMap[each.UID] = each.Message
	}

	g, err := global.WxSelf.Groups(true)
	if err == nil {
		for _, each := range g {
			if message, ok := idMap[each.ID()]; ok {
				gMsgMap[message] = each
			}
		}
	}

	f, err := global.WxSelf.Friends(true)
	if err == nil {
		for _, each := range f {
			if message, ok := idMap[each.ID()]; ok {
				fMsgMap[message] = each
			}
		}
	}

	if len(fMsgMap) > 0 {
		for message, each := range fMsgMap {
			_, err = global.WxSelf.SendTextToFriend(each, message)
			if err != nil {
				err = errors.Wrapf(err, "sendNotice SendTextToFriend err")
				logrus.Error(err.Error())
			}
		}
	}

	if len(gMsgMap) > 0 {
		for message, each := range gMsgMap {
			_, err = global.WxSelf.SendTextToGroup(each, message)
			if err != nil {
				err = errors.Wrapf(err, "sendNotice SendTextToGroup err")
				logrus.Error(err.Error())
			}
		}
	}

}
