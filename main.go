package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myrouter/funcs/update_ip"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

const (
	// 服务的端口
	port = "9090"
)

func init() {
	whenInterrupt()

	update_ip.Update()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(UseLogin())

	router.GET("/", Index)
	router.GET("/api/status", Status)

	// 控制路由器
	router.POST("/api/reboot", Reboot)

	// 控制周边
	router.POST("/api/wol", WakeupPC)

	fmt.Printf("开始服务 :%s\n", port)
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		fmt.Printf("开启服务出错：%s\n", err)
		os.Exit(1)
	}
}

// 中断处理程序
func whenInterrupt() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		killServer()
		os.Exit(0)
	}()
}

// 不知什么原因，在 vn007+ 结束程序后，占用的端口不会释放，需要 kill -9 本程序的进程才能释放端口
func killServer() {
	// 仅在 Linux 下 kill 本程序
	if runtime.GOOS != "linux" {
		return
	}

	path, err := os.Executable()
	if err != nil {
		fmt.Printf("无法执行中断处理程序：获取执行文件的路径出错：%s\n", err)
		return
	}
	cmd := exec.Command("kill", "-9", fmt.Sprintf("$(pidof %s)", filepath.Base(path)))
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("终结本程序出错：", err.Error())
		return
	}

	// Print the output
	fmt.Println("已终结本程序\n", string(stdout))
}
