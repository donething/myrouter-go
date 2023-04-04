package shell

import (
	"bufio"
	"fmt"
	"github.com/donething/utils-go/dolog"
	"io"
	"log"
	"myrouter/config"
	"net"
	"os/exec"
	"strings"
)

const defaultPort = 23040

// StartShell 启用 Shell 服务
func StartShell() {
	err := startShell()
	if err != nil {
		fmt.Println(err)
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

	log.Printf("shell Listening on 0.0.0.0:%d\n", port)

	for {
		conn, err := listener.Accept()
		log.Println("shell Received connection")
		if err != nil {
			return fmt.Errorf("shell unable to accept connection: %w", err)
		}

		go handlePipe(conn)
	}
}

// 输入输出的交互
func handlePipe(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("run time panic: %v", err)
		}
	}()

	defer conn.Close()

	_, err := conn.Write([]byte("Enter password: "))
	dolog.CkPanic(err)

	reader := bufio.NewReader(conn)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != config.Conf.Shell.Passwd {
		_, err = conn.Write([]byte("Incorrect password.\n"))
		dolog.CkPanic(err)
		return
	}

	_, err = conn.Write([]byte("Access granted.\n"))
	dolog.CkPanic(err)

	cmd := exec.Command("/bin/sh", "-i")

	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)

	err = cmd.Run()
	dolog.CkPanic(err)
}
