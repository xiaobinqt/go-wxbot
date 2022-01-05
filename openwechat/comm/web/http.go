package web

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// HTTP .
func HTTP(reqURL, method string, header map[string]string,
	timeout time.Duration, body string) (respBody string, statusCode int, err error) {
	timeBegin := time.Now()

	defer func() {
		logrus.Infof("url:%s, cost:%d ms, body:%s", reqURL, time.Since(timeBegin).Milliseconds(), body)
	}()

	logrus.Infof("req:%s", reqURL)

	newReq, err := http.NewRequest(method, reqURL, strings.NewReader(body))
	if err != nil {
		err = errors.Wrapf(err, "NewRequest error:%s", reqURL)
		return "", http.StatusInternalServerError, err
	}

	if header != nil {
		for k, v := range header {
			if strings.EqualFold(k, "host") {
				newReq.Host = v
			}
			newReq.Header.Set(k, v)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	newReq = newReq.WithContext(ctx)
	logrus.Tracef("newReq:%+v", newReq)

	// 忽略对证书的校验
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}

	newResp, err := (&http.Client{
		Transport: tr,
	}).Do(newReq)
	if err != nil {
		err = errors.Wrapf(err, "request error:%s", reqURL)
		return "", http.StatusInternalServerError, err
	}
	defer newResp.Body.Close()

	newBody, err := ioutil.ReadAll(newResp.Body)
	return string(newBody), newResp.StatusCode, err
}
