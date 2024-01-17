package routers

import (
	"fmt"
	"myrouter/comm/logger"
	. "myrouter/config"
	"myrouter/funcs/wol"
	"myrouter/models"
	"myrouter/routers/redmi"
	"myrouter/routers/vn007p"
	"net/http"
)

// MyRouter 登录路由器使用的账号
var MyRouter Router

func init() {
	switch Conf.Router.Logo {
	case redmi.Logo:
		MyRouter = &redmi.Ax6000{Username: Conf.Router.Username, Passwd: Conf.Router.Passwd}

	case vn007p.Logo:
		MyRouter = &vn007p.Vn007{Username: Conf.Router.Username, Passwd: Conf.Router.Passwd}

	default:
		panic(fmt.Errorf("未知的路由器标识'%s'", Conf.Router.Logo))
	}
}

// Reboot 重启路由器
func Reboot(w http.ResponseWriter, _ *http.Request) {
	var result models.Result
	err := MyRouter.Reboot()
	if err != nil {
		logger.Error.Printf("重启路由器出错：%s\n", err)
		result = models.Result{Code: 1500, Msg: fmt.Sprintf("重启路由器出错：%s", err)}
	} else {
		result = models.Result{Code: 0, Msg: "正在重启路由器…"}
	}

	result.Response(w)
}

// WakeupPC 唤醒网络设备
func WakeupPC(w http.ResponseWriter, _ *http.Request) {
	if Conf.WOL.MACAddr == "" {
		logger.Warn.Printf("无法网络唤醒电脑：目标 MAC 地址为空\n")
		result := models.Result{Code: 1600, Msg: "无法网络唤醒电脑：目标 MAC 地址为空"}
		result.Response(w)
		return
	}

	var result models.Result
	err := wol.Wakeup(Conf.WOL.MACAddr)
	if err != nil {
		logger.Error.Printf("网络唤醒电脑出错：%s\n", err)
		result = models.Result{Code: 1700, Msg: fmt.Sprintf("网络唤醒电脑出错：%s", err)}
	} else {
		result = models.Result{Code: 0, Msg: "正在网络唤醒电脑…"}
	}

	result.Response(w)
}
