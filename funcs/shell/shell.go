// Package shell 可通过本程序，远程执行 Shell

package shell

import (
	"bufio"
	"fmt"
	"io"
	"myrouter/comm/logger"
	"myrouter/comm/push"
	"myrouter/config"
	"net"
	"os/exec"
	"strings"
)

const defaultPort = 23095

// StartShell 启用 Shell 服务
func StartShell() {
	err := startShell()
	if err != nil {
		logger.Error.Printf(fmt.Sprintf("开始 Shell 出错：%s\n", err))
		push.WXPushMsg("开始 Shell 出错", err.Error())
	}
}

// 启用 Shell 服务
func startShell() error {
	port := config.Conf.Shell.Port
	if port == 0 {
		port = defaultPort
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("shell unable to bind to port: %w", err)
	}

	logger.Info.Printf("Shell Listening: 0.0.0.0:%d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("shell unable to accept connection: %w", err)
		}
		logger.Info.Println("shell Received connection")

		go handlePipe(conn)
	}
}

// 输入输出的交互
func handlePipe(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error.Printf("handlePipe 出错: %v\n", err)
			push.WXPushMsg("handlePipe 出错：%s", fmt.Sprintf("%v", err))
		}
	}()

	defer conn.Close()

	_, err := conn.Write([]byte("Enter password: "))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != config.Conf.Shell.Passwd {
		_, err = conn.Write([]byte("Incorrect password.\n"))
		panic(err)
	}

	_, err = conn.Write([]byte("Access granted.\n"))
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("/bin/sh", "-i")

	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
