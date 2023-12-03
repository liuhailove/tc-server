package service

import (
	"time"

	"github.com/liuhailove/tc-server/pkg/rtc/types"
)

const (
	pingFrequency = 10 * time.Second // 每 10s ping一次
	pingTimeout   = 2 * time.Second  // ping 超时事件
)

// WSSignalConnection websocket信号连接
type WSSignalConnection struct {
	conn types.WebsocketClient
}
