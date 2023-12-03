package types

import "time"

// WebsocketClient websocket客户端
type WebsocketClient interface {
	// ReadMessage 读取消息
	ReadMessage() (messageType int, p []byte, err error)
	// WriterMessage 写入消息
	WriterMessage(messageType int, data []byte) error
	// WriteControl 写入控制消息
	WriteControl(messageType int, data []byte, deadline time.Time) error
}
