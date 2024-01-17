package models

import (
	"encoding/json"
	"fmt"
	"myrouter/comm/logger"
	"myrouter/comm/push"
	"net/http"
)

// Result 响应内容
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Response 响应JResult的JSON到客户端
func (r *Result) Response(w http.ResponseWriter) {
	bs, err := json.Marshal(*r)
	if err != nil {
		return
	}

	// 写入客户端
	_, err = w.Write(bs)
	if err != nil {
		logger.Error.Printf("写入响应内容出错：'%s'。内容中的消息：'%s'\n", err, r.Msg)
		push.WXPushMsg("写入响应内容出错", fmt.Sprintf("内容中的消息：'%s'", r.Msg))
		return
	}
}
