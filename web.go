package main

import (
	"embed"
	"fmt"
	"html/template"
	"myrouter/comm"
	. "myrouter/configs"
	"myrouter/funcs/wol"
	"myrouter/models/vn007plus"
	"net/http"
	"strings"
)

// IndexData 传递到首页模板的数据
type IndexData struct {
	IPv6 string
}

//go:embed "templates/*.html"
var templatesFS embed.FS

// 登录路由器使用的账号
var admin = vn007plus.Get(Conf.Admin.Username, Conf.Admin.Passwd)

// UseAuth 验证请求的中间件
// 中间件 https://www.alexedwards.net/blog/making-and-using-middleware
func UseAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 放行首页
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// 其它需要操作验证码
		auth := r.Header.Get("Authorization")
		if auth == "" || auth != Conf.Auth {
			fmt.Printf("错误的操作验证码：%s\n", auth)
			return
		}

		fmt.Printf("已验证操作码，继续下一步\n")
		next.ServeHTTP(w, r)
	}
}

// UseLogin 登录路由器的中间件
// 仅对 /api/ 中需要发送Auth的功能，才先模拟登录路由器
func UseLogin(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 不是 /api/ 的请求，直接下一步
		if strings.Index(r.URL.Path, "/api/") != 0 {
			next.ServeHTTP(w, r)
			return
		}

		// 登录路由器
		err := admin.Login()

		// 登录失败，无法继续下一步
		if err != nil {
			_, errW := w.Write([]byte(fmt.Sprintf("登录路由器失败：%s", err)))
			comm.Panic(errW)
			return
		}

		// 登录成功，继续下一步
		fmt.Printf("已登录路由器\n")
		next.ServeHTTP(w, r)
	}
}

// Index 首页
func Index(w http.ResponseWriter, _ *http.Request) {
	tpl, err := template.ParseFS(templatesFS, "templates/index.html")
	comm.Panic(err)
	err = tpl.Execute(w, "Hello.")
	comm.Panic(err)
}

// Reboot 重启路由器
func Reboot(w http.ResponseWriter, _ *http.Request) {
	msg := "正在重启路由器…"
	err := admin.Reboot()
	if err != nil {
		msg = fmt.Sprintf("重启路由器出错：%s", err)
	}

	_, err = w.Write([]byte(msg))
	comm.Panic(err)
}

// WakeupPC 唤醒网络设备
func WakeupPC(w http.ResponseWriter, _ *http.Request) {
	if Conf.WOL.MACAddr == "" {
		_, err := w.Write([]byte("无法网络唤醒电脑：目标 MAC 地址为空"))
		comm.Panic(err)
		return
	}

	msg := "正在网络唤醒电脑…"
	err := wol.Wakeup(Conf.WOL.MACAddr)
	if err != nil {
		msg = fmt.Sprintf("无法网络唤醒电脑：%s", err)
	}
	_, err = w.Write([]byte(msg))
	comm.Panic(err)
}
