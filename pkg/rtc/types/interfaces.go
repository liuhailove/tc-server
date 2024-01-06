package types

import (
	"fmt"
	"github.com/liuhailove/tc-base-go/protocol/auth"
	"github.com/pion/rtcp"
	"time"

	"github.com/pion/webrtc/v3"

	"github.com/liuhailove/tc-base-go/protocol/logger"
	"github.com/liuhailove/tc-base-go/protocol/tc"
	"github.com/liuhailove/tc-base-go/protocol/utils"
	"github.com/liuhailove/tc-server/pkg/routing"
	"github.com/liuhailove/tc-server/pkg/sfu/buffer"
	"github.com/liuhailove/tc-server/pkg/sfu/pacer"
)

// WebsocketClient websocket客户端
type WebsocketClient interface {
	// ReadMessage 读取消息
	ReadMessage() (messageType int, p []byte, err error)
	// WriteMessage 写入消息
	WriteMessage(messageType int, data []byte) error
	// WriteControl 写入控制消息
	WriteControl(messageType int, data []byte, deadline time.Time) error
}

// AddSubscriberParams 添加订阅者参数
type AddSubscriberParams struct {
	AllTracks bool
	TrackIDs  []tc.TrackID
}

// ---------------------------------------------

type MigrateState int32 // 迁移状态

const (
	MigrateStateInit MigrateState = iota
	MigrateStateSync
	MigrateStateComplete
)

func (m MigrateState) String() string {
	switch m {
	case MigrateStateInit:
		return "MIGRATE_STATE_INIT"
	case MigrateStateSync:
		return "MIGRATE_STATE_SYNC"
	case MigrateStateComplete:
		return "MIGRATE_STATE_COMPLETE"
	default:
		return fmt.Sprintf("%d", int(m))
	}
}

// ---------------------------------------------

type SubscribedCodecQuality struct {
	CodecMime string
	Quality   tc.VideoQuality
}

// ---------------------------------------------

type ParticipantCloseReason int // 参与者关闭原因

const (
	ParticipantCloseReasonClientRequestLeave ParticipantCloseReason = iota
	ParticipantCloseReasonRoomManagerStop
	ParticipantCloseReasonRoomClose
	ParticipantCloseReasonVerifyFailed
	ParticipantCloseReasonJoinFailed
	ParticipantCloseReasonJoinTimeout
	ParticipantCloseReasonStateDisconnected
	ParticipantCloseReasonPeerConnectionDisconnected
	ParticipantCloseReasonDuplicateIdentity
	ParticipantCloseReasonMigrationComplete
	ParticipantCloseReasonStale
	ParticipantCloseReasonServiceRequestRemoveParticipant
	ParticipantCloseReasonServiceRequestDeleteRoom
	ParticipantCloseReasonSimulateMigration
	ParticipantCloseReasonSimulateNodeFailure
	ParticipantCloseReasonSimulateServerLeave
	ParticipantCloseReasonNegotiateFailed
	ParticipantCloseReasonMigrationRequested
	ParticipantCloseReasonOvercommitted
	ParticipantCloseReasonPublicationError
	ParticipantCloseReasonSubscriptionError
)

func (p ParticipantCloseReason) String() string {
	switch p {
	case ParticipantCloseReasonClientRequestLeave:
		return "CLIENT_REQUEST_LEAVE"
	case ParticipantCloseReasonRoomManagerStop:
		return "ROOM_MANAGER_STOP"
	case ParticipantCloseReasonRoomClose:
		return "ROOM_CLOSE"
	case ParticipantCloseReasonVerifyFailed:
		return "VERIFY_FAILED"
	case ParticipantCloseReasonJoinFailed:
		return "JOIN_FAILED"
	case ParticipantCloseReasonJoinTimeout:
		return "JOIN_TIMEOUT"
	case ParticipantCloseReasonStateDisconnected:
		return "STATE_DISCONNECTED"
	case ParticipantCloseReasonPeerConnectionDisconnected:
		return "PEER_CONNECTION_DISCONNECTED"
	case ParticipantCloseReasonDuplicateIdentity:
		return "DUPLICATE_IDENTITY"
	case ParticipantCloseReasonMigrationComplete:
		return "MIGRATION_COMPLETE"
	case ParticipantCloseReasonStale:
		return "STALE"
	case ParticipantCloseReasonServiceRequestRemoveParticipant:
		return "SERVICE_REQUEST_REMOVE_PARTICIPANT"
	case ParticipantCloseReasonServiceRequestDeleteRoom:
		return "SERVICE_REQUEST_DELETE_ROOM"
	case ParticipantCloseReasonSimulateMigration:
		return "SIMULATE_MIGRATION"
	case ParticipantCloseReasonSimulateNodeFailure:
		return "SIMULATE_NODE_FAILURE"
	case ParticipantCloseReasonSimulateServerLeave:
		return "SIMULATE_SERVER_LEAVE"
	case ParticipantCloseReasonNegotiateFailed:
		return "NEGOTIATE_FAILED"
	case ParticipantCloseReasonMigrationRequested:
		return "MIGRATION_REQUESTED"
	case ParticipantCloseReasonOvercommitted:
		return "OVERCOMMITTED"
	case ParticipantCloseReasonPublicationError:
		return "PUBLICATION_ERROR"
	case ParticipantCloseReasonSubscriptionError:
		return "SUBSCRIPTION_ERROR"
	default:
		return fmt.Sprintf("%d", int(p))
	}
}

