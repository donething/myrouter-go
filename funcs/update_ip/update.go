package update_ip

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"myrouter/configs"
	"myrouter/entities"
	"myrouter/push"
	"net"
	"strings"
	"time"
)

// 临时保存上次获取的 IP 地址信息，以便与本次获取的相比较
var myIPAddrs *entities.IPAddrs

// Update 推送 IP 地址到远程服务端
//
// **重启服务**将触发立即推送 IP 地址到远程服务端
func Update() {
	if configs.Conf.Remote.UpdateIPURL == "" {
		fmt.Printf("服务端推送 IP 的地址没有配置，无法推送 IP 地址\n")
		push.WXPushCard("[路由器] 服务端地址没有配置", "无法推送 IP 地址", "", "")
		return
	}
	fmt.Printf("服务端推送 IP 的地址：'%s'\n", configs.Conf.Remote.UpdateIPURL)

	// 执行程序后先推送一次 IP 地址
	err := up()
	if err != nil {
		fmt.Printf("推送 IP 地址时出错：%s\n", err)
		push.WXPushCard("[路由器] 推送 IP 地址时出错", err.Error(), "", "")
		return
	}

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			err := up()
			if err != nil {
				fmt.Printf("推送 IP 地址时出错：%s\n", err)
				push.WXPushCard("[路由器] 推送 IP 地址时出错", err.Error(), "", "")
				continue
			}
		}
	}()
}

// 推送 IP 地址
func up() error {
	ip, err := GetLocalIPAddr()
	if err != nil {
		return err
	}

	fmt.Printf("此次获取的 IP 地址：%+v\n", ip)
	if myIPAddrs == nil || ip.IPv4 != myIPAddrs.IPv4 || ip.IPv6 != myIPAddrs.IPv6 {
		myIPAddrs = ip
		fmt.Printf("IP 地址已改变，向远程发送新的地址\n")
		// 发送推送请求
		var client = dohttp.New(30*time.Second, false, false)
		bs, err := client.PostJSONObj(configs.Conf.Remote.UpdateIPURL, *myIPAddrs, nil)
		if err != nil {
			return err
		}

		// 分析结果
		var result entities.JResult
		err = json.Unmarshal(bs, &result)
		if err != nil {
			return fmt.Errorf("解析远程响应出错 '%s' ==> '%s'", err, string(bs))
		}
		if result.Code != 0 {
			return fmt.Errorf("%s：%s", result.Msg, result.Data)
		}

		push.WXPushCard("[路由器] 已推送 IP 地址", "已推送路由器 IP 地址到远程服务器", "", "")
		fmt.Printf("已推送路由器 IP 地址到远程服务器\n")
		return nil
	}

	fmt.Printf("IP 地址信息没有变化，无需发送到远程\n")
	return nil
}

// GetLocalIPAddr 获取本地的 IP 地址
//
// @see https://www.cnblogs.com/hirampeng/p/11478995.html
func GetLocalIPAddr() (*entities.IPAddrs, error) {
	var ipAddrs = new(entities.IPAddrs)

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
