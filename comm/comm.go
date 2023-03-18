package comm

import (
	"github.com/donething/utils-go/dohttp"
	"myrouter/comm/push"
)

var Client = dohttp.New(true, false)

// Panic 出错后退出
func Panic(err error) {
	if err != nil {
		push.WXPushCard("[路由器] 发生崩溃错误", err.Error(), "", "")
		panic(err)
	}
}
