package tian

import (
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/web"
)

// 不是 vip 只能多注册几个账号，哈哈哈

func GetMessageV1(stype string, word ...string) (message string, err error) {
	var (
		surl, respBody string
		statusCode     int
		info           Info1
	)

	surl = fmt.Sprintf("http://api.tianapi.com/%s/index?key=%s", stype, global.Conf.Keys.TianapiKey1)

	respBody, statusCode, err = web.HTTP(surl, http.MethodGet, map[string]string{}, 30*time.Second, "")
	if err != nil {
		err = errors.Wrapf(err, "GetMessage http err")
		logrus.Error(err.Error())
		return "", err
	}

	if statusCode != http.StatusOK {
		err = fmt.Errorf("GetMessage http statusCode not 200 is %d", statusCode)
		logrus.Error(err.Error())
		return "", err
	}

	err = jsoniter.Unmarshal([]byte(respBody), &info)
	if err != nil {
		err = errors.Wrapf(err, "GetMessage http Unmarshal err")
		logrus.Error(err.Error())
		return "", err
	}

	if stype == C_lizhiguyan {
		if len(info.Newslist) > 0 && info.Newslist[0].Saying != "" {
			message = fmt.Sprintf("%s【翻译：%s】。", info.Newslist[0].Saying, info.Newslist[0].Transl)
			return message, nil
		}
	}

	err = fmt.Errorf("GetMessage http content empty")
	logrus.Error(err.Error())
	return "", err
}