func (p ParticipantCloseReason) ToDisconnectReason() tc.DisconnectReason {
	switch p {
	case ParticipantCloseReasonClientRequestLeave:
		return tc.DisconnectReason_CLIENT_INITIATED
	case ParticipantCloseReasonRoomManagerStop:
		return tc.DisconnectReason_SERVER_SHUTDOWN
	case ParticipantCloseReasonVerifyFailed, ParticipantCloseReasonJoinFailed, ParticipantCloseReasonJoinTimeout:
		// expected to be connected but is not
		return tc.DisconnectReason_JOIN_FAILURE
	case ParticipantCloseReasonPeerConnectionDisconnected:
		return tc.DisconnectReason_STATE_MISMATCH
	case ParticipantCloseReasonDuplicateIdentity, ParticipantCloseReasonMigrationComplete, ParticipantCloseReasonStale:
		return tc.DisconnectReason_DUPLICATE_IDENTITY
	case ParticipantCloseReasonServiceRequestRemoveParticipant:
		return tc.DisconnectReason_PARTICIPANT_REMOVED
	case ParticipantCloseReasonServiceRequestDeleteRoom:
		return tc.DisconnectReason_ROOM_DELETED
	case ParticipantCloseReasonSimulateMigration:
		return tc.DisconnectReason_DUPLICATE_IDENTITY
	case ParticipantCloseReasonSimulateNodeFailure:
		return tc.DisconnectReason_SERVER_SHUTDOWN
	case ParticipantCloseReasonSimulateServerLeave:
		return tc.DisconnectReason_SERVER_SHUTDOWN
	case ParticipantCloseReasonOvercommitted:
		return tc.DisconnectReason_SERVER_SHUTDOWN
	case ParticipantCloseReasonNegotiateFailed, ParticipantCloseReasonPublicationError:
		return tc.DisconnectReason_STATE_MISMATCH
	default:
		// the other types will map to unknown reason
		return tc.DisconnectReason_UNKNOWN_REASON
	}
}

// ---------------------------------------------

type SignallingCloseReason int // 信令原因

const (
	SignallingCloseReasonUnknown SignallingCloseReason = iota
	SignallingCloseReasonMigration
	SignallingCloseReasonResume
	SignallingCloseReasonTransportFailure
	SignallingCloseReasonFullReconnectPublicationError
	SignallingCloseReasonFullReconnectSubscriptionError
	SignallingCloseReasonFullReconnectNegotiateFailed
	SignallingCloseReasonParticipantClose
)

