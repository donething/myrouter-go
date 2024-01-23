package clash

import (
	"fmt"
	"myrouter/config"
	"myrouter/funcs/shell"
	"path/filepath"
)

// 常用文件的路径
var (
	rulesPath  = filepath.Join(config.Conf.Clash.DirPath, "yamls", "rules.yaml")
	configPath = filepath.Join(config.Conf.Clash.DirPath, "yamls", "config.yaml")
	execPath   = filepath.Join(config.Conf.Clash.DirPath, "menu.sh")
)

// ExecRestartClash 执行重启 clash
func ExecRestartClash() ([]byte, error) {
	return shell.Exec("-c", fmt.Sprintf("%s -s start", execPath))
}
