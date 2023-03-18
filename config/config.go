package config

import (
	"flag"
	"fmt"
	"github.com/donething/utils-go/doconf"
	"github.com/donething/utils-go/dolog"
	"myrouter/interfaces/jdc"
	"os"
	"strings"
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
	dolog.CkPanic(err)

	if !exist {
		fmt.Printf("已创建配置文件，请填写后，重新运行程序\n")
		os.Exit(0)
	}

	// 设置默认值
	if strings.TrimSpace(Conf.Router.From) == "" {
		Conf.Router.From = jdc.From
	}
}
