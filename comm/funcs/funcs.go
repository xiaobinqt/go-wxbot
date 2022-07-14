package funcs

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Wd() string {
	wd, _ := os.Getwd()
	return wd
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Img2base64(path string) (resultBase64 string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.Wrapf(err, "Img2base64 ReadFile err")
		logrus.Error(err.Error())
		return "", err
	}

	var base64Encoding string

	mimeType := http.DetectContentType(bytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	return base64Encoding, nil
}
