package main

import (
	"fmt"
	"html/template"
	"myrouter/comm/logger"
	"myrouter/comm/push"
	"myrouter/config"
	"myrouter/models"
	"myrouter/routers"
	"net/http"
	"strings"
)

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
		err := routers.MyRouter.Login()

		// 登录失败，无法继续下一步
		if err != nil {
			result := models.Result{Code: 0, Msg: fmt.Sprintf("登录路由器失败：%s", err)}
			result.Response(w)
			return
		}

		// 登录成功，继续下一步
		logger.Info.Printf("已登录路由器，继续下一步\n")
		next.ServeHTTP(w, r)
	}
}

// Index 首页
func Index(w http.ResponseWriter, _ *http.Request) {
	tpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		logger.Error.Printf("解析首页模板出错：%s\n", err)
		push.WXPushMsg("解析首页模板出错", err.Error())
		return
	}

	err = tpl.Execute(w, "Hello.")
	if err != nil {
		logger.Error.Printf("执行首页模板出错：%s\n", err)
		push.WXPushMsg("执行首页模板出错", err.Error())
	}
}

// UseAuth 验证请求的中间件
// 中间件 https://www.alexedwards.net/blog/making-and-using-middleware
func UseAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 不是 /api/ 的请求，直接下一步
		if strings.Index(r.URL.Path, "/api/") != 0 {
			next.ServeHTTP(w, r)
			return
		}

		const BearerSchema = "Bearer "
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < len(BearerSchema) {
			abortAuth(w, r.URL.Path, "没有携带验证信息")
			return
		}

		logger.Info.Printf("访问'%s'的完整验证信息：'%s'\n", r.URL.Path, authHeader)

		// Trim Bearer prefix to get the token
		token := authHeader[len(BearerSchema):]
		if token == "" {
			abortAuth(w, r.URL.Path, fmt.Sprintf("非法的验证信息"))
			return
		}

		// 遍历已设置的授权码
		for key, auth := range config.Conf.Auths {
			if token == auth {
				// 验证通过，继续下一步
				logger.Info.Printf("'%s' 已通过验证'%s'('%s')，继续下一步\n", r.URL.Path, key, auth)
				next.ServeHTTP(w, r)
				return
			}
		}

		// 没有匹配到有效的验证码，禁止访问
		abortAuth(w, r.URL.Path, fmt.Sprintf("无效的验证信息"))
		return
	}
}

// 拒绝访问
func abortAuth(w http.ResponseWriter, path string, msg string) {
	logger.Warn.Printf("'%s' 拒绝访问: %s\n", path, msg)

	result := models.Result{Code: 1000, Msg: fmt.Sprintf("拒绝访问: %s", msg)}
	result.Response(w)
}
