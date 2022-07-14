package qweather

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/web"
)

const QWeatherHOST = "https://geoapi.qweather.com"

type QWeatherResp struct {
	Code     string `json:"code"`
	Location []struct {
		Name      string `json:"name"`
		Id        string `json:"id"`
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Adm2      string `json:"adm2"`
		Adm1      string `json:"adm1"`
		Country   string `json:"country"`
		Tz        string `json:"tz"`
		UtcOffset string `json:"utcOffset"`
		IsDst     string `json:"isDst"`
		Type      string `json:"type"`
		Rank      string `json:"rank"`
		FxLink    string `json:"fxLink"`
	} `json:"location"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

func GetLocationID(cityName string) (id string, err error) {
	var (
		respBody, surl string
		statusCode     int
		params         = url.Values{}
		qweatherResp   QWeatherResp
	)

	params.Set("location", cityName)
	params.Set("key", global.Conf.Keys.QweatherKey)
	params.Set("range", "cn")
	//params.Set("adm", cityName)

	surl = fmt.Sprintf("%s/v2/city/lookup?%s", QWeatherHOST, params.Encode())
	respBody, statusCode, err = web.HTTP(surl, http.MethodGet, map[string]string{}, 10*time.Second, "")

	if err != nil {
		err = errors.Wrapf(err, "GetLocationID HTTP error")
		logrus.Error(err.Error())
		return "", err
	}

	if statusCode != http.StatusOK {
		err = fmt.Errorf("GetLocationID HTTP status code error: %d", statusCode)
		logrus.Error(err.Error())
		return "", err
	}

	err = jsoniter.Unmarshal([]byte(respBody), &qweatherResp)
	if err != nil {
		err = errors.Wrapf(err, "GetLocationID jsoniter.Unmarshal error")
		logrus.Error(err.Error())
		return "", err
	}

	if qweatherResp.Code != "200" {
		err = fmt.Errorf("GetLocationID qweatherResp.Code error: %s", qweatherResp.Code)
		logrus.Error(err.Error())
		return "", err
	}

	if len(qweatherResp.Location) == 0 {
		err = fmt.Errorf("GetLocationID qweatherResp.Location empyt")
		logrus.Error(err.Error())
		return "", err
	}

	// 匹配
	for _, v := range qweatherResp.Location {
		if v.Name == cityName {
			return v.Id, nil
		}
		if v.Adm2 == cityName {
			return v.Id, nil
		}
		if strings.Contains(v.Adm1, cityName) {
			return v.Id, nil
		}
	}

	return "", fmt.Errorf("GetLocationID not found")
}

type QWeatherDetailResp struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Wind360   string `json:"wind360"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		WindSpeed string `json:"windSpeed"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Pressure  string `json:"pressure"`
		Vis       string `json:"vis"`
		Cloud     string `json:"cloud"`
		Dew       string `json:"dew"`
	} `json:"now"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

func GetQWeatherDetail(cityID, cityName string) (detail string, err error) {
	var (
		respBody, surl     string
		statusCode         int
		params             = url.Values{}
		qweatherDetailResp QWeatherDetailResp
	)

	params.Set("location", cityID)
	params.Set("key", global.Conf.Keys.QweatherKey)

	surl = fmt.Sprintf("https://devapi.qweather.com/v7/weather/now?%s", params.Encode())
	respBody, statusCode, err = web.HTTP(surl, http.MethodGet, map[string]string{}, 10*time.Second, "")

	if err != nil {
		err = errors.Wrapf(err, "GetQWeatherDetail HTTP error")
		logrus.Error(err.Error())
		return "", err
	}

	if statusCode != http.StatusOK {
		err = fmt.Errorf("GetQWeatherDetail HTTP status code error: %d", statusCode)
		logrus.Error(err.Error())
		return "", err
	}

	err = jsoniter.Unmarshal([]byte(respBody), &qweatherDetailResp)
	if err != nil {
		err = errors.Wrapf(err, "GetQWeatherDetail jsoniter.Unmarshal error")
		logrus.Error(err.Error())
		return "", err
	}

	if qweatherDetailResp.Code != "200" {
		err = fmt.Errorf("GetQWeatherDetail qweatherResp.Code error: %s", qweatherDetailResp.Code)
		logrus.Error(err.Error())
		return "", err
	}

	detail = fmt.Sprintf(`%s今天天气，温度 %s 度，%s，%s %s 级，相对湿度 %s`, cityName,
		qweatherDetailResp.Now.Temp, qweatherDetailResp.Now.Text,
		qweatherDetailResp.Now.WindDir, qweatherDetailResp.Now.WindScale,
		qweatherDetailResp.Now.Humidity,
	)

	return detail, nil
}
