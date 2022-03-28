package update_ip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"io/ioutil"
	"myrouter/configs"
	"myrouter/entities"
	"myrouter/push"
	"net"
	"net/http"
	"strings"
	"time"
)

// 临时保存上次获取的 IP 地址信息，以便与本次获取的相比较
var myIPAddrs *entities.IPAddrs

// Update 更新 IP 地址到远程服务端
//
// **重启服务**将触发立即更新 IP 地址到远程服务端
func Update() {
	if configs.Conf.Remote.UpdateIPURL == "" {
		fmt.Printf("服务端域名没有配置，无需即时更新 IP 地址\n")
		return
	}

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			err := up()
			if err != nil {
				fmt.Printf("更新 IP 地址时出错：%s\n", err)
				push.WXPushCard("更新 IP 地址时出错", err.Error(), "", "")
				continue
			}
		}
	}()
}

// 更新 IP 地址
func up() error {
	ip, err := getLocalIPAddr()
	if err != nil {
		return err
	}

	fmt.Printf("此次获取的 IP 地址：%+v\n", ip)
	if myIPAddrs == nil || ip.IPv4 != myIPAddrs.IPv4 || ip.IPv6 != myIPAddrs.IPv6 {
		myIPAddrs = ip
		fmt.Printf("IP 地址已改变，向远程发送新的地址\n")
		// 发送更新请求
		ipBS, err := json.Marshal(ip)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", configs.Conf.Remote.UpdateIPURL, bytes.NewReader(ipBS))
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		// 解析响应
		bs, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
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

		fmt.Printf("远程服务器已更新 IP 地址信息\n")
	}

	fmt.Printf("IP 地址信息没有变化，无需发送到远程\n")
	return nil
}

// 获取本地的 IP 地址
//
// @see https://www.cnblogs.com/hirampeng/p/11478995.html
func getLocalIPAddr() (*entities.IPAddrs, error) {
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
			ipStr := ipNet.IP.String()
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

	return ipAddrs, nil
}
