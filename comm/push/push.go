package push

import (
	"fmt"
	"github.com/donething/utils-go/dowx"
	"myrouter/comm/logger"
	. "myrouter/config"
)

// TAG 发送消息时的来源提示
var TAG = fmt.Sprintf("路由器[%s]", Conf.Router.Logo)

// WXQiYe 微信推送
var WXQiYe *dowx.QiYe

// WXPushCard 推送微信卡片消息
func WXPushCard(title string, description string, url string, btnText string) {
	if !initPush() {
		return
	}

	t := fmt.Sprintf("%s %s", TAG, title)
	err := WXQiYe.PushCard(Conf.WXPush.Agentid, t, description, Conf.WXPush.ToUser, url, btnText)
	if err != nil {
		logger.Error.Printf("微信推送消息出错：%s。消息标题：'%s'\n", err, title)
		return
	}
}

// WXPushMsg 推送微信文本消息
func WXPushMsg(title, content string) {
	if !initPush() {
		return
	}

	t := fmt.Sprintf("%s %s", TAG, title)
	err := WXQiYe.PushTextMsg(Conf.WXPush.Agentid, t, content, Conf.WXPush.ToUser)
	if err != nil {
		logger.Error.Printf("微信推送消息出错：%s。消息标题：'%s'\n", err, title)
		return
	}
}

// 初始化
func initPush() bool {
	// 初始化微信推送
	if WXQiYe == nil {
		if Conf.WXPush.Appid == "" || Conf.WXPush.Secret == "" {
			logger.Warn.Printf("无法推送消息：没有设置企业微信的 token\n")
			return false
		}

		WXQiYe = dowx.NewQiYe(Conf.WXPush.Appid, Conf.WXPush.Secret)
	}

	return true
}
