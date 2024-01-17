// Package update 更新 本路由器的 IP

package update

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"myrouter/comm/httpclient"
	"myrouter/comm/logger"
	"myrouter/comm/push"
	"myrouter/config"
	"myrouter/models"
	"net"
	"os/exec"
	"strings"
	"time"
)

// 临时保存上次获取的 IP 地址，以便与本次获取的相比较
var myIPAddrs *models.IPAddr

// Update 推送 IP 地址到远程服务器
//
// **重启服务**将触发立即推送 IP 地址到远程服务器
//
// 在获取出错后，暂停获取
func Update() {
	if config.Conf.Remote.UpdateIPAddr == "" {
		logger.Warn.Printf("无法推送 IP 地址：没有配置服务器地址\n")
		push.WXPushMsg("无法推送IP 地址", "没有配置服务器地址")
		return
	}

	// 执行程序后先推送一次 IP 地址
	err := up()
	if err != nil {
		logger.Error.Printf("推送 IP 地址出错：%s\n", err)
		push.WXPushMsg("推送 IP 地址出错", err.Error())
		return
	}

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			err = up()
			if err != nil {
				ticker.Stop()
				logger.Error.Printf("推送 IP 地址出错(暂停定时获取IP地址)：%s\n", err)
				push.WXPushMsg("推送 IP 地址出错", err.Error())
			}
		}
	}()
}

// 推送 IP 地址
func up() error {
	ip, err := GetLocalIPAddr()
	// ip, err := GetLocalIPAddrWithCmd()
	if err != nil {
		return err
	}

	// 都为空
	if ip.IPv4 == "" && ip.IPv6 == "" {
		return fmt.Errorf("获取到的 IPv4、IPv6 都为空")
	}

	ip.From = config.Conf.Router.Logo
	logger.Info.Printf("此次获取的 IP 地址：%+v\n", ip)
	if myIPAddrs == nil || ip.IPv4 != myIPAddrs.IPv4 || ip.IPv6 != myIPAddrs.IPv6 {
		myIPAddrs = ip
		logger.Info.Printf("IP 地址已改变，向远程发送新的地址\n")

		// 发送推送请求
		bs, err := httpclient.Client.PostJSONObj(config.Conf.Remote.UpdateIPAddr, *myIPAddrs, nil)
		if err != nil {
			return fmt.Errorf("发送推送 IP 的请求出错：%w", err)
		}

		// 分析结果
		var result models.Result
		err = json.Unmarshal(bs, &result)
		if err != nil {
			return fmt.Errorf("解析远程响应出错 '%s' ==> '%s'", err, string(bs))
		}
		if result.Code != 0 {
			return fmt.Errorf("推送 IP 失败：%s", result.Msg)
		}

		push.WXPushMsg("已推送 IP 地址", "已推送 IP 地址到远程服务器")
		logger.Info.Printf("已推送路由器 IP 地址到远程服务器\n")
		return nil
	}

	logger.Info.Printf("无需发送到远程：IP 地址信息没有变化\n")
	return nil
}

// GetLocalIPAddr 获取本地的 IP 地址
//
// 在运营商重新分配IP地址后，将无法获取到新的IP信息，需要用 GetLocalIPAddrWithCmd()
//
// @see https://www.cnblogs.com/hirampeng/p/11478995.html
func GetLocalIPAddr() (*models.IPAddr, error) {
	var ipAddrs = new(models.IPAddr)

	// 获取所有网卡
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	// 遍历
	for _, addr := range addrs {
		// 取网络地址的网卡的信息
		ipNet, isIpNet := addr.(*net.IPNet)
		// 是网卡并且不是本地环回网卡
		if isIpNet && !ipNet.IP.IsLoopback() {
			ipStr := strings.TrimSpace(ipNet.IP.String())
			// ipv4
			if len(strings.Split(ipStr, ".")) == 4 && dohttp.IsPublicIP(net.ParseIP(ipStr)) &&
				ipAddrs.IPv4 == "" {
				ipAddrs.IPv4 = ipStr
			}
			// ipv6
			if len(strings.Split(ipStr, ":")) == 8 && ipAddrs.IPv6 == "" {
				ipAddrs.IPv6 = ipStr
			}
		}
	}

	if ipAddrs.IPv4 == "" && ipAddrs.IPv6 == "" {
		return nil, fmt.Errorf("获取到的 IPV4、IPV6 地址都为空")
	}

	return ipAddrs, nil
}

// GetLocalIPAddrWithCmd 通过 Linux 命令获取地址信息
//
// 由于 GetLocalIPAddr() 在运营商重新分配IP地址后，将无法获取到新的IP信息，所以用Linux命令的方式获取
//
// @see https://superuser.com/a/1057290
//
// @see https://stackoverflow.com/a/41038684/8179418
func GetLocalIPAddrWithCmd() (*models.IPAddr, error) {
	var ipAddrs = new(models.IPAddr)
	outV6, errV6 := exec.Command("bash", "-c",
		"ip -6 addr | grep inet6 | awk -F '[ \\t]+|/' '{print $3}' | grep -v ^::1 | grep -v ^fe80",
	).Output()
	outV4, errV4 := exec.Command("bash", "-c",
		"ip -4 addr | grep inet | awk -F '[ \\t]+|/' '{print $3}' | grep -v ^127 | grep -v ^192",
	).Output()

	if errV6 != nil && errV4 != nil {
		return nil, fmt.Errorf("获取 IPV6、IPV4 地址都出错：IPv6(%s)，IPv4(%s)", errV6, errV4)
	}

	ipAddrs.IPv6 = string(outV6)
	ipAddrs.IPv4 = string(outV4)
	return ipAddrs, nil
}
