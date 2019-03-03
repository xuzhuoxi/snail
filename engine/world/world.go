//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

import "sync"

type IWorldEntity interface {
	IEntity
	IInitEntity
	IChannelBehavior
	IVariableSupport

	Zones() []string
	GetZone(zoneId string) IZoneEntity
	GetRooms(zoneId string) []IRoomEntity

	CreateZone(zone *ZoneConfig) (IZoneEntity, error)
	CreateRoomAt(zoneId string, room *RoomConfig) (IRoomEntity, error)

	JoinWorld(user *UserEntity, zoneId, roomId string) error
	TransferWorld(user *UserEntity, zoneId, roomId string) error
	TransferRoom(user *UserEntity, roomId string) error
}

//-----------------------------------------------

type WorldEntity struct {
	WorldId   string
	WorldName string

	zoneList []string
	zoneMu   sync.RWMutex

	VariableSupport VariableSupport
	ChannelEntity   ChannelEntity
}

func (w *WorldEntity) UID() string {
	return w.WorldId
}

func (w *WorldEntity) NickName() string {
	return w.WorldName
}

func (w *WorldEntity) InitEntity() {
	w.VariableSupport = NewVariableSupport()
	w.ChannelEntity = NewChannelEntity(w.WorldId, w.WorldName)
	w.ChannelEntity.InitEntity()
}

func (w *WorldEntity) ChannelId() string {
	return w.ChannelEntity.ChannelId()
}

func (e *WorldEntity) TouchChannel(subscriber string) {
	e.ChannelEntity.TouchChannel(subscriber)
}

func (e *WorldEntity) UnTouchChannel(subscriber string) {
	e.ChannelEntity.UnTouchChannel(subscriber)
}

func (e *WorldEntity) Broadcast(speaker string, handler func(receiver string)) int {
	return e.ChannelEntity.Broadcast(speaker, handler)
}

func (e *WorldEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	return e.ChannelEntity.BroadcastSome(speaker, receiver, handler)
}

func (e *WorldEntity) SetVar(key string, value interface{}) {
	e.VariableSupport.SetVar(key, value)
}

func (e *WorldEntity) GetVar(key string) interface{} {
	return e.VariableSupport.GetVar(key)
}

func (e *WorldEntity) CheckVar(key string) bool {
	return e.VariableSupport.CheckVar(key)
}

func (e *WorldEntity) RemoveVar(key string) {
	e.VariableSupport.RemoveVar(key)
}

func (w *WorldEntity) Zones() []string {
	panic("implement me")
}

func (w *WorldEntity) GetZone(zoneId string) IZoneEntity {
	panic("implement me")
}

func (w *WorldEntity) GetRooms(zoneId string) []IRoomEntity {
	panic("implement me")
}

func (w *WorldEntity) CreateZone(zone *ZoneConfig) (IZoneEntity, error) {
	panic("implement me")
}

func (w *WorldEntity) CreateRoomAt(zoneId string, room *RoomConfig) (IRoomEntity, error) {
	panic("implement me")
}

func (w *WorldEntity) JoinWorld(user *UserEntity, zoneId, roomId string) error {
	panic("implement me")
}

func (w *WorldEntity) TransferWorld(user *UserEntity, zoneId, roomId string) error {
	panic("implement me")
}

func (w *WorldEntity) TransferRoom(user *UserEntity, roomId string) error {
	panic("implement me")
}
