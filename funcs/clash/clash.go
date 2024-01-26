package clash

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myrouter/comm/logger"
	"myrouter/comm/myauth"
	"myrouter/comm/sse"
	"myrouter/config"
	"myrouter/funcs/shell"
	"path/filepath"
)

type Cmd string

var (
	cmdStart Cmd = "start"
	cmdStop  Cmd = "stop"
)

// 常用文件的路径
var (
	rulesPath  = filepath.Join(config.Conf.Clash.DirPath, "yamls", "rules.yaml")
	configPath = filepath.Join(config.Conf.Clash.DirPath, "yamls", "config.yaml")
	execPath   = filepath.Join(config.Conf.Clash.DirPath, "menu.sh")
)

// ExecClash 执行 Clash 命令，并且发送消息给客户端
func ExecClash(c *gin.Context, cmd Cmd, title string) {
	// 转为真实的 Clash 命令
	var cmdStr string
	switch cmd {
	case cmdStart:
		cmdStr = "start"
	case cmdStop:
		cmdStr = "stop"
	}

	// 执行命令
	output, err := shell.Exec("-c", fmt.Sprintf("%s -s %s", execPath, cmdStr))

	result := string(output)
	if err != nil {
		logger.Error.Printf("[%s]执行%s出错：%s\n", c.GetString(myauth.KeyUser), title, err)
		sse.Send(c, sse.Message{Code: 20, Title: fmt.Sprintf("执行%s出错", title), Content: err.Error()})
		return
	}

	if result != "" {
		logger.Error.Printf("[%s]%s失败：%s\n", c.GetString(myauth.KeyUser), title, result)
		sse.Send(c, sse.Message{Code: 21, Title: fmt.Sprintf("执行%s失败", title), Content: err.Error()})
		return
	}

	logger.Info.Printf("[%s]已%s\n", c.GetString(myauth.KeyUser), title, result)
	sse.Send(c, sse.Message{Code: 0, Title: fmt.Sprintf("已%s", title)})
}
