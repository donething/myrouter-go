package redmi

import (
	"myrouter/config"
)

// Logo 所属的路由器。用于根据路由器生成不同的对象
const Logo config.Logo = "RedmiAX6000"

// Ax6000 京东云无线宝 路由器
type Ax6000 struct {
	Username string
	Passwd   string
}

// Login 登录路由器
func (r *Ax6000) Login() error {
	return nil
}

// Reboot 重启路由器
func (r *Ax6000) Reboot() error {
	return nil
}
