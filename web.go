package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "myrouter/configs"
	"myrouter/entities"
	"myrouter/funcs/wol"
	"myrouter/vn007plus"
	"net/http"
	"time"
)

// 登录路由器的账号
var admin = vn007plus.Get(Conf.Admin.Username, Conf.Admin.Passwd)

// UseLogin 登录路由器
func UseLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 访问首页直接通过
		if c.FullPath() == "/" {
			c.Next()
		}

		// 登录失败
		if err := admin.Login(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, entities.JResult{
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

// Index 首页
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World! %s", time.Now().String())
}

// Reboot 重启
func Reboot(c *gin.Context) {
	err := admin.Reboot()
	if err != nil {
		c.JSON(http.StatusOK, entities.JResult{
			Code: 3000,
			Msg:  fmt.Sprintf("重启路由器失败：%s", err),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, entities.JResult{
		Code: 0,
		Msg:  "正在重启路由器",
		Data: nil,
	})
}

// WakeupPC 唤醒网络设备
func WakeupPC(c *gin.Context) {
	if Conf.WOL.MACAddr == "" {
		c.JSON(http.StatusOK, entities.JResult{
			Code: 3100,
			Msg:  "网络唤醒失败：目标 MAC 地址为空",
			Data: nil,
		})
		return
	}

	err := wol.Wakeup(Conf.WOL.MACAddr)
	if err != nil {
		c.JSON(http.StatusOK, entities.JResult{
			Code: 3110,
			Msg:  fmt.Sprintf("唤醒网络设备出错：%s", err),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, entities.JResult{
		Code: 0,
		Msg:  "成功唤醒目标网络设备",
		Data: nil,
	})
}
