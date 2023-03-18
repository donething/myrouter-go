package models

// IPAddrs IP 地址信息
type IPAddrs struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`

	// 所属的路由器
	From string `json:"from"`
}
