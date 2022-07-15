package image

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/web"
)

const ImageSourceURL = "https://imgegg.qianxiaoduan.com/wallpaper?offset=%d&limit=3&type=1"
const ImageURL = "https://img.qianxiaoduan.com/"

func GetRandomOffset() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10000)
}

type ImgSourceResp struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Count int `json:"count"`
		Rows  []struct {
			Id        string      `json:"id"`
			Thumbnail string      `json:"thumbnail"`
			Url       string      `json:"url"`
			Hot       int         `json:"hot"`
			Width     int         `json:"width"`
			Height    int         `json:"height"`
			Name      interface{} `json:"name"`
			Type      string      `json:"type"`
			Scale     string      `json:"scale"`
			Tag       string      `json:"tag"`
			CreatedAt string      `json:"createdAt"`
			UpdatedAt string      `json:"updatedAt"`
		} `json:"rows"`
	} `json:"result"`
}

func GetImage() (imgURL string, err error) {
	var (
		surl, respBody string
		statusCode     int
		imgSourceResp  ImgSourceResp
	)
	surl = fmt.Sprintf(ImageSourceURL, GetRandomOffset())
	respBody, statusCode, err = web.HTTP(surl, http.MethodGet, map[string]string{}, 30*time.Second, "")
	if err != nil {
		err = errors.Wrapf(err, "GetImage request error:%s", surl)
		logrus.Error(err.Error())
		return "", err
	}

	if statusCode != http.StatusOK {
		err = fmt.Errorf("GetImage request error:%s, statusCode:%d", surl, statusCode)
		logrus.Error(err.Error())
		return "", err
	}

	err = jsoniter.Unmarshal([]byte(respBody), &imgSourceResp)
	if err != nil {
		err = errors.Wrapf(err, "GetImage unmarshal error:%s", surl)
		logrus.Error(err.Error())
		return "", err
	}

	if len(imgSourceResp.Result.Rows) == 0 {
		err = fmt.Errorf("GetImage Result.Rows len is 0 error:%s, ", surl)
		logrus.Error(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s%s", ImageURL, imgSourceResp.Result.Rows[0].Url), nil
}

func SaveEncourageImg(imgURL string) (savePath string, err error) {
	savePath = fmt.Sprintf("%s/%d%s", os.TempDir(), time.Now().UnixNano(), path.Ext(imgURL))
	logrus.Debugf("SendEncourageImg savePath: %s", savePath)

	request, err := http.NewRequest(http.MethodGet, imgURL, nil)
	if err != nil {
		err = errors.Wrapf(err, "SendEncourageImg newRequest err")
		logrus.Error(err.Error())
		return "", err
	}

	newResp, err := (&http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		},
	}).Do(request)
	if err != nil {
		err = errors.Wrapf(err, "SendEncourageImg client do err")
		logrus.Error(err.Error())
		return "", err
	}
	defer newResp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		err = errors.Wrapf(err, "SendEncourageImg os.Create err")
		logrus.Error(err.Error())
		return "", err
	}

	defer out.Close()

	io.Copy(out, newResp.Body)

	return savePath, nil
}
