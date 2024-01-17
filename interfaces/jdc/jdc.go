package jdc

import "myrouter/models"

// From 所属的路由器。用于根据路由器生成不同的对象
const From models.Router = "JDC"

// JDC 京东云无线宝 路由器
type JDC struct {
	Username string
	Passwd   string
}

// Login 登录路由器
func (r *JDC) Login() error {
	return nil
}

// Reboot 重启路由器
func (r *JDC) Reboot() error {
	return nil
}
