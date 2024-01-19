package models

// RenderData 发送给客户端渲染网页的数据
type RenderData struct {
	Rules       []string `json:"rules"`       // 自定义的规则
	ProxyGroups []string `json:"proxyGroups"` // 所有代理组
}
