package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myrouter/comm/logger"
	"myrouter/comm/myauth"
	"myrouter/config"
	"myrouter/models"
	"myrouter/routers"
	"net/http"
	"strings"
)

const BearerSchema = "Bearer "

// UseAuth 验证请求的中间件
//
// 验证信息可以在请求头`Authorization`中，也可以在查询字符串`auth`中
//
// 中间件 https://www.alexedwards.net/blog/making-and-using-middleware
func UseAuth(c *gin.Context) {
	// 不是 /api/ 的请求，直接下一步
	if !strings.HasPrefix(c.FullPath(), "/api/") {
		c.Next()
		return
	}

	var token string
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		if len(authHeader) < len(BearerSchema) {
			abortAuth(c, fmt.Sprintf("非法的验证信息'%s'", authHeader))
			return
		}

		// Trim Bearer prefix to get the token
		token = authHeader[len(BearerSchema):]
	} else if authStr := c.Query("auth"); authStr != "" {
		token = authStr
	} else {
		abortAuth(c, "没有携带验证信息")
		return
	}

	logger.Info.Printf("访问'%s'的完整验证信息：'%s'\n", c.FullPath(), token)

	// 遍历已设置的授权码
	for name, auth := range config.Conf.Auths {
		if token == auth {
			// 验证通过，继续下一步
			logger.Info.Printf("[%s]已通过验证('%s')，继续下一步 '%s'\n", name, auth, c.FullPath())
			c.Set(myauth.KeyUser, name)
			return
		}
	}

	// 没有匹配到有效的验证码，禁止访问
	abortAuth(c, fmt.Sprintf("无效的验证信息'%s'", token))
}

// 拒绝访问
func abortAuth(c *gin.Context, msg string) {
	logger.Warn.Printf("'%s' 拒绝访问: %s\n", c.FullPath(), msg)

	c.AbortWithStatusJSON(http.StatusUnauthorized, models.Result{
		Code: 1000,
		Msg:  fmt.Sprintf("拒绝访问: %s", msg),
	})
}

// UseLoginRouter 登录路由器的中间件
func UseLoginRouter(c *gin.Context) {
	// 登录路由器
	err := routers.MyRouter.Login()

	// 登录失败，无法继续下一步
	if err != nil {
		logger.Info.Printf("[%s]登录路由器失败：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 1000, Msg: fmt.Sprintf("登录路由器失败：%s", err)})
		return
	}

	// 登录成功，继续下一步
	logger.Info.Printf("[%s]已登录路由器，继续下一步\n", c.GetString(myauth.KeyUser))
	c.Next()
}
