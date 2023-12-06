package types

type ProtocolVersion int

const CurrentProtocol = 10

func (v ProtocolVersion) SupportsPackedStreamId() bool {
	return v > 0
}

func (v ProtocolVersion) SupportsProtobuf() bool {
	return v > 0
}

func (v ProtocolVersion) HandlesDataPackets() bool {
	return v > 1
}

// SubscriberAsPrimary 指示客户端将订阅者连接作为主要发起者
func (v ProtocolVersion) SubscriberAsPrimary() bool {
	return v > 2
}

// SupportsSpeakerChanged 如果客户端处理发言人信息增量，而不是完整列表
func (v ProtocolVersion) SupportsSpeakerChanged() bool {
	return v > 2
}

// SupportsTransceiverReuse 如果支持收发器复用，则优化 SDP 大小
func (v ProtocolVersion) SupportsTransceiverReuse() bool {
	return v > 3
}

// SupportsConnectionQuality - 避免频繁发送较低协议版本的 ConnectionQuality 更新
func (v ProtocolVersion) SupportsConnectionQuality() bool {
	return v > 4
}

// SupportsSessionMigrate 是否支持会话迁移
func (v ProtocolVersion) SupportsSessionMigrate() bool {
	return v > 5
}

func (v ProtocolVersion) SupportsICELite() bool {
	return v > 5
}

// SupportsUnpublish 是否支持取消发布
func (v ProtocolVersion) SupportsUnpublish() bool {
	return v > 6
}
