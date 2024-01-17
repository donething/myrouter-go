// Package wol 网络唤醒

package wol

import (
	"fmt"
	"myrouter/comm/logger"
	"myrouter/config"
	"net"
)

// Go发送UDP广播消息_蜗牛撵大象的博客-CSDN博客_golang udp广播
// https://blog.csdn.net/weixin_42651014/article/details/95636511

// 网络唤醒端口，默认为 9 或 7
const port = "9"

// Wakeup 唤醒指定 Mac 地址的设置
func Wakeup(macAddr string) error {
	// The address to broadcast to is usually the default `255.255.255.255` but
	// can be overloaded by specifying an override in the CLI arguments.
	// 目标地址，默认广播到当前网段中的所有地址
	udpAddrStr := fmt.Sprintf("%s:%s", "255.255.255.255", port)
	udpAddr, err := net.ResolveUDPAddr("udp", udpAddrStr)
	if err != nil {
		return err
	}

	// 本地地址。在路由器上网络唤醒电脑时需要此项，否则发送 UDP 广播无效，其它平台可为 nil
	localAddrStr := fmt.Sprintf("%s:%s", config.Conf.Gateway, port)
	localAddr, err := net.ResolveUDPAddr("udp", localAddrStr)
	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := New(macAddr)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	logger.Info.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	logger.Info.Printf("... Broadcasting to: %s\n", udpAddrStr)
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	logger.Info.Printf("Magic packet sent successfully to %s\n", macAddr)
	return nil
}
