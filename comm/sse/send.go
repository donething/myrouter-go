package sse

import (
	"github.com/gin-gonic/gin"
	"myrouter/comm/logger"
	"myrouter/comm/myauth"
	"sync"
)

// 已连接的 SSE 客户端列表
var clients = map[string]*EventSender{}
var mu = sync.Mutex{}

// Send 向当前 SSE 客户端发送消息
//
// 客户端发起的请求携带`clientID`请求头，以便服务端返回 SSE 消息
func Send(c *gin.Context, message Message) {
	clientID := c.GetHeader("clientID")

	mu.Lock()
	sender, ok := clients[clientID]
	mu.Unlock()
	if !ok {
		logger.Warn.Printf("未知的 SSE clientID：'%s'\n", clientID)
		return
	}

	sender.send(message)
	logger.Info.Printf("[%s]已发送 SSE 消息'%s'\n", c.GetString(myauth.KeyUser), message.Title)
}

// LogClient 记录当前的 SSE 客户端
func LogClient(c *gin.Context, clientID string, sender *EventSender) {
	mu.Lock()
	defer mu.Unlock()

	clients[clientID] = sender
	logger.Info.Printf("[%s]已添加新 SSE 客户端'%s'。当前总共有 %d 个客户端\n",
		c.GetString(myauth.KeyUser), clientID, len(clients))
}

// DelClient 删除当前 SSE 客户端
func DelClient(c *gin.Context, clientID string) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, clientID)
	logger.Info.Printf("[%s]已删除 SSE 客户端'%s'。当前总共有 %d 个客户端\n",
		c.GetString(myauth.KeyUser), clientID, len(clients))
}
