package service

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/liuhailove/tc-base-go/protocol/logger"
	"github.com/liuhailove/tc-base-go/protocol/tc"
	"github.com/liuhailove/tc-server/pkg/rtc/types"
)

const (
	pingFrequency = 10 * time.Second // 每 10s ping一次
	pingTimeout   = 2 * time.Second  // ping 超时事件
)

// WSSignalConnection websocket信号连接
type WSSignalConnection struct {
	conn    types.WebsocketClient
	mu      sync.Mutex
	useJSON bool
}

// NewWSSignalConnection  新建WS信号连接
func NewWSSignalConnection(conn types.WebsocketClient) *WSSignalConnection {
	wsc := &WSSignalConnection{
		conn:    conn,
		mu:      sync.Mutex{},
		useJSON: false,
	}
	go wsc.pingWorker()
	return wsc
}

// ReadRequest 读取请求
func (c *WSSignalConnection) ReadRequest() (*tc.SignalRequest, int, error) {
	for {
		//处理特殊消息并传递其余消息
		messageType, payload, err := c.conn.ReadMessage()
		if err != nil {
			return nil, 0, err
		}

		msg := &tc.SignalRequest{}
		switch messageType {
		case websocket.BinaryMessage:
			if c.useJSON {
				c.mu.Lock()
				//如果客户端支持，则切换到protobuf
				c.useJSON = false
				c.mu.Unlock()
			}
			// protobuf 编码
			err := proto.Unmarshal(payload, msg)
			return msg, len(payload), err
		case websocket.TextMessage:
			c.mu.Lock()
			// json编码，也写回json
			c.useJSON = true
			c.mu.Unlock()
			err := protojson.Unmarshal(payload, msg)
			return msg, len(payload), err
		default:
			logger.Debugw("unsupported message", "messageType", messageType)
			return nil, len(payload), nil
		}
	}
}

// WriteResponse 写入响应
func (c *WSSignalConnection) WriteResponse(msg *tc.SignalResponse) (int, error) {
	var msgType int
	var payload []byte
	var err error

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.useJSON {
		msgType = websocket.TextMessage
		payload, err = protojson.Marshal(msg)
	} else {
		msgType = websocket.BinaryMessage
		payload, err = proto.Marshal(msg)
	}
	if err != nil {
		return 0, err
	}

	return len(payload), c.conn.WriteMessage(msgType, payload)
}

// pingWorker 定时ping
func (c *WSSignalConnection) pingWorker() {
	for {
		<-time.After(pingFrequency)
		err := c.conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(pingTimeout))
		if err != nil {
			logger.Errorw("ping worker error", err)
			return
		}
	}
}
