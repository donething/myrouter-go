package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "myrouter/configs"
	"myrouter/vn007plus"
	"net/http"
	"time"
)

// 登录路由器的账号
var admin = vn007plus.Get(Conf.Admin.Username, Conf.Admin.Passwd)

// Index 首页
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World! %s", time.Now().String())
}

// Reboot 重启
func Reboot(c *gin.Context) {
	err := admin.Reboot()
	if err != nil {
		c.JSON(http.StatusOK, JResult{
			Code: 3000,
			Msg:  fmt.Sprintf("重启路由器失败：%s", err),
			Data: nil,
		})
	}

	c.JSON(http.StatusOK, JResult{
		Code: 0,
		Msg:  "正在重启路由器",
		Data: nil,
	})
}

// UseLogin 登录路由器
func UseLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 访问首页直接通过
		if c.FullPath() == "/" {
			c.Next()
		}

		// 登录失败
		if err := admin.Login(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, JResult{
				Code: 2000,
				Msg:  fmt.Sprintf("登录路由器失败：%s", err),
				Data: nil,
			})
			return
		}

		// 登录成功，下一步
		c.Next()
	}
}
