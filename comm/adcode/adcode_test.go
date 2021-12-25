package adcode

import (
	"fmt"
	"testing"
)

func TestLoadAdcodeInfo(t *testing.T) {
	loadAdcodeInfo()
}

func TestGetAdcodeByCityName(t *testing.T) {
	fmt.Println("山亭区 adcode = ", GetAdcodeByCityName("山亭区"))
	fmt.Println("泾县 adcode = ", GetAdcodeByCityName("泾县"))
	fmt.Println("鱼台 adcode = ", GetAdcodeByCityName("鱼台"))
}
