package urltool

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/url"
	"path"
)

func GetBasePath(targetUrl string) (string, error) {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		logx.Errorw("url.Parse failed", logx.LogField{Key: "err", Value: err.Error()})
		return "", err
	}

	if len(myUrl.Host) == 0 {
		return "", errors.New("url.Host is empty")
	}
	return path.Base(myUrl.Path), err
}
