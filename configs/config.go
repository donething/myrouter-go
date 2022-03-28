package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"os"
)

// 配置
type config struct {
	// 路由器 IP 地址，如"192.168.0.1"
	IP string `json:"ip"`

	// 管理员账号
	Admin struct {
		Username string `json:"username"`
		Passwd   string `json:"passwd"`
	} `json:"admin"`

	// 微信推送
	WXPush struct {
		Appid   string `json:"appid"`   // 组织 ID
		Secret  string `json:"secret"`  // 秘钥
		Agentid int    `json:"agentid"` // 应用（频道） ID
	} `json:"wx_push"`

	// 网络唤醒
	WOL struct {
		// 需要网络唤醒的 Mac 地址，如"89:0A:CD:EF:00:12"、"89:0a:cd:ef:00:12"或"01-23-45-56-67-89"
		MACAddr string `json:"mac_addr"`
	} `json:"wol"`
}

const (
	// 配置的默认文件名
	name = "myrouter.json"
)

var (
	// 配置文件的路径
	confPath string

	// Conf 配置的实例
	Conf config
	// PostURL 服务端 POST 地址
	PostURL = "http://%s/cgi-bin/http.cgi"
)

func init() {
	flag.StringVar(&confPath, "c", name, "指定配置文件的路径")
	flag.Parse()

	// 读取配置
	exist, err := dofile.Exists(confPath)
	Fatal(err)
	if !exist {
		bs, errMarshal := json.MarshalIndent(Conf, "", "  ")
		Fatal(errMarshal)
		_, errWrite := dofile.Write(bs, confPath, os.O_CREATE, 0600)
		Fatal(errWrite)

		fmt.Printf("已创建配置文件：'%s'，请填写配置后，重新运行\n", confPath)
		os.Exit(0)
	}

	fmt.Printf("读取配置文件：'%s'\n", confPath)
	bs, err := dofile.Read(confPath)
	Fatal(err)
	err = json.Unmarshal(bs, &Conf)
	Fatal(err)

	// 配置解析完成后，进行依赖配置的初始化
	PostURL = fmt.Sprintf(PostURL, Conf.IP)
}

// Fatal 出错时，强制关闭程序
func Fatal(err error) {
	if err != nil {
		fmt.Printf("运行出错：%s\n", err)
		os.Exit(1)
	}
}
