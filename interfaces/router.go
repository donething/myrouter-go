package interfaces

// Router 路由器接口
type Router interface {
	Login() error

	Reboot() error
}
