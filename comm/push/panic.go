package push

// Panic 出错后退出
func Panic(err error) {
	if err != nil {
		WXPushMsg("出现崩溃错误", err.Error())
		panic(err)
	}
}
