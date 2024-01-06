package sfu

import (
	"errors"
	"fmt"
	"github.com/liuhailove/tc-server/pkg/sfu/pacer"
	"github.com/pion/rtcp"
	"go.uber.org/atomic"
	"sync"
	"time"

	"github.com/liuhailove/tc-base-go/protocol/logger"
	"github.com/liuhailove/tc-base-go/protocol/tc"
	"github.com/liuhailove/tc-server/pkg/sfu/buffer"
	"github.com/pion/webrtc/v3"
)

// TrackSender 定义一个将媒体发送到远程对等点的接口
type TrackSender interface {
	// UpTrackLayerChange 上轨道图层更改
	UpTrackLayerChange()
	// UpTrackBitrateAvailabilityChange 音轨比特率可用性更改
	UpTrackBitrateAvailabilityChange()
	// UpTrackMaxPublishedLayerChange 上行音轨已发布的最大图层更改
	UpTrackMaxPublishedLayerChange(maxPublishedLayer int32)
	// UpTrackMaxTemporalLayerSeenChange 上行音轨最大时间层改变
	UpTrackMaxTemporalLayerSeenChange(maxTemporalLayerSeen int32)
	// UpTrackBitrateReport 上行音轨的比特率报告
	UpTrackBitrateReport(availableLayers []int32, bitrates Bitrates)
	// WriteRTP 写入RTP包
	WriteRTP(p *buffer.ExtPacket, layer int32) error
	// Close 关闭发送
	Close()
	// IsClosed 音轨发送是否已经关闭
	IsClosed() bool
	// ID 是该 Track 的全局唯一标识符
	ID() string
	// SubscriberID 订阅者ID
	SubscriberID() tc.ParticipantID
	TrackInfoAvailable()
	// HandleRTCPSenderReportData 处理RTCP发送报告数据
	HandleRTCPSenderReportData(payloadType webrtc.PayloadType, layer int32, srData *buffer.RTCPSenderReportData) error
}

// -------------------------------------------------------------------

const (
	RTPPaddingMaxPayloadSize      = 255
	RTPPaddingEstimatedHeaderSize = 20
	RTPBlankFramesMuteSeconds     = float32(1.0)
	RTPBlankFramesCloseSeconds    = float32(0.2)

	FlagStopRTXOnPLI = true

	KeyFrameIntervalMin = 200
	KeyFrameIntervalMax = 1000
	flushTimeout        = 1 * time.Second

	maxPadding = 2000

	waitBeforeSendPaddingOnMute = 100 * time.Millisecond
	maxPaddingOnMuteDuration    = 5 * time.Second
)

// -------------------------------------------------------------------

var (
	ErrUnknownKind                       = errors.New("unknown kind of codec")
	ErrOutOfOrderSequenceNumberCacheMiss = errors.New("out-of-order sequence number not found in cache")
	ErrPaddingOnlyPacket                 = errors.New("padding only packet that need not be forwarded")
	ErrDuplicatePacket                   = errors.New("duplicate packet")
	ErrPaddingNotOnFrameBoundary         = errors.New("padding cannot send on non-frame boundary")
	ErrDownTrackAlreadyBound             = errors.New("already bound")
)

