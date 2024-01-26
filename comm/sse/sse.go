package sse

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"myrouter/comm/logger"
	"net/http"
)

// Message 将要发送给客户端的消息
type Message struct {
	Code    int    `json:"code"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// EventSender 结构中包含了一个事件通道，用于将消息发送给客户端
type EventSender struct {
	eventCh chan Message
}

// NewEventSender 用于初始化一个事件发送器
func NewEventSender() *EventSender {
	return &EventSender{
		eventCh: make(chan Message, 10),
	}
}

// Send 向客户端发送消息
func (s *EventSender) send(event Message) {
	s.eventCh <- event
}

// start 初始化。将从 eventCh 读取消息并发送到客户端
//
// 如果发送消息到客户端出错，说明客户端可能已关闭，需要退出并删除客户端的信息
func (s *EventSender) start(c *gin.Context, w gin.ResponseWriter, clientID string) {
	defer func() {
		DelClient(c, clientID)
	}()

	for {
		event := <-s.eventCh
		data, err := json.Marshal(event)
		if err != nil {
			logger.Error.Printf("序列化数据出错：%s：\n", err)
			_, err = w.Write([]byte(fmt.Sprintf("序列化数据出错：%s\n\n", err)))
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("data: %s\n\n", string(data))))
		if err != nil {
			logger.Error.Printf("向客户端发送消息出错，将关闭该通道：%s。数据：%s\n", err, string(data))
			return
		}

		w.Flush()
	}
}

// UseSSEvents 向客户端发送消息
//
// GET /api/sse/events?auth=<authorization>&clientID=<uuid>
//
// 客户端发起 SSE 后，服务端将根据`clientID`记录当前客户端的 SSE 连接
func UseSSEvents(c *gin.Context) {
	clientID := c.Query("clientID")
	if len(clientID) != 36 {
		logger.Warn.Printf("clientID(UUID)不合法：'%s'\n", clientID)
		return
	}

	w := c.Writer
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	var eventSender = NewEventSender()
	LogClient(c, clientID, eventSender)

	eventSender.start(c, w, clientID)
}
