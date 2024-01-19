package models

import (
	"myrouter/config"
)

// IPAddr IP 地址信息
type IPAddr struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

// Mem 内存使用情况
type Mem struct {
	Total       string  `json:"total"`
	Available   string  `json:"available"`
	Used        string  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// CPU CPU使用情况
type CPU struct {
	// 总使用率，不是每个核心的
	UsedPercent string `json:"usedPercent"`
}

// RouterStatus 路由器的状态
type RouterStatus struct {
	IPAddr IPAddr `json:"ipAddr"`
	Mem    Mem    `json:"mem"`
	CPU    CPU    `json:"cpu"`

	// 所属的路由器
	Logo config.Logo `json:"logo"`
}
