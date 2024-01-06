package buffer

import (
	"github.com/gammazero/deque"
	"sync"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"go.uber.org/atomic"

	"github.com/liuhailove/tc-base-go/mediatransportutil/pkg/bucket"
	"github.com/liuhailove/tc-base-go/mediatransportutil/pkg/nack"
	"github.com/liuhailove/tc-base-go/mediatransportutil/pkg/twcc"
	"github.com/liuhailove/tc-base-go/protocol/logger"
	"github.com/liuhailove/tc-server/pkg/sfu/audio"
	sutils "github.com/liuhailove/tc-server/pkg/utils"
)

// 同步信源(SSRC)标识符：占32位，用于标识同步信源。该标识符是随机选择的，参加同一视频会议的两个同步信源不能有相同的SSRC。
//
//·特约信源(CSRC)标识符：每个CSRC标识符占32位，可以有0～15个。每个CSRC标识了包含在该RTP报文有效载荷中的所有特约信源。
//
// 这里的同步信源是指产生媒体流的信源，例如麦克风、摄像机、RTP混合器等；它通过RTP报头中的一个32位数字SSRC标识符来标识，而不依赖于网络地址，接收者将根据SSRC标识符来区分不同的信源，进行RTP报文的分组。
//
// 特约信源是指当混合器接收到一个或多个同步信源的RTP报文后，经过混合处理产生一个新的组合RTP报文，并把混合器作为组合RTP报文的 SSRC，而将原来所有的SSRC都作为CSRC传送给接收者，使接收者知道组成组合报文的各个SSRC。

const (
	ReportDelta = time.Second
)

// pendingPacket 待处理的数据包
type pendingPacket struct {
	arrivalTime time.Time // 到达时间
	packet      []byte    // 包数据
}

// ExtPacket 扩展包
type ExtPacket struct {
	VideoLayer                                    // 视频编码层
	Arrival              time.Time                // 到达时间
	Packet               *rtp.Packet              // rtp包
	Payload              interface{}              // 负载
	KeyFrame             bool                     // 是否为关键帧
	RawPacket            []byte                   // 原始包
	DependencyDescriptor *ExtDependencyDescriptor // 依赖
}

// Buffer 包含所有数据包
type Buffer struct {
	sync.RWMutex
	bucket        *bucket.Bucket
	nacker        *nack.NackQueue
	videoPool     *sync.Pool
	audioPool     *sync.Pool
	codecType     webrtc.RTPCodecType
	extPackets    deque.Deque[*ExtPacket]
	pPackets      []pendingPacket
	closeOnce     sync.Once
	mediaSSRC     uint32
	clockRate     uint32
	lastReport    time.Time
	twccExt       uint8
	audioLevelExt uint8
	bound         bool
	closed        atomic.Bool
	mime          string

	// supported feedbacks
	latestTSForAudioLevelInitialized bool
	latestTSForAudioLevel            uint32

	twcc             *twcc.Responder
	audioLevelParams audio.AudioLevelParams
	audioLevel       *audio.AudioLevel

	lastPacketRead int

	pliThrottle int64

	rtpStats             *RTPStats
	rrSnapshotId         uint32
	deltaStatsSnapshotId uint32

	lastFractionLostToReport uint8 // 订阅者丢失的最后一部分，应向发布者报告；仅限音频

	// 回调
	onClose            func()
	onRtcpFeedback     func([]rtcp.Packet)
	onRtcpSenderReport func()
	onFpsChanged       func()
	onFinalRtpStats    func(*RTPStats)

	// logger
	logger logger.Logger

	// 依赖描述符
	ddExt    uint8
	ddParser *DependencyDescriptorParser

	paused              bool
	frameRateCalculator [DefaultMaxLayerSpatial + 1]FrameRateCalculator
	frameRateCalculated bool
}

// NewBuffer 构建一个新的Buffer
func NewBuffer(ssrc uint32, vp, ap *sync.Pool) *Buffer {
	// 将通过 SetLogger 重置正确的上下文
	l := logger.GetLogger()
	b := &Buffer{
		mediaSSRC:   ssrc,
		videoPool:   vp,
		audioPool:   ap,
		pliThrottle: int64(500 * time.Millisecond),
		logger:      l.WithComponent(sutils.ComponentPub).WithComponent(sutils.ComponentSFU),
	}
	b.extPackets.SetMinCapacity(7)
	return b
}

func (b *Buffer) SetLogger(logger logger.Logger) {
	b.Lock()
	defer b.Unlock()

	b.logger = logger.WithComponent(sutils.ComponentSFU)
	if b.rtpStats != nil {
		b.rtpStats.SetLogger(b.logger)
	}
}
