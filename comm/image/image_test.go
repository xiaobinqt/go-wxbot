package image

import "testing"

func TestGetImage(t *testing.T) {
	imgURL, err := GetImage()
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(imgURL)
}

func TestSendEncourageImg(t *testing.T) {
	path, err := SaveEncourageImg("https://img.qianxiaoduan.com/image/wallpaper/5126.jpg")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(path)
}
