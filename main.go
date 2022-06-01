package main

import (
	"embed"
	"fmt"
	"myrouter/comm"
	"myrouter/funcs/update_ip"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

//go:embed "templates/*.html"
var templatesFS embed.FS

//go:embed "static"
var staticFS embed.FS

func init() {
	whenInterrupt()

	update_ip.Update()
}

// 后台服务
// 使用 "0.0.0.0"可以同时监听 IPv4、IPv6
const addr = "0.0.0.0:20220"

func main() {
	fmt.Printf("开始服务：http://%s\n", addr)
	server := http.Server{
		Addr: addr,
	}

	// http.FS can be used to create a http Filesystem
	var sFileSystem = http.FS(staticFS)
	sfs := http.FileServer(sFileSystem)
	// Serve static files
	// 注意最后的"/"不能省略，否则会返回首页
	http.Handle("/static/", sfs)

	http.Handle("/", UseAuth(UseLogin(http.HandlerFunc(Index))))
	http.Handle("/api/reboot", UseAuth(UseLogin(http.HandlerFunc(Reboot))))
	http.Handle("/api/wol", UseAuth(UseLogin(http.HandlerFunc(WakeupPC))))

	comm.Panic(server.ListenAndServe())
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
// 强制结束本程序的命令：kill -9 $(pidof myrouter)
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
