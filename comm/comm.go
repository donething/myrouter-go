package comm

import "myrouter/push"

// Panic 出错后退出
func Panic(err error) {
	if err != nil {
		push.WXPushCard("[路由器] 发生崩溃错误", err.Error(), "", "")
		panic(err)
	}
}
