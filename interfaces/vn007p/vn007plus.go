package vn007p

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"myrouter/comm"
	"myrouter/config"
)

// Belongs 所属的路由器。用于根据路由器生成不同的对象
const Belongs = "Vn007"

// Vn007 路由器
type Vn007 struct {
	Username string
	Passwd   string

	UserLevel string
	// 自动设置
	SessionId string
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
	url = "http://%s/cgi-bin/http.cgi"

	headers = map[string]string{
		"Host": config.Conf.Gateway,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, " +
			"like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"X-Requested-With": "XMLHttpRequest",
	}
)

// Login 登录路由器
//
// 在文件中 /js/login.js
func (r *Vn007) Login() error {
	// 先需要从服务端获取登录 Token
	bs, err := comm.Client.PostJSONObj(url, genPostBasic(cmdNextLoginTime, mGet, ""), headers)
	if err != nil {
		return fmt.Errorf("获取登录前的 Token 出错：%w\n", err)
	}

	// 解析，获取 Token
	var loginTimeResp loginTimeResp
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
	_, err = io.WriteString(h, loginTimeResp.Token+r.Passwd)
	if err != nil {
		return fmt.Errorf("加密登录密码出错：%w\n", err)
	}
	pwdHash = fmt.Sprintf("%x", h.Sum(nil))

	//
	// 执行登录
	var ss = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%f", rand.Float32())))) +
		fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%f", rand.Float32()))))
	var data = genPostBasic(cmdLogin, mPost, ss)
	var loginData = loginData{
		postBasic:     data,
		Username:      r.Username,
		Passwd:        pwdHash,
		IsAutoUpgrade: "0",
	}

	bs, err = comm.Client.PostJSONObj(url, loginData, headers)
	if err != nil {
		return fmt.Errorf("发送登录请求出错：%w\n", err)
	}

	// 解析登录结果
	var loginResp loginResp
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
	r.SessionId = loginResp.SessionID
	return nil
}

// Reboot 重启路由器
func (r *Vn007) Reboot() error {
	// 执行请求
	bs, err := comm.Client.PostJSONObj(url, genPostBasic(cmdReboot, mPost, r.SessionId), headers)
	if err != nil {
		return fmt.Errorf("发送重启请求出错：%w\n", err)
	}

	// 解析
	var rebootResp rebootResp
	err = json.Unmarshal(bs, &rebootResp)
	// 成功执行重启时不会返回任何内容，所以排除响应长度为 0 的错误
	if err != nil && len(bs) != 0 {
		return fmt.Errorf("解析重启的响应出错：%w\n", err)
	}

	// 重启失败
	if rebootResp.Message == "NO_AUTH" {
		return fmt.Errorf("验证码有误 'NO_AUTH'")
	}

	// 重启成功
	return nil
}
