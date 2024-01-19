package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myrouter/comm/logger"
	"myrouter/comm/myauth"
	. "myrouter/config"
	"myrouter/funcs/status"
	"myrouter/funcs/wol"
	"myrouter/models"
	"myrouter/routers/redmi"
	"myrouter/routers/vn007p"
	"net/http"
)

// MyRouter 当前路由器的实例（用于需要登录路由器后执行的操作）
var MyRouter Router

func init() {
	// 根据配置文件，创建相应路由器的实例
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
func Reboot(c *gin.Context) {
	// 解析 JSON
	var data models.PostData[bool]
	err := c.BindJSON(&data)
	if err != nil {
		logger.Error.Printf("[%s]解析请求中的数据出错：%s\n", c.GetString(myauth.Key), err)
		c.JSON(http.StatusOK, models.Result{Code: 1500, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)})
		return
	}

	// 仅在为 true 时恢复到上次的规则
	if !data.Data {
		logger.Error.Printf("[%s]根据传递的参数'%t'，不重启路由器\n", c.GetString(myauth.Key), data.Data)
		c.JSON(http.StatusOK, models.Result{Code: 1510,
			Msg: fmt.Sprintf("根据传递的参数'%v'，不重启路由器：", data.Data)},
		)
		return
	}

	err = MyRouter.Reboot()
	if err != nil {
		logger.Error.Printf("[%s]重启路由器出错：%s\n", c.GetString(myauth.Key), err)
		c.JSON(http.StatusOK, models.Result{Code: 1520, Msg: fmt.Sprintf("重启路由器出错：%s", err)})
		return
	}

	logger.Info.Printf("[%s]正在重启路由器…\n", c.GetString(myauth.Key))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "正在重启路由器…"})
}

// WakeupPC 唤醒网络设备
//
// POST /api/router/reboot
//
// JSON 表单数据类型为 PostData[bool]
func WakeupPC(c *gin.Context) {
	if Conf.WOL.MACAddr == "" {
		logger.Warn.Printf("[%s]无法网络唤醒电脑：目标 MAC 地址为空\n", c.GetString(myauth.Key))
		c.JSON(http.StatusOK, models.Result{Code: 1600, Msg: "无法网络唤醒电脑：目标 MAC 地址为空"})
		return
	}

	// 解析 JSON
	var data models.PostData[bool]
	err := c.BindJSON(&data)
	if err != nil {
		logger.Error.Printf("[%s]解析请求中的数据出错：%s\n", c.GetString(myauth.Key), err)
		c.JSON(http.StatusOK, models.Result{Code: 1610, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)})
		return
	}

	// 仅在为 true 时恢复到上次的规则
	if !data.Data {
		logger.Error.Printf("[%s]根据传递的参数'%t'，不网络唤醒电脑\n", c.GetString(myauth.Key), data.Data)
		c.JSON(http.StatusOK, models.Result{Code: 1620,
			Msg: fmt.Sprintf("根据传递的参数'%v'，不网络唤醒电脑：", data.Data)},
		)
		return
	}

	err = wol.Wakeup(Conf.WOL.MACAddr)
	if err != nil {
		logger.Error.Printf("[%s]网络唤醒电脑出错：%s\n", c.GetString(myauth.Key), err)
		c.JSON(http.StatusOK, models.Result{Code: 1630, Msg: fmt.Sprintf("网络唤醒电脑出错：%s", err)})
		return
	}

	logger.Info.Printf("[%s]已发送唤醒电脑的网络包\n", c.GetString(myauth.Key))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已发送唤醒电脑的网络包"})
}

// RouterStatus 获取路由器的状态
//
// GET /api/router/ip
//
// 返回 Result[models.RouterStatus]
func RouterStatus(c *gin.Context) {
	s := status.GetRouterStatus()

	logger.Info.Printf("[%s]已获取路由器的状态\n", c.GetString(myauth.Key))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已获取路由器的状态", Data: s})
}