var (
	VP8KeyFrame8x8 = []byte{
		0x10, 0x02, 0x00, 0x9d, 0x01, 0x2a, 0x08, 0x00,
		0x08, 0x00, 0x00, 0x47, 0x08, 0x85, 0x85, 0x88,
		0x85, 0x84, 0x88, 0x02, 0x02, 0x00, 0x0c, 0x0d,
		0x60, 0x00, 0xfe, 0xff, 0xab, 0x50, 0x80,
	}

	H264KeyFrame2x2SPS = []byte{
		0x67, 0x42, 0xc0, 0x1f, 0x0f, 0xd9, 0x1f, 0x88,
		0x88, 0x84, 0x00, 0x00, 0x03, 0x00, 0x04, 0x00,
		0x00, 0x03, 0x00, 0xc8, 0x3c, 0x60, 0xc9, 0x20,
	}
	H264KeyFrame2x2PPS = []byte{
		0x68, 0x87, 0xcb, 0x83, 0xcb, 0x20,
	}
	H264KeyFrame2x2IDR = []byte{
		0x65, 0x88, 0x84, 0x0a, 0xf2, 0x62, 0x80, 0x00,
		0xa7, 0xbe,
	}
	H264KeyFrame2x2 = [][]byte{H264KeyFrame2x2SPS, H264KeyFrame2x2PPS, H264KeyFrame2x2IDR}

	OpusSilenceFrame = []byte{
		0xf8, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
)

// -------------------------------------------------------------------

type DownTrackState struct {
	RTPStats                       *buffer.RTPStats
	DeltaStatsSnapshotId           uint32
	DeltaStatsOverriddenSnapshotId uint32
	ForwarderState                 ForwarderState
}

func (d DownTrackState) String() string {
	return fmt.Sprintf("DownTrackState{rtpStats: %s, delta: %d, deltaOverridden: %d, forwarder: %s}",
		d.RTPStats.ToString(), d.DeltaStatsSnapshotId, d.DeltaStatsOverriddenSnapshotId, d.ForwarderState.String())
}

// -------------------------------------------------------------------

type NackInfo struct {
	Timestamp      uint32
	SequenceNumber uint16
	Attempts       uint8
}

// DownTrackStreamAllocatorListener 下行流分配器监听
type DownTrackStreamAllocatorListener interface {
	// RTCP received
	OnREMB(dt *DownTrack, remb *rtcp.ReceiverEstimatedMaximumBitrate)
	OnTransportCCFeedback(dt *DownTrack, cc *rtcp.TransportLayerCC)

	// OnAvailableLayersChanged 视频层可用性已更改
	OnAvailableLayersChanged(dt *DownTrack)

	// OnBitrateAvailabilityChanged 视频层比特率可用性已更改
	OnBitrateAvailabilityChanged(dt *DownTrack)

	// OnMaxPublishedSpatialChanged 最大已发布空间层已更改
	OnMaxPublishedSpatialChanged(dt *DownTrack)

	// OnMaxPublishedTemporalChanged 已更改最大发布时间层
	OnMaxPublishedTemporalChanged(dt *DownTrack)

	// OnSubscriptionChanged 订阅已更改-静音/取消静音
	OnSubscriptionChanged(dt *DownTrack)

	// OnSubscribedLayerChanged 订阅的最大视频层已更改
	OnSubscribedLayerChanged(dt *DownTrack, layers buffer.VideoLayer)

	// OnResume 流已恢复
	OnResume(dt *DownTrack)

	// OnPacketSent 发送包事件
	OnPacketSent(dt *DownTrack, size int)

	// OnNACK 收到NACK
	OnNACK(dt *DownTrack, nackInfos []NackInfo)

	// OnRTCPReceiverReport 收到RTCP接收器报告
	OnRTCPReceiverReport(dt *DownTrack, rr rtcp.ReceiverReport)
}

type ReceiverReportListener func(dt *DownTrack, report *rtcp.ReceiverReport)

type DownTrackParams struct {
	Codec             []webrtc.RTPCodecParameters
	Receiver          TrackReceiver
	BufferFactory     *buffer.Factory
	SubID             tc.ParticipantID
	StreamID          string
	MaxTrack          int
	PlayoutDelayLimit *tc.PlayoutDelay
	Pacer             pacer.Pacer
	Logger            logger.Logger
	Trailer           []byte
}

// DownTrack 实现TrackLocal，是用于写入数据包的轨道
// 对于SFU用户，轨道处理数据包进行简单的联播
// 和SVC发布器。
// DownTrack具有以下生命周期
// - new
// - bound / unbound
// - closed
// 一旦关闭，DownTrack就不能再使用。
type DownTrack struct {
	params        DownTrackParams
	logger        logger.Logger
	id            tc.TrackID
	subscriberID  tc.ParticipantID
	kind          webrtc.RTPCodecType
	mime          string
	ssrc          uint32
	streamID      string
	maxTrack      int
	payloadType   uint8
	sequencer     *sequencer
	bufferFactory *buffer.Factory

	forwarder *Forwarder

	upstreamCodecs            []webrtc.RTPCodecParameters
	codec                     webrtc.RTPCodecCapability
	absSendTimeExtID          int
	transportWideExtID        int
	dependencyDescriptorExtID int
	playoutDelayExtID         int
	transceiver               atomic.Value
	writeStream               webrtc.TrackLocalWriter
	rtcpReader                *buffer.RTCPReader

	listenerLock            sync.Mutex
	receiverReportListeners []ReceiverReportListener

	bindLock sync.Mutex
	bound    atomic.Bool

	onBinding func(error)

	isClosed             atomic.Bool
	connected            atomic.Bool
	bindAndConnectedOnce atomic.Bool

	rtpStats *buffer.RTPStats

	totalRepeatedNACKs atomic.Uint32

	keyFrameRequestGeneration atomic.Uint32

	blankFramesGeneration atomic.Uint32

	connectionStats                *connectionquality.ConnectionStats
	deltaStatsSnapshotId           uint32
	deltaStatsOverriddenSnapshotId uint32

	isNACKThrottled atomic.Bool

	streamAllocatorLock             sync.RWMutex
	streamAllocatorListener         DownTrackStreamAllocatorListener
	streamAllocatorReportGeneration int
	streamAllocatorBytesCount       atomic.Uint32
	bytesSent                       atomic.Uint32
	bytesRetransmitted              atomic.Uint32

	playoutDelayBytes atomic.Value // 编组播放延迟的字节数
	playoutDelayAcked atomic.Bool

	pacer pacer.Pacer

	maxLayerNotifierCh chan struct{}

	cbMu                        sync.RWMutex
	onStatsUpdate               func(dt *DownTrack, stats *tc.AnalyticsStat)
	onMaxSubscribedLayerChanged func(dt *DownTrack, layer int32)
	onRttUpdate                 func(dt *DownTrack, rtt uint32)
	onCloseHandler              func(willBeResumed bool)
}
