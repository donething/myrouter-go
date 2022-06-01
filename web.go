package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"math"
	"myrouter/comm"
	. "myrouter/configs"
	"myrouter/funcs/wol"
	"myrouter/models/vn007plus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 登录路由器使用的账号
var admin = vn007plus.Get(Conf.Admin.Username, Conf.Admin.Passwd)

// UseAuth 验证请求的中间件
// 中间件 https://www.alexedwards.net/blog/making-and-using-middleware
func UseAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 放行内网访问
		if strings.Index(r.RemoteAddr, "127.0.0.1:") == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// 放行图标等资源文件
		if r.URL.Path == "/favicon.ico" || r.URL.Path == "/static/" {
			next.ServeHTTP(w, r)
			return
		}

		// 提取验证参数
		t := r.URL.Query().Get("t")
		s := r.URL.Query().Get("s")

		// 缺少验证参数，直接不通过验证，抛弃此次请求
		if strings.TrimSpace(t) == "" || strings.TrimSpace(s) == "" {
			abortAuth(w, r, "验证信息为空")
			return
		}

		// 验证时间戳，如果和系统时间误差多于 30*1000毫秒(30秒)，即验证不通过
		reqUnixMilli, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			abortAuth(w, r, "时间戳无法转换为数字")
			return
		}
		if math.Abs(float64(time.Now().UnixMilli()-reqUnixMilli)) > 30*1000 {
			abortAuth(w, r, "时间戳已过期")
			return
		}

		// 根据时间戳 t(毫秒)，计算 sha256。计算目标为 (操作验证码 + t + 操作验证码)
		sum := sha256.Sum256([]byte(Conf.Auth + t + Conf.Auth))
		sumStr := fmt.Sprintf("%x", sum)
		// 验证不通过，抛弃此次请求
		if strings.ToLower(sumStr) != strings.ToLower(s) {
			abortAuth(w, r, "验证失败")
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

// 抛弃请求
func abortAuth(w http.ResponseWriter, r *http.Request, msg string) {
	fmt.Printf("操作验证不通过 '%s'：%s\n", r.URL.Path, msg)
	_, err := w.Write([]byte(msg))
	comm.Panic(err)
}
