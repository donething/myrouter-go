package config

import "myrouter/models"

// Config 配置
type Config struct {
	// 操作验证码
	Auth string `json:"auth"`

	// 路由器网关地址，如"192.168.0.1"
	Gateway string `json:"gateway"`

	// 管理员账号
	Router struct {
		// From 所属的路由器。用于根据路由器生成不同的对象
		// 可选择 "From"、"JDC"。留空时默认为 "JDC"
		From models.Router `json:"from"`

		Username string `json:"username"`
		Passwd   string `json:"passwd"`
	} `json:"router"`

	// 连接以执行远程命令的端口和密码
	Shell struct {
		// 连接端口。默认 23040
		Port   int    `json:"port"`
		Passwd string `json:"passwd"`
	} `json:"shell"`

	// 微信推送
	WXPush struct {
		Appid   string `json:"appid"`   // 组织 ID
		Secret  string `json:"secret"`  // 秘钥
		Agentid int    `json:"agentid"` // 应用（频道） ID
		ToUser  string `json:"toUser"`  // 接收者。默认 "@all"
	} `json:"wxPush"`

	// 网络唤醒
	WOL struct {
		// 需要网络唤醒的 Mac 地址，如"89:0A:CD:EF:00:12"、"01-23-45-56-67-89"
		MACAddr string `json:"macAddr"`
	} `json:"wol"`

	// 保存 IP 地址的远程服务器。如 "https://example.com/api/router/ip/update"
	// 本地测试可设为后台地址，如"http://127.0.0.1:1234/api/router/ip/update"
	Remote struct {
		UpdateIPAddr string `json:"updateIPAddr"`
	} `json:"remote"`
}