func (s SignallingCloseReason) String() string {
	switch s {
	case SignallingCloseReasonUnknown:
		return "UNKNOWN"
	case SignallingCloseReasonMigration:
		return "MIGRATION"
	case SignallingCloseReasonResume:
		return "RESUME"
	case SignallingCloseReasonTransportFailure:
		return "TRANSPORT_FAILURE"
	case SignallingCloseReasonFullReconnectPublicationError:
		return "FULL_RECONNECT_PUBLICATION_ERROR"
	case SignallingCloseReasonFullReconnectSubscriptionError:
		return "FULL_RECONNECT_SUBSCRIPTION_ERROR"
	case SignallingCloseReasonFullReconnectNegotiateFailed:
		return "FULL_RECONNECT_NEGOTIATE_FAILED"
	case SignallingCloseReasonParticipantClose:
		return "PARTICIPANT_CLOSE"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}

// ---------------------------------------------

type Participant interface {
	// ID 获取参与者ID
	ID() tc.ParticipantID
	// Identity 参与者的身份ID
	Identity() tc.ParticipantIdentity
	// State 参与者状态
	State() tc.ParticipantInfo_State

	// CanSkipBroadcast 能否跳过广播
	CanSkipBroadcast() bool
	// ToProto 转化为协议对象
	ToProto() *tc.ParticipantInfo

	// SetName 设置参与者名称
	SetName(name string)
	// SetMetadata 设置元数据
	SetMetadata(metadata string)

	// IsPublisher 是否为发布者
	IsPublisher() bool
	// GetPublishedTrack 获取发布的媒体音轨
	GetPublishedTrack(sid tc.TrackID) MediaTrack
	// GetPublishedTracks 获取全部发布的音轨
	GetPublishedTracks() []MediaTrack
	// RemovePublishedTrack 移除发布的音轨
	RemovePublishedTrack(track MediaTrack, willBeResumed bool, shouldClose bool)

	//HasPermission 通过标识检查订阅者的权限。如果允许订阅者订阅，则返回true
	//到trackID为的曲目
	HasPermission(trackID tc.TrackID, subIdentity tc.ParticipantIdentity) bool

	// 权限

	// Hidden 是否隐藏
	Hidden() bool
	// IsRecorder 是否为录制者
	IsRecorder() bool

	// Start 开始
	Start()
	// Close 关闭
	Close(sendLeave bool, reason ParticipantCloseReason, isExpectedToResume bool) error

	// SubscriptionPermission 订阅者权限
	SubscriptionPermission() (*tc.SubscriptionPermission, utils.TimedVersion)

	// UpdateSubscriptionPermission 更新参与者权限
	UpdateSubscriptionPermission(subscriptionPermission *tc.SubscriptionPermission, timedVersion utils.TimedVersion, resolverByIdentity func(participantIdentity tc.ParticipantIdentity) LocalParticipant,
		resolverBySid func(participantID tc.ParticipantID) LocalParticipant) error

	// UpdateVideoLayers 更新VideoLayers
	UpdateVideoLayers(updateVideoLayers *tc.UpdateVideoLayers) error

	// DebugInfo 调试信息
	DebugInfo() map[string]interface{}
}

// -------------------------------------------------------

type ICEConnectionType string // ICE 连接类型

const (
	ICEConnectionTypeUDP     ICEConnectionType = "udp"
	ICEConnectionTypeTCP     ICEConnectionType = "tcp"
	ICEConnectionTypeTURN    ICEConnectionType = "turn"
	ICEConnectionTypeUnknown ICEConnectionType = "unknown"
)

// AddTrackParams 添加音轨参数
type AddTrackParams struct {
	// Stereo 是否为立体声
	Stereo bool
	// Red 它通常指的是一种压缩算法，即"Redundant Audio Data"（冗余音频数据）。
	// RED是一种用于音频数据传输的编码方案，旨在减少网络传输中的数据量和带宽需求。
	Red bool
}

// LocalParticipant local参与者接口
type LocalParticipant interface {
	Participant

	// ToProtoWithVersion 转换为协议对象，同时包含TimedVersion
	ToProtoWithVersion() (*tc.ParticipantInfo, utils.TimedVersion)

	GetTrailer() []byte
	GetLogger() logger.Logger
	// GetAdaptiveStream 获取是否支持自适应流
	GetAdaptiveStream() bool
	// ProtocolVersion 协议版本
	ProtocolVersion() ProtocolVersion
	// SupportSyncStreamID 是否支持同步流ID
	SupportSyncStreamID() bool
	// ConnectedAt 连接事件
	ConnectedAt() time.Time
	// IsClosed 是否连接关闭
	IsClosed() bool
	// IsReady 是否就绪
	IsReady() bool
	// IsDisconnected 是否断开连接
	IsDisconnected() bool
	// IsIdle 是否空闲
	IsIdle() bool
	// SubscriberAsPrimary 是否作为主订阅者
	SubscriberAsPrimary() bool
	// GetClientInfo 获取客户端信息
	GetClientInfo() *tc.ClientInfo
	// GetClientConfiguration 获取客户端配置信息
	GetClientConfiguration() *tc.ClientConfiguration
	// GetICEConnectionType 获取ICE连接类型
	GetICEConnectionType() ICEConnectionType
	GetBufferFactory() *buffer.Factory
	// GetPlayoutDelayConfig 获取播放延迟配置
	GetPlayoutDelayConfig() *tc.PlayoutDelay

	// SetResponseSink 设置响应接收器
	SetResponseSink(sink routing.MessageSink)
	// CloseSignalConnection 关闭信道连接
	CloseSignalConnection(reason SignallingCloseReason)
	// 	UpdateLastSeenSignal 更新最后看到的信道
	UpdateLastSeenSignal()
	// SetSignalSourceValid 设置信道源的有效性
	SetSignalSourceValid(valid bool)
	// HandleSignalSourceClose 处理信道源关闭
	HandleSignalSourceClose()

	//---------------权限-------------

	// ClaimGrants 声明的权限
	ClaimGrants() *auth.ClaimGrants
	// SetPermission 谁参与者权限
	SetPermission(permission *tc.ParticipantPermission) bool
	// CanPublishSource 判断是否可以发布音轨源
	CanPublishSource(source tc.TrackSource) bool
	// CanSubscribe 是否可以订阅
	CanSubscribe() bool
	// CanPublishData 是否能够发布数据
	CanPublishData() bool

	//-------------对等连接----------------

	// AddICECandidate 添加ICE候选
	AddICECandidate(candidate webrtc.ICECandidateInit, target tc.SignalTarget)
	// HandleOffer 处理提供的会话描述
	HandleOffer(sdp webrtc.SessionDescription)
	// AddTrack 添加音轨
	AddTrack(req *tc.AddTrackRequest)
	// SetTrackMuted 设置音轨静音
	SetTrackMuted(trackId tc.TrackID, muted bool, fromAdmin bool)

	// HandleAnswer 处理会话描述应答
	HandleAnswer(sdp webrtc.SessionDescription)
	// Negotiate 协商
	Negotiate(force bool)
	// ICERestart ICE重启
	ICERestart(iceConfig *tc.ICEConfig)
	// AddTrackToSubscriber 向订阅者添加音轨
	AddTrackToSubscriber(trackLocal webrtc.TrackLocal, params AddTrackParams) (*webrtc.RTPSender, *webrtc.RTPTransceiver, error)
	// AddTransceiverFromTrackToSubscriber 将收发器从音轨添加到订阅者
	AddTransceiverFromTrackToSubscriber(trackLocal webrtc.TrackLocal, params AddTrackParams) (*webrtc.RTPSender, *webrtc.RTPTransceiver, error)
	// RemoveTrackFromSubscriber 从订阅者移除音轨
	RemoveTrackFromSubscriber(sender *webrtc.RTPSender) error

	//-------------订阅-----------------------

	// SubscribeToTrack 订阅音轨
	SubscribeToTrack(trackID tc.TrackID)
	// UnsubscribeFromTrack 解除订阅
	UnsubscribeFromTrack(trackID tc.TrackID)
	// UpdateSubscribedTrackSettings 更新订阅者音轨设置
	UpdateSubscribedTrackSettings(trackID tc.TrackID, settings *tc.UpdateTrackSettings)
	// GetSubscribedTracks 获取订阅的音轨
	GetSubscribedTracks() []SubscribedTrack
	// VerifySubscribeParticipantInfo 更新订阅的参与者信息
	VerifySubscribeParticipantInfo(pID tc.ParticipantID, version uint32)
	//WaitUntilSubscribed 等待直到所有订阅都已解决，或者如果超时
	//已到达。如果超时过期，它将返回错误。
	WaitUntilSubscribed(timeout time.Duration) error

	//GetSubscribedParticipants 返回当前参与者订阅的参与者标识列表
	GetSubscribedParticipants() []tc.ParticipantID
	// IsSubscribedTo 是否订阅了sid
	IsSubscribedTo(sid tc.ParticipantID) bool

	// GetAudioLevel 获取音频级别
	GetAudioLevel() (smoothLevel float64, active bool)
	// GetConnectionQuality 获取连接质量
	GetConnectionQuality() *tc.ConnectionQualityInfo

	// --------------server发送消息-----------------

	// SendJoinResponse 发送连接响应
	SendJoinResponse(joinResponse *tc.JoinResponse) error
	// SendParticipantUpdate 发送参与者更新
	SendParticipantUpdate(participants []*tc.ParticipantInfo) error
	// SendSpeakerUpdate 发送发言者更新
	SendSpeakerUpdate(speakers []*tc.SpeakerInfo, force bool) error
	// SendDataPacket 发送数据包
	SendDataPacket(packet *tc.DataPacket, data []byte) error
	// SendRoomUpdate 发送房间更新
	SendRoomUpdate(room *tc.Room) error
	// SendConnectionQualityUpdate 发送连接质量更新
	SendConnectionQualityUpdate(update *tc.ConnectionQualityUpdate) error
	// SubscriptionPermissionUpdate 订阅者权限更新
	SubscriptionPermissionUpdate(publisherID tc.ParticipantIdentity, trackID tc.TrackID, allowed bool)
	// SendRefreshToken 刷新Token
	SendRefreshToken(token string) error
	// HandleReconnectAndSendResponse 处理重连并发送响应
	HandleReconnectAndSendResponse(reconnectionReason tc.ReconnectReason, reconnectResponse *tc.ReconnectResponse) error
	// IssueFullReconnect 发出完全重新连接
	IssueFullReconnect(reason ParticipantCloseReason)

	// ----------------------回调---------------------------

	// OnStateChange 状态变更事件
	OnStateChange(func(p LocalParticipant, oldState tc.ParticipantInfo_State))
	// OnMigrateStateChange 状态迁移变更事件
	OnMigrateStateChange(func(p LocalParticipant, migrateState MigrateState))
	// OnTrackPublished 远端发送音轨事件
	OnTrackPublished(func(LocalParticipant, MediaTrack))
	// OnTrackUpdated 其中的一个发布的音轨状态发生变化
	OnTrackUpdated(callback func(participant LocalParticipant, track MediaTrack))
	// OnTrackUnpublished 一个曲目不在发布
	OnTrackUnpublished(callback func(participant LocalParticipant, track MediaTrack))
	// OnParticipantUpdate 元数据或权限已更新
	OnParticipantUpdate(callback func(LocalParticipant))
	// OnDataPacket 收到数据包事件
	OnDataPacket(callback func(LocalParticipant, *tc.DataPacket))
	// OnSubscribedStatusChanged 订阅者状态发生变化
	OnSubscribedStatusChanged(fn func(publisherID tc.ParticipantID, subscribed bool))
	// OnClose 参与者关闭事件
	OnClose(callback func(participant LocalParticipant))
	// OnClaimsChanged 声明变更
	OnClaimsChanged(callback func(LocalParticipant))
	// OnReceiverReport 接受者报告事件
	OnReceiverReport(dt *sfu.DownTrack, report *rtcp.ReceiverReport)

	//---------------------------会话迁移--------------------

	// MaybeStartMigration 可能开始会话迁移
	MaybeStartMigration(force bool, onStart func()) bool
	// SetMigrateState 设置迁移状态
	SetMigrateState(s MigrateState)
	// MigrateState 获取迁移状态
	MigrateState() MigrateState
	// SetMigrateInfo 设置迁移信息
	SetMigrateInfo(previousOffer, previousAnswer *webrtc.SessionDescription)

	// UpdateMediaRTT 更新媒体的RTT
	UpdateMediaRTT(rtt uint32)
	// UpdateSignalRTT 更新信号RTT
	UpdateSignalRTT(rtt uint32)
	// CacheDownTrack 缓存下行音轨
	CacheDownTrack(trackID tc.TrackID, rtpTransceiver *webrtc.RTPTransceiver, downTrackState sfu.DownTrackState)
	// UncacheDownTrack 取消缓存下行音轨
	UncacheDownTrack(rtpTransceiver *webrtc.RTPTransceiver)

	// SetICEConfig 设置ICE配置
	SetICEConfig(iceConfig *tc.ICEConfig)
	// OnICEConfigChanged ICE配置变化
	OnICEConfigChanged(callback func(participant LocalParticipant, iceConfig *tc.ICEConfig))
	// UpdateSubscribedQuality 更新订阅者的连接质量
	UpdateSubscribedQuality(nodeID tc.NodeID, trackID tc.TrackID, maxQualities []SubscribedCodecQuality) error
	// UpdateMediaLoss 更新媒体丢包率
	UpdateMediaLoss(nodeID tc.NodeID, trackID tc.TrackID, fractionLoss uint32) error

	// ---------------下行带宽管理-------------------

	// SetSubscriberAllowPause 更新订阅者是否允许暂停
	SetSubscriberAllowPause(allowPause bool)
	// SetSubscriberChannelCapacity 设置订阅者频道容量
	SetSubscriberChannelCapacity(channelCapacity int64)

	GetPacer() pacer.Pacer
}

// MediaTrack 表示媒体曲目
type MediaTrack interface {
}
