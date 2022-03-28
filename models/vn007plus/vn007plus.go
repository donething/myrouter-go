package vn007plus

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"io"
	"math/rand"
	. "myrouter/configs"
	"time"
)

// 路由器账号实体，每个账号密码对应一个实体
type account struct {
	username  string
	passwd    string
	userLevel string
	sessionId string
	client    dohttp.DoClient
}

const (
	// 发给给路由器的命令中，指定操作的类型
	mGet  = "GET"
	mPost = "POST"

	// 发给给路由器的命令中，指定操作
	cmdLogin         = 100
	cmdNextLoginTime = 232
	cmdReboot        = 6
)

var (
	headers = map[string]string{
		"Host": Conf.IP,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, " +
			"like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"X-Requested-With": "XMLHttpRequest",
	}
)

// Get 获取账号实体，每个账号密码对应一个路由器账号实体
func Get(username string, passwd string) *account {
	return &account{
		username: username,
		passwd:   passwd,
		client:   dohttp.New(10*time.Second, true, false),
	}
}

// Login 登录
// 在文件中 /js/login.js
func (a *account) Login() error {
	// 先需要从服务端获取登录 Token
	bs, err := a.client.PostJSONObj(PostURL, genPostBasic(cmdNextLoginTime, mGet, ""), headers)
	if err != nil {
		return fmt.Errorf("获取登录前的 Token 出错：%w\n", err)
	}

	// 解析，获取 Token
	var loginTimeResp LoginTimeResp
	err = json.Unmarshal(bs, &loginTimeResp)
	if err != nil {
		return fmt.Errorf("解析获取登录前的 Token 的响应出错：%w\n", err)
	}
	if loginTimeResp.Token == "" {
		return fmt.Errorf("获取登录前的 Token 为空：%s\n", string(bs))
	}

	// 加密密码，以供登录
	var pwdHash string
	h := sha256.New()
	_, err = io.WriteString(h, loginTimeResp.Token+a.passwd)
	if err != nil {
		return fmt.Errorf("加密登录密码出错：%w\n", err)
	}
	pwdHash = fmt.Sprintf("%x", h.Sum(nil))

	//
	// 执行登录
	var ss = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%f", rand.Float32())))) +
		fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%f", rand.Float32()))))
	var data = genPostBasic(cmdLogin, mPost, ss)
	var loginData = LoginData{
		PostBasic:     data,
		Username:      a.username,
		Passwd:        pwdHash,
		IsAutoUpgrade: "0",
	}
	bs, err = a.client.PostJSONObj(PostURL, loginData, headers)
	if err != nil {
		return fmt.Errorf("发送登录请求出错：%w\n", err)
	}

	// 解析登录结果
	var loginResp LoginResp
	err = json.Unmarshal(bs, &loginResp)
	if err != nil {
		return fmt.Errorf("解析登录响应出错：%w\n", err)
	}

	// 判断登录是否成功
	if loginResp.LoginFail == "fail" {
		bodyBS, _ := json.Marshal(loginData)
		return fmt.Errorf("登录路由器失败：'%s' ==> token: %s, data: '%s'\n",
			string(bs), loginTimeResp.Token, string(bodyBS))
	}

	// 登录成功
	a.sessionId = loginResp.SessionID
	return nil
}

// Reboot 重启路由器
func (a *account) Reboot() error {
	// 执行请求
	bs, err := a.client.PostJSONObj(PostURL, genPostBasic(cmdReboot, mPost, a.sessionId), headers)
	if err != nil {
		return fmt.Errorf("发送重启请求出错：%w\n", err)
	}

	// 解析
	var rebootResp RebootResp
	err = json.Unmarshal(bs, &rebootResp)
	// 成功执行重启时不会返回任何内容，所以排除响应长度为 0 的错误
	if err != nil && len(bs) != 0 {
		return fmt.Errorf("解析重启的响应出错：%w\n", err)
	}

	// 重启失败
	if rebootResp.Message == "NO_AUTH" {
		return fmt.Errorf("重启路由器失败：'%s' ==> SessionID: '%s'\n", string(bs), a.sessionId)
	}

	// 重启成功
	return nil
}
