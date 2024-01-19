package config

import "encoding/json"

// Logo 路由器标识（名子/型号）
type Logo string

// Auths Bearer Authentication 提示用户正确输入
type Auths map[string]string

func (a *Auths) UnmarshalJSON(data []byte) error {
	// 需要定义别名，避免无限循环一直解析
	type Alias Auths

	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// 没有值时设为提示值，有值时解析出值
	if alias == nil {
		*a = map[string]string{"the_site_name": "example_32bit_bearer_authentication"}
	} else {
		*a = Auths(alias)
	}

	return nil
}

// Config 配置
type Config struct {
	// 操作授权码
	Auths Auths `json:"auths"`

	// 路由器网关地址，如"192.168.0.1"
	Gateway string `json:"gateway"`

	// clash 的辅助功能
	Clash struct {
		// 自定义规则文件的路径。为空时默认为"/data/clash/yamls/rules.yaml"
		RulesPath string `json:"rulesPath"`
		// config.yaml 配置文件的路径。为空时默认为"/data/clash/config.yaml"
		ConfigPath string `json:"configPath"`
	}

	// 管理员账号
	Router struct {
		// Logo 运行所在的路由器标识，用于根据路由器生成不同的实例
		// 可选择 routers 包中的 "RedmiAX6000"、"JDC"等。可留空，默认值查看 config.go
		Logo Logo `json:"logo"`

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

	// 发送IP到远程服务器的URL，如 "https://example.com/api/router/ip/update"
	// 本地测试可设为后台地址，如"http://127.0.0.1:1234/api/router/ip/update"
	Remote struct {
		UpdateIPAddr string `json:"updateIPAddr"`
	} `json:"remote"`
}
