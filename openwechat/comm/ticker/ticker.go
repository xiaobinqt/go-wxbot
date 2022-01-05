package ticker

func Ticker() {
	go LoveTicker()
	go FansTicker()
}
