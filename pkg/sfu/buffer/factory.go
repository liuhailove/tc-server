package buffer

import "sync"

// Factory Buffer工厂
type Factory struct {
	sync.RWMutex
	videoPool   *sync.Pool // 视频Pool
	audioPool   *sync.Pool // 音频Pool
	rtpBuffers  map[uint32]*Buffer
	rtcpReaders map[uint32]*RTCPReader
}
