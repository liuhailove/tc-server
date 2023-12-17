package buffer

import (
	"github.com/liuhailove/tc-base-go/mediatransportutil/pkg/bucket"
	"io"
	"sync"

	"github.com/pion/transport/v2/packetio"
)

type FactoryOfBufferFactory struct {
	videoPool *sync.Pool
	audioPool *sync.Pool
}

func NewFactoryOfBufferFactory(trackingPackets int) *FactoryOfBufferFactory {
	return &FactoryOfBufferFactory{
		videoPool: &sync.Pool{
			New: func() any {
				b := make([]byte, trackingPackets*bucket.MaxPktSize)
				return &b
			},
		},
		audioPool: &sync.Pool{
			New: func() any {
				b := make([]byte, bucket.MaxPktSize*200)
				return &b
			},
		},
	}
}

// Factory Buffer工厂
type Factory struct {
	sync.RWMutex
	videoPool   *sync.Pool // 视频Pool
	audioPool   *sync.Pool // 音频Pool
	rtpBuffers  map[uint32]*Buffer
	rtcpReaders map[uint32]*RTCPReader
}

func (f *Factory) GetOrNew(packetType packetio.BufferPacketType, ssrc uint32) io.ReadWriteCloser {
	f.Lock()
	defer f.Unlock()
	switch packetType {
	case packetio.RTCPBufferPacket:
		if reader, ok := f.rtcpReaders[ssrc]; ok {
			return reader
		}
		reader := NewRTCPReader(ssrc)
		f.rtcpReaders[ssrc] = reader
		reader.OnClose(func() {
			f.Lock()
			delete(f.rtcpReaders, ssrc)
			f.Unlock()
		})
		return reader
	case packetio.RTPBufferPacket:
		if reader, ok := f.rtpBuffers[ssrc]; ok {
			return reader
		}
		buffer := NewBuffer(ssrc, f.videoPool, f.audioPool)
		f.rtpBuffers[ssrc] = buffer
		buffer.OnClose(func() {
			f.Lock()
			delete(f.rtpBuffers, ssrc)
			f.Unlock()
		})
		return buffer
	}
	return nil
}

func (f *Factory) GetBufferPair(ssrc uint32) (*Buffer, *RTCPReader) {
	f.RLock()
	defer f.RUnlock()
	return f.rtpBuffers[ssrc], f.rtcpReaders[ssrc]
}

func (f *Factory) GetBuffer(ssrc uint32) *Buffer {
	f.RLock()
	defer f.RUnlock()
	return f.rtpBuffers[ssrc]
}

func (f *Factory) GetRTCPReader(ssrc uint32) *RTCPReader {
	f.RLock()
	defer f.RUnlock()
	return f.rtcpReaders[ssrc]
}
