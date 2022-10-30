package entities

// JResult 响应内容
type JResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// IPAddrs IP 地址信息
type IPAddrs struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}
