package models

// ClashConfig 解析 Clash 的配置 `/data/clash/config.yaml`
type ClashConfig struct {
	Port               int      `yaml:"port"`                // 端口, 如 7890
	SocksPort          int      `yaml:"socks-port"`          // Socks端口, 如 7891
	AllowLan           bool     `yaml:"allow-lan"`           // 是否允许局域网, 如"true"
	Mode               string   `yaml:"mode"`                // 模式, 如"Rule"
	LogLevel           string   `yaml:"log-level"`           // 日志级别, 如"info"
	ExternalController string   `yaml:"external-controller"` // 外部控制器, 如":9090"
	Proxies            []Proxy  `yaml:"proxies"`             // 代理列表
	ProxyGroups        []Group  `yaml:"proxy-groups"`        // 代理组列表
	Rules              []string `yaml:"rules"`               // 规则列表, 如 [" - DOMAIN-SUFFIX,acl4.ssr,🎯 全球直连"]
}

type Proxy struct {
	Name    string `yaml:"name"`    // 代理名称, 如"自用（路由器）"
	Server  string `yaml:"server"`  // 服务器地址, 如"example.com"
	Port    int    `yaml:"port"`    // 端口, 如"23199"
	Type    string `yaml:"type"`    // 类型, 如"vmess"
	Uuid    string `yaml:"uuid"`    // UUID, 如"aa7ded50-9508-4bd2-c44e-ebdc5c7aa9bb"
	AlterId int    `yaml:"alterId"` // 额外ID, 如"0"
	Cipher  string `yaml:"cipher"`  // 加密模式, 如"auto"
	TLS     bool   `yaml:"tls"`     // 是否启用TLS, 如"false"
	Network string `yaml:"network"` // 网络类型, 如"ws"
	UDP     bool   `yaml:"udp"`     // 是否启用UDP, 如"true"
}

type Group struct {
	Name     string   `yaml:"name"`     // 组名称, 如"🚀 节点选择"
	Type     string   `yaml:"type"`     // 类型, 如"select"
	Proxies  []string `yaml:"proxies"`  // 代理列表, 如["🔯 故障转移","♻️ 自动选择","自用（路由器）"]
	URL      string   `yaml:"url"`      // URL, 如"http://www.gstatic.com/generate_204"
	Interval int      `yaml:"interval"` // 时间间隔, 如"300"
}
