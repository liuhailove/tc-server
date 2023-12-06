package routing

import (
	"context"
	"encoding/json"
	"github.com/liuhailove/tc-base-go/protocol/logger"

	"github.com/golang/protobuf/proto"
	"github.com/redis/go-redis/v9"

	"github.com/liuhailove/tc-base-go/protocol/tc"
)

// MessageSink 是编写 protobuf 消息并让 MessageSource 读取它们的抽象，
// 可能通过传输位于不同的节点上
//
//counterfeiter:generate . MessageSink
type MessageSink interface {
	// WriteMessage 写入消息
	WriteMessage(msg proto.Message) error
	// IsClosed sink是否已经关闭
	IsClosed() bool
	// Close 关闭sink
	Close()
	// ConnectionID 连接ID
	ConnectionID() tc.ConnectionID
}

//counterfeiter:generate . MessageSource
type MessageSource interface {
	// ReadChan 导出一个单向通道，使其更易于与 select 一起使用
	ReadChan() <-chan proto.Message
	IsClosed() bool
	Close()
	ConnectionID() tc.ConnectionID
}

// ParticipantInit 参与者初始化
type ParticipantInit struct {
	// Identity 参与者标识
	Identity tc.ParticipantIdentity
	// Name 参与者名称
	Name tc.ParticipantName
	// Reconnect 是否重连
	Reconnect bool
	// ReconnectReason 重连原因
	ReconnectReason tc.ReconnectReason
	// AutoSubscribe 是否自动订阅
	AutoSubscribe bool
	// Client 客户端信息
	Client *tc.ClientInfo
	// Grants 授权
	Grants *auth.ClaimsGrants
	// Region 地区
	Region string
	// AdaptiveStream 自适应流
	AdaptiveStream bool
	// ID 参与者ID
	ID tc.ParticipantID
	// SubscriberAllowPause 订阅者是否允许暂停
	SubscriberAllowPause *bool
}

// NewParticipantCallback 新参与者回调
type NewParticipantCallback func(
	ctx context.Context,
	roomName tc.RoomName,
	pi ParticipantInit,
	requestSource MessageSource,
	responseSink MessageSink,
) error

// RTCMessageCallback RTC消息回调
type RTCMessageCallback func(
	ctx context.Context,
	roomName tc.RoomName,
	identity tc.ParticipantIdentity,
	msg *tc.RTCNodeMessage,
)

// Router 允许多个节点协调参与者会话
//
//counterfeiter:generate . Router
type Router interface {
	MessageSource

	// RegisterNode 注册节点
	RegisterNode() error
	// UnregisterNode 解除注册
	UnregisterNode() error
	// RemoveDeadNodes 移除不活跃节点
	RemoveDeadNodes() error

	// ListNodes 列举节点
	ListNodes() ([]*tc.Node, error)

	// GetNodeForRoom 获取房间所在节点
	GetNodeForRoom(ctx context.Context, roomName tc.RoomName) (*tc.Node, error)
	// SetNodeForRoom 为房间设置Node
	SetNodeForRoom(ctx context.Context, roomName tc.RoomName, nodeId tc.NodeID) error
	// ClearRoomState 清除房间状态
	ClearRoomState(ctx context.Context, roomName tc.RoomName) error

	// GetRegion 获取所在地区
	GetRegion() string

	// Start 开始
	Start() error
	// Drain 取出
	Drain()
	// Stop 停止
	Stop()

	// OnNewParticipantRTC 调用 OnNewParticipantRTC 来启动新参与者的RTC连接
	OnNewParticipantRTC(callback NewParticipantCallback)

	// OnRTCMessage 执行RTC node的操作时被调用
	OnRTCMessage(callback RTCMessageCallback)
}

// MessageRouter 消息路由
type MessageRouter interface {
	// StartParticipantSignal 参与者信号连接已准备好启动
	StartParticipantSignal(ctx context.Context, roomName tc.RoomName, pi ParticipantInit) (connectionID tc.ConnectionID, reqSink MessageSink, resSource MessageSource, err error)
	// WriteParticipantRTC 向参与者或房间写入消息
	WriteParticipantRTC(ctx context.Context, roomName tc.RoomName, identity tc.ParticipantIdentity, msg *tc.RTCNodeMessage) error
	WriteRoomRTC(ctx context.Context, roomName tc.RoomName, msg *tc.RTCNodeMessage) error
}

// CreateRouter 创建路由
func CreateRouter(config *config.Config, rc redis.UniversalClient, node LocalNode, signalClient SignalClient) Router {
	lr := NewLocalRouter(node, signalClient)

	if rc != nil {
		return NewRedisRouter(config, lr, rc)
	}

	// 本地路由和存储
	logger.Infow("using single-node routing")
	return lr
}

// ToStartSession 开启会话
func (pi *ParticipantInit) ToStartSession(roomName tc.RoomName, connectID tc.ConnectionID) (*tc.StartSession, error) {
	claims, err := json.Marshal(pi.Grants)
	if err != nil {
		return nil, err
	}

	ss := &tc.StartSession{
		RoomName: string(roomName),
		Identity: string(pi.Identity),
		Name:     string(pi.Name),
		// connection ID 是为了让 RTC 节点识别将消息路由回哪里
		ConnectionId:    string(connectID),
		Reconnect:       pi.Reconnect,
		ReconnectReason: pi.ReconnectReason,
		AutoSubscribe:   pi.AutoSubscribe,
		Client:          pi.Client,
		GrantsJson:      string(claims),
		AdaptiveStream:  pi.AdaptiveStream,
		ParticipantId:   string(pi.ID),
	}
	if pi.SubscriberAllowPause != nil {
		subscriberAllowPause := *pi.SubscriberAllowPause
		// 此处原先为option可选，修改为了非option
		ss.SubscriberAllowPause = subscriberAllowPause
	}

	return ss, nil
}

// ParticipantInitFromStartSession 从会话session初始化参与者Init
func ParticipantInitFromStartSession(ss *tc.StartSession, region string) (*ParticipantInit, error) {
	claims := &auth.ClaimGrants{}
	if err := json.Unmarshal([]byte(ss.GrantsJson), claims); err != nil {
		return nil, err
	}

	pi := &ParticipantInit{
		Identity:        tc.ParticipantIdentity(ss.Identity),
		Name:            tc.ParticipantName(ss.Name),
		Reconnect:       ss.Reconnect,
		ReconnectReason: ss.ReconnectReason,
		Client:          ss.Client,
		AutoSubscribe:   ss.AutoSubscribe,
		Grants:          claims,
		Region:          region,
		AdaptiveStream:  ss.AdaptiveStream,
		ID:              tc.ParticipantID(ss.ParticipantId),
	}
	pi.SubscriberAllowPause = &ss.SubscriberAllowPause

	return pi, nil
}
