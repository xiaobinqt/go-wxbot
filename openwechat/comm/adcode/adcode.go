package adcode

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	adcodeMap = make(map[string]string) // map[adcode]LocInfo
)

func loadAdcodeInfo() {
	// 读取文件
	f, err := os.Open("comm/adcode/adcode.txt")
	if err != nil {
		err = errors.Wrapf(err, "open adcode.txt err")
		logrus.Error(err.Error())
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		line = strings.TrimRight(line, "\n")
		if err != nil || io.EOF == err {
			break
		}
		adcodeArr := strings.Split(line, ",")
		if len(adcodeArr) < 9 {
			continue
		}

		adcodeMap[adcodeArr[1]] = adcodeArr[0]
		if adcodeArr[8] != "" {
			adcodeMap[adcodeArr[8]] = adcodeArr[0]
		} else {
			adcodeMap[adcodeArr[7]] = adcodeArr[0]
		}
	}

	return
}

// GetAdcodeByCityName 通过城市获取 adcode, 支持到区县和简称，比如 山亭区/山亭
func GetAdcodeByCityName(cityName string) (adcode string) {
	if len(adcodeMap) == 0 {
		loadAdcodeInfo()
	}

	return adcodeMap[cityName]
}
