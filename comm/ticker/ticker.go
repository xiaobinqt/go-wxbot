package ticker

func Ticker() {
	go LoveTicker()
	//go FansTicker()
	go BubeiGroupTicker()
	go EncourageTicker()
	go MasterTicker()
	go ScheduleNoticeTicker()
}
