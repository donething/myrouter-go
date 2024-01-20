package main

import (
	"embed"
	"fmt"
	"github.com/donething/utils-go/dolog"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"myrouter/comm/logger"
	"myrouter/comm/push"
	"myrouter/funcs/clash"
	"myrouter/funcs/shell"
	"myrouter/routers"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func init() {
	whenInterrupt()

	go shell.StartShell()

	// go status.Tick()
}

// 后台服务
// 使用 "0.0.0.0"可以同时监听 IPv4、IPv6
const port = "25816"

//go:embed assets/*
var embededFiles embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 使用 gzip 压缩，语句"gzip.Gzip(gzip.DefaultCompression)"不能放在Middleware()中，否则无效
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// 验证权限
	router.Use(UseAuth)

	// 内嵌静态资源，不必带着静态资源文件夹
	// 从嵌入的文件系统中获取assets文件夹
	assets, err := fs.Sub(embededFiles, "assets")
	dolog.CkPanic(err)
	// 将 assets 文件夹设置为 GIN 的静态文件系统
	router.StaticFS("/assets", http.FS(assets))
	// 解析模板文件
	tmpl, err := template.ParseFS(embededFiles, "assets/index.html")
	dolog.CkPanic(err)
	router.SetHTMLTemplate(tmpl)
	// 定义路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 路由器功能
	router.POST("/api//router/wol", routers.WakeupPC)
	router.POST("/api/router/reboot", routers.Reboot)
	router.GET("/api/router/status", routers.RouterStatus)

	// Clash 辅助工具
	router.GET("/api/clash/rules/all", clash.GetRules)
	router.POST("/api/clash/rule/add", clash.AddRule)
	router.POST("/api/clash/rule/del", clash.DelRule)
	router.POST("/api/clash/rule/override", clash.OverrideRules)
	router.POST("/api/clash/rule/backtolast", clash.BackToLastRules)
	router.GET("/api/clash/config/proxygroups", clash.GetProxyGroups)
	router.POST("/api/clash/manager/restart", clash.RestartClash)
	router.GET("/api/clash/data/render", clash.GetClashRenderData)

	// 开始服务
	logger.Info.Printf("开始服务：http://127.0.0.1:%s\n", port)
	push.Panic(http.ListenAndServe(":"+port, router))
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
