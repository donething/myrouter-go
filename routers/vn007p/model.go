package vn007p

// BasicPost POST 请求需要发送的基础数据
type BasicPost struct {
	Cmd    int    `json:"cmd"`    // 执行的操作，对应 vn007plus 中"cmd"开头的属性
	Method string `json:"method"` // "POST"、"GET"，对应 vn007plus 中"m"开头的属性
	// 执行登录时为 Md5.md5(Math.random().toString()) + Md5.md5(Math.random().toString())
	// 登录成功后，设为响应中的 sessionId
	SessionID string `json:"sessionId"`
	Language  string `json:"language"` // 默认 "CN"
}

// BasicResp 响应的基础数据
type BasicResp struct {
	Success bool `json:"success"`
	Cmd     int  `json:"cmd"`
}

// 具体操作

// LoginTimeResp 实际登录前先获取 Token 的响应
type LoginTimeResp struct {
	*BasicResp
	Buffer        string `json:"buffer"` // "3"
	Token         string `json:"token"`  // 用于登录时加密密码
	NetxLoginTime string `json:"netx_login_time"`
}

// LoginData 登录时发送的数据
type LoginData struct {
	*BasicPost
	Username      string `json:"username"`
	Passwd        string `json:"passwd"` // 实际为 sha256_digest(LoginTimeResp.Token + psw)
	IsAutoUpgrade string `json:"isAutoUpgrade"`
}

// LoginResp 登录时的响应
type LoginResp struct {
	*BasicResp
	UserLevel string `json:"user_level"` // 登录用户的级别，3：管理员
	AUTH      string `json:"AUTH"`       // 登录后，值为"AUTH"
	SessionID string `json:"sessionId"`  // 登录以后执行其它操作时需要使用

	LoginFail string `json:"login_fail"` // 登录失败时值为"fail"，成功时没有该值（即为""），用于判断登录是否成功
}

// RebootData 重启时发送的数据
type RebootData struct {
	*BasicPost
	RebootType int `json:"rebootType"` // 一般设为 1
}

// RebootResp 重启的响应，仅在出错时可用；成功执行重启时不会返回任何内容
type RebootResp struct {
	*BasicResp
	Message string `json:"message"`
}

// 生成基础 POST body 的数据
func genPostBasic(cmd int, method string, session string) *BasicPost {
	return &BasicPost{
		Cmd:       cmd,
		Method:    method,
		SessionID: session,
		Language:  "CN",
	}
}
