package connect

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		// 禁止KeepAlives，避免长连接
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

func Get(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorw("connect client.Get failed", logx.LogField{Key: "err", Value: err.Error()})
		return false
	}
	// 断开连接
	err = resp.Body.Close()
	if err != nil {
		logx.Errorw("connect resp.Body.Close() failed", logx.LogField{Key: "err", Value: err.Error()})
		return false
	}
	return resp.StatusCode == http.StatusOK // 必须得是200，其他的状态码都算失败
}
