package shell

import (
	"os/exec"
)

// Exec 执行命令
//
// 注意第一个参数"-c": shell.Exec("-c", "/data/clash/clash.sh -s start")
//
// /bin/sh -c 后面的内容被设计为一条完整的命令，它会把所有的内容视为一个完整的单独命令
func Exec(args ...string) ([]byte, error) {
	// 创建Cmd
	cmd := exec.Command("/bin/sh", args...)

	// 返回获取执行结果
	return cmd.CombinedOutput()
}
