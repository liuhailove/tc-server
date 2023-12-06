package pacer

import (
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"time"
)

// ExtensionData 扩展数据
type ExtensionData struct {
	ID      uint8
	Payload []byte
}

// Packet 包
type Packet struct {
	Header             *rtp.Header
	Extensions         []ExtensionData
	Payload            []byte
	AbsSendTimeExtID   uint8
	TransportWideExtID uint8
	WriteStream        webrtc.TrackLocalWriter
	Metadata           interface{}
	OnSent             func(md interface{}, sentHeader *rtp.Header, payloadSize int, sentTime time.Time, sendError error)
}

// Pacer 控制发送端的数据发送速率
type Pacer interface {
	// Enqueue 包入队
	Enqueue(p Packet)
	// Stop 停止
	Stop()

	// SetInterval 设置间隔
	SetInterval(interval time.Duration)
	// SetBitrate 设置bit率
	SetBitrate(bitrate int)
}
