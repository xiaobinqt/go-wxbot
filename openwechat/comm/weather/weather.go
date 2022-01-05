package weather

import (
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	adcode2 "go-wxbot/openwechat/comm/adcode"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/web"
)

// https://restapi.amap.com/v3/weather/weatherInfo?city=110101&key=<用户key>
func GetWeatherInfo(cityName string) (info WeatherInfo, err error) {
	var (
		surl,
		respBody, adcode string
		statusCode int
	)

	adcode = adcode2.GetAdcodeByCityName(cityName)
	if adcode == "" {
		return info, fmt.Errorf("get adcode empty,cityName: %s", cityName)
	}

	surl = fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s",
		adcode, global.Conf.Keys.WeatherKey)

	respBody, statusCode, err = web.HTTP(surl, http.MethodGet, map[string]string{}, 30*time.Second, "")
	if err != nil {
		err = errors.Wrapf(err, "GetWeatherInfo http err")
		logrus.Error(err.Error())
		return info, err
	}

	if statusCode != http.StatusOK {
		return info, fmt.Errorf("GetWeatherInfo http statusCode not 200 is %d", statusCode)
	}

	err = jsoniter.Unmarshal([]byte(respBody), &info)
	if err != nil {
		err = errors.Wrapf(err, "GetWeatherInfo Unmarshal err")
		return info, err
	}

	if info.Info != "OK" {
		err = fmt.Errorf("GetWeatherInfo resp info not OK is %s", info.Info)
		return info, err
	}

	return info, nil
}

func GetFormatWeatherMessage(cityName string) (format string, err error) {
	var (
		info WeatherInfo
	)
	info, err = GetWeatherInfo(cityName)
	if err != nil {
		return "", err
	}

	format = fmt.Sprintf(`
%s%s今日天气 %s，温度 %s 摄氏度，空气湿度 %s。
`,
		info.Lives[0].Province,
		info.Lives[0].City,
		info.Lives[0].Weather,
		info.Lives[0].Temperature,
		info.Lives[0].Humidity,
	)

	return format, nil
}
