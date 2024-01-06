package sfu

import (
	"errors"
	"github.com/liuhailove/tc-base-go/protocol/logger"
	"github.com/liuhailove/tc-base-go/protocol/tc"
	"github.com/liuhailove/tc-server/pkg/sfu/buffer"
	"github.com/pion/webrtc/v3"
)

var (
	ErrReceiverClosed        = errors.New("receiver closed")
	ErrDownTrackAlreadyExist = errors.New("DownTrack already exist")
	ErrBufferNotFound        = errors.New("buffer not found")
)

type AudioLevelHandle func(level uint8, duration uint32)

type Bitrates [buffer.DefaultMaxLayerSpatial + 1][buffer.DefaultMaxLayerTemporal + 1]int64

// TrackReceiver 定义从远程对等方接收媒体的接口
type TrackReceiver interface {
	TrackID() tc.TrackID
	StreamID() string
	Codec() webrtc.RTPCodecParameters
	HeaderExtensions() []webrtc.RTPHeaderExtensionParameter
	IsClosed() bool

	ReadRTP(buf []byte, layer uint8, sn uint16) (int, error)
	GetLayeredBitrate() ([]int32, Bitrates)

	GetAudioLevel() (float64, bool)

	SendPLI(layer int32, force bool)

	SetUpTrackPaused(paused bool)
	SetMaxExpectedSpatialLayer(layer int32)

	AddDownTrack(track TrackSender) error
	DeleteDownTrack(participantID tc.ParticipantID)

	DebugInfo() map[string]interface{}

	TrackInfo() *tc.TrackInfo

	// GetPrimaryReceiverForRed 如果该接收器代表 RED 编解码器，则获取主接收器；否则它会自行返回
	GetPrimaryReceiverForRed() TrackReceiver

	//GetRedReceiver 获取主编解码器的Red接收器，由仅用于opus编解码器的前向Red编码使用
	GetRedReceiver() TrackReceiver

	// GetTemporalLayerFpsForSpatial 获取主编解码器的红色接收器，用于仅 opus 编解码器的前向红色编码	GetRedReceiver() TrackReceiver
	GetTemporalLayerFpsForSpatial(layer int32) []float32

	GetCalculatedClockRate(layer int32) uint32
	GetReferenceLayerRTPTimestamp(ts uint32, layer int32, referenceLayer int32) (uint32, error)
}

// WebRTCReceiver 接收媒体轨道
type WebRTCReceiver struct {
	logger logger.Logger
}
