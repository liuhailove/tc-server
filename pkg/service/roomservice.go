package service

import (
	"context"

	"github.com/liuhailove/tc-base-go/protocol/tc"
)

// RoomService 单节点的room服务，通过twirp协议提供Room的增删改查功能
// 此服务注册到http server上提供服务
type RoomService struct {
}

// CreateRoom 创建房间
func (r RoomService) CreateRoom(ctx context.Context, request *tc.CreateRoomRequest) (*tc.Room, error) {
	//TODO implement me
	panic("implement me")
}

// ListRooms 列举房间
func (r RoomService) ListRooms(ctx context.Context, request *tc.ListRoomsRequest) (*tc.ListRoomsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteRoom 删除房间
func (r RoomService) DeleteRoom(ctx context.Context, request *tc.DeleteRoomRequest) (*tc.DeleteRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ListParticipants 列举参与者
func (r RoomService) ListParticipants(ctx context.Context, request *tc.ListParticipantsRequest) (*tc.ListParticipantsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetParticipant 获取参与者
func (r RoomService) GetParticipant(ctx context.Context, identity *tc.RoomParticipantIdentity) (*tc.ParticipantInfo, error) {
	//TODO implement me
	panic("implement me")
}

// RemoveParticipant 移除参与者
func (r RoomService) RemoveParticipant(ctx context.Context, identity *tc.RoomParticipantIdentity) (*tc.RemoveParticipantResponse, error) {
	//TODO implement me
	panic("implement me")
}

// MutePublishedTrack 发布音轨静音
func (r RoomService) MutePublishedTrack(ctx context.Context, request *tc.MuteRoomTrackRequest) (*tc.MuteRoomTrackResponse, error) {
	//TODO implement me
	panic("implement me")
}

// UpdateParticipant 更新参与者
func (r RoomService) UpdateParticipant(ctx context.Context, request *tc.UpdateParticipantRequest) (*tc.ParticipantInfo, error) {
	//TODO implement me
	panic("implement me")
}

// UpdateSubscriptions 更新订阅者
func (r RoomService) UpdateSubscriptions(ctx context.Context, request *tc.UpdateSubscriptionsRequest) (*tc.UpdateSubscriptionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// SendData 发送数据
func (r RoomService) SendData(ctx context.Context, request *tc.SendDataRequest) (*tc.SendDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

// UpdateRoomMetadata 更新房间原数据
func (r RoomService) UpdateRoomMetadata(ctx context.Context, request *tc.UpdateRoomMetadataRequest) (*tc.Room, error) {
	//TODO implement me
	panic("implement me")
}
