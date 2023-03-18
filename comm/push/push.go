package push

import (
	"fmt"
	"github.com/donething/utils-go/dowx"
	"myrouter/config"
	. "myrouter/config"
)

var (
	// WXQiYe 微信推送
	WXQiYe *dowx.QiYe
)

// WXPushCard 推送微信卡片消息
func WXPushCard(title string, description string, url string, btnText string) {
	if !initPush() {
		return
	}

	err := WXQiYe.PushCard(config.Conf.WXPush.Agentid, title, description, Conf.WXPush.ToUser, url, btnText)
	if err != nil {
		fmt.Printf("微信推送消息出错：%s\n", err)
		return
	}
}

// WXPushText 推送微信文本消息
func WXPushText(content string) {
	if !initPush() {
		return
	}
	err := WXQiYe.PushText(Conf.WXPush.Agentid, content, Conf.WXPush.ToUser)
	if err != nil {
		fmt.Printf("微信推送消息出错：%s\n", err)
		return
	}
}

// 初始化
func initPush() bool {
	// 初始化微信推送
	if WXQiYe == nil {
		if Conf.WXPush.Appid == "" || Conf.WXPush.Secret == "" {
			fmt.Printf("微信推送的 Token 为空，无法推送消息\n")
			return false
		}
		WXQiYe = dowx.NewQiYe(Conf.WXPush.Appid, Conf.WXPush.Secret)
	}

	return true
}
