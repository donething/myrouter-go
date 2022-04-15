package main

import (
	"fmt"
	"myrouter/comm"
	. "myrouter/configs"
	"myrouter/funcs/wol"
	"myrouter/models/vn007plus"
	"net/http"
	"strings"
	"time"
)

// 登录路由器使用的账号
var admin = vn007plus.Get(Conf.Admin.Username, Conf.Admin.Passwd)

// UseAuth 验证请求的中间件
// 中间件 https://www.alexedwards.net/blog/making-and-using-middleware
func UseAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 不限制跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 仅对 /api/ 路径进行操作验证，其它访问直接放行下一步
		if strings.Index(r.URL.Path, "/api/") != 0 {
			next.ServeHTTP(w, r)
			return
		}

		// 缺少验证参数、不正确，则抛弃此次请求
		auth := r.Header.Get("Authorization")
		if auth == "" || auth != Conf.Auth {
			abortAuth(w, r, "操作验证码有误，不通过验证")
			return
		}

		// 验证通过，继续下一步
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
	_, err := w.Write([]byte(fmt.Sprintf("Hello, %s", time.Now().String())))
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

// 抛弃请求
func abortAuth(w http.ResponseWriter, r *http.Request, msg string) {
	fmt.Printf("操作验证不通过 '%s'：%s\n", r.URL.Path, msg)
	_, err := w.Write([]byte(msg))
	comm.Panic(err)
}
