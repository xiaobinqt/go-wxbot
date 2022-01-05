package ticker

import (
	"fmt"
	"testing"
)

func TestSendLoveMessage(t *testing.T) {
	var done = make(chan struct{})
	go LoveTicker()

	fmt.Println("done...")
	<-done
}
