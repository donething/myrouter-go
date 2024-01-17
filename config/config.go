package config

import (
	"flag"
	"github.com/donething/utils-go/doconf"
	"myrouter/comm/logger"
	"os"
)

const (
	// 配置的默认文件名
	name = "myrouter.json"
)

var (
	// 配置文件的路径
	confPath string

	// Conf 配置的实例
	Conf Config
)

func init() {
	flag.StringVar(&confPath, "c", name, "指定配置文件的路径")
	flag.Parse()

	// 读取配置
	exist, err := doconf.Init(confPath, &Conf)
	if err != nil {
		// 不能用 push.Panic() 会和 push 包导致"import cycle not allowed"
		panic(err)
	}

	if !exist {
		logger.Warn.Printf("已创建配置文件，请填写后，重新运行程序\n")
		os.Exit(0)
	}

	// 设置默认值
	if Conf.Router.Logo == "" {
		// 不使用常量 redmi.Logo，避免循环导入包
		Conf.Router.Logo = "RedmiAX6000"
	}
	if Conf.Clash.RulesPath == "" {
		Conf.Clash.RulesPath = "/data/clash/yamls/rules.yaml"
	}
}
