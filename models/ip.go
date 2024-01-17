package models

import (
	"myrouter/config"
)

// IPAddr IP 地址信息
type IPAddr struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`

	// 所属的路由器
	From config.Logo `json:"from"`
}
