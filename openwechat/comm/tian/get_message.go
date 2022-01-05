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

func GetMessage(stype string, word ...string) (message string, err error) {
	var (
		surl, respBody string
		statusCode     int
		info           Info1
	)

	if stype == C_caipu {
		surl = fmt.Sprintf("http://api.tianapi.com/%s/index?key=%s&word=%s",
			stype, global.Conf.Keys.TianapiKey, word[0])
	} else {
		surl = fmt.Sprintf("http://api.tianapi.com/%s/index?key=%s", stype, global.Conf.Keys.TianapiKey)
	}

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

	if stype == C_godreply { // 神回复比较特殊一点
		if len(info.Newslist) > 0 && info.Newslist[0].Content != "" && info.Newslist[0].Title != "" {
			message = fmt.Sprintf(`
%s
%s
`,
				info.Newslist[0].Title,
				info.Newslist[0].Content,
			)
			return message, nil
		}
	} else if stype == C_caipu { // 菜谱
		if len(info.Newslist) > 0 && info.Newslist[0].Zuofa != "" && info.Newslist[0].Yuanliao != "" {
			message = fmt.Sprintf(`
原料 ：%s
%s做法 ：%s
`,
				info.Newslist[0].Yuanliao,
				"\n",
				info.Newslist[0].Zuofa,
			)
			return message, nil
		}
		return "", ErrNotfoundCaiPu
	} else if stype == C_englishSentence { // 英语一句话
		if len(info.Newslist) > 0 && info.Newslist[0].Zh != "" && info.Newslist[0].En != "" {
			message = fmt.Sprintf(`
%s
%s
`,
				info.Newslist[0].En,
				info.Newslist[0].Zh,
			)
			return message, nil
		}
	} else {
		if len(info.Newslist) > 0 && info.Newslist[0].Content != "" {
			return info.Newslist[0].Content, nil
		}
	}

	err = fmt.Errorf("GetMessage http content empty")
	logrus.Error(err.Error())
	return "", err
}
