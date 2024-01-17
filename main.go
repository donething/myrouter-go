package main

import (
	"embed"
	"fmt"
	"myrouter/comm/logger"
	"myrouter/comm/push"
	"myrouter/funcs/clash"
	"myrouter/funcs/shell"
	"myrouter/funcs/update"
	"myrouter/routers"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

// 内嵌资源都不要套入文件夹。如 "/html/templates"、"/html/static" 是错误的，会导致访问时 404 Not Found

//go:embed "templates/*.html"
var templatesFS embed.FS

//go:embed "static"
var staticFS embed.FS

func init() {
	whenInterrupt()

	go shell.StartShell()

	go update.Update()
}

// 后台服务
// 使用 "0.0.0.0"可以同时监听 IPv4、IPv6
const addr = "0.0.0.0:25816"

func main() {
	server := http.Server{
		Addr: addr,
	}

	// http.FS can be used to create a http Filesystem
	var sFileSystem = http.FS(staticFS)
	sfs := http.FileServer(sFileSystem)
	// Serve static files
	// 注意最后的"/"不能省略，否则会返回首页
	http.Handle("/static/", sfs)

	// 基本功能
	http.Handle("/", UseAuth(UseLogin(http.HandlerFunc(Index))))
	http.Handle("/api/reboot", UseAuth(UseLogin(http.HandlerFunc(routers.Reboot))))
	http.Handle("/api/wol", UseAuth(UseLogin(http.HandlerFunc(routers.WakeupPC))))

	// Clash 辅助工具
	http.Handle("/api/clash/rules/get", UseAuth(http.HandlerFunc(clash.GetRules)))
	http.Handle("/api/clash/rules/save", UseAuth(http.HandlerFunc(clash.SaveRules)))
	http.Handle("/api/clash/rules/backtolast", UseAuth(http.HandlerFunc(clash.BackToLastRules)))

	// 开始服务
	logger.Info.Printf("开始服务: //%s\n", addr)
	push.Panic(server.ListenAndServe())
}

// 中断处理程序
func whenInterrupt() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info.Println("\r- Ctrl+C pressed in Terminal")
		killServer()
		os.Exit(0)
	}()
}

// 不知什么原因，在 vn007+ 结束程序后，占用的端口不会释放，需要 kill -9 本程序的进程才能释放端口
// 强制结束本程序的命令：kill -9 $(pidof myrouter)
func killServer() {
	path, err := os.Executable()
	if err != nil {
		logger.Error.Printf("无法执行中断处理程序：获取执行文件的路径出错：%s\n", err)
		return
	}
	cmd := exec.Command("kill", "-9", fmt.Sprintf("$(pidof %s)", filepath.Base(path)))
	stdout, err := cmd.Output()

	if err != nil {
		logger.Error.Println("终结本程序出错：", err.Error())
		return
	}

	// Print the output
	logger.Warn.Println("已终结本程序\n", string(stdout))
}
