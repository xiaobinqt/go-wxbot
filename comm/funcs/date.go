package funcs

import (
	"fmt"
	"time"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"go-wxbot/openwechat/comm/global"
)

func ImportDateFormatMsg() (msg string) {
	yuanDan := "01-01"      // 元旦
	valentineDay := "02-14" // 情人节
	anniversary := "09-17"  // 纪念日
	//qiXiDay := "七月初七"       // 七夕，七月初七
	day520 := "05-20"
	//springFestival := "正月初一" // 春节
	//loverBirthday := "七月初十"

	//logrus.Debugf(yuanDan, valentineDay, anniversary, qiXiDay, day520, springFestival, loverBirthday)

	msg = fmt.Sprintf("%s。\n距离元旦【01-01】还有 %d 天\n距离春节还有 %d 天\n距离【02-14】情人节还有 %d 天\n距离【520】还有 %d 天\n距离七夕，七月初七还有 %d 天\n"+
		"距离纪念日还有 %d 天\n距离%s的生日还有 %d 天",
		global.Conf.Keys.RemindMsg,
		GetDiffDaysSolar(getCurrentDate(), yuanDan),
		GetDiffDaysLunar(getCurrentDate(), getLunar2SolarDate(int64(time.Now().Year()), 1, 0), 1, 1),
		GetDiffDaysSolar(getCurrentDate(), valentineDay),
		GetDiffDaysSolar(getCurrentDate(), day520),
		GetDiffDaysLunar(getCurrentDate(), getLunar2SolarDate(int64(time.Now().Year()), 7, 7), 7, 7),
		GetDiffDaysSolar(getCurrentDate(), anniversary), global.Conf.Keys.LoverChName,
		GetDiffDaysLunar(getCurrentDate(), getLunar2SolarDate(int64(time.Now().Year()), 7, 10), 7, 10),
	)

	return msg
}

const DefaultDateFormat = "2006-01-02"

func getYearDay(year int, date string) string {
	return fmt.Sprintf(`%d-%s`, year, date)
}

func getCurrentDate() string {
	return time.Now().Format("01-02")
}

// 获取两个时间相差的天数
func GetDiffDaysSolar(curDate, futureDate string) (dd int) {
	curD, _ := time.ParseInLocation(DefaultDateFormat, getYearDay(time.Now().Year(), curDate), time.Local)
	furD, _ := time.ParseInLocation(DefaultDateFormat, getYearDay(time.Now().Year(), futureDate), time.Local)

	// 如果是负数说明已经过去了，加一年再计算
	dd = int(furD.Sub(curD).Hours() / 24)
	if dd > 0 {
		return dd
	}

	furD, _ = time.ParseInLocation(DefaultDateFormat,
		getYearDay(time.Now().AddDate(1, 0, 0).Year(), futureDate), time.Local)

	return int(furD.Sub(curD).Hours() / 24)
}

func GetDiffDaysLunar(curDate, futureDate string, lunarMonth, lunarDay int64) (dd int) {
	curD, _ := time.ParseInLocation(DefaultDateFormat, getYearDay(time.Now().Year(), curDate), time.Local)
	furD, _ := time.ParseInLocation(DefaultDateFormat, getYearDay(time.Now().Year(), futureDate), time.Local)

	// 如果是负数说明已经过去了，加一年再计算
	//fmt.Println(curD.String(), furD.String(), int(furD.Sub(curD).Hours()/24))
	dd = int(furD.Sub(curD).Hours() / 24)
	if dd > 0 {
		return dd
	}

	futureDate = getLunar2SolarDate(int64(time.Now().AddDate(1, 0, 0).Year()),
		lunarMonth, lunarDay)
	furD, _ = time.ParseInLocation(DefaultDateFormat,
		getYearDay(time.Now().AddDate(1, 0, 0).Year(), futureDate), time.Local)

	return int(furD.Sub(curD).Hours() / 24)
}

// TODO 这里可能不准,需要计算闰月
func getLunar2SolarDate(year, month, day int64) string {
	c := calendar.ByLunar(year, month, day, 0, 0, 0, false)
	return fmt.Sprintf("%02d-%02d", c.Solar.GetMonth(), c.Solar.GetDay())
}

func RemainingDays() int {
	curD, _ := time.ParseInLocation(DefaultDateFormat, time.Now().Format("2006-01-02"), time.Local)
	furD, _ := time.ParseInLocation(DefaultDateFormat, fmt.Sprintf("%d-12-31", time.Now().Year()), time.Local)

	return int(furD.Sub(curD).Hours() / 24)
}
