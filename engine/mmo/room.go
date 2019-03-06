//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"sync"
)

//房间实体
type IRoomEntity interface {
	IEntity
	IEntityOwner
	IInitEntity
	IChannelBehavior
	IVariableSupport

	//进入房间
	EnterRoom(userId string) error
	//离开房间
	LeaveRoom(userId string) error
	//包含用户
	ContainUser(userId string) bool
	//用户列表
	UserList() []string
}

//房间索引
type IRoomIndex interface {
	//检查Room是否存在
	CheckRoom(roomId string) bool
	//获取Room
	GetRoom(roomId string) IRoomEntity
	//添加一个新Room到索引中
	AddRoom(room IRoomEntity) error
	//从索引中移除一个Room
	RemoveRoom(roomId string) (IRoomEntity, error)
	//从索引中更新一个Room
	UpdateRoom(room IRoomEntity) error
}

//-----------------------------------------------

func NewIRoomEntity(roomId string, roomName string) IRoomEntity {
	return &RoomEntity{RoomId: roomId, RoomName: roomName, MaxMember: 0}
}

func NewIAOBRoomEntity(roomId string, roomName string) IRoomEntity {
	return &AOBRoomEntity{RoomEntity: RoomEntity{RoomId: roomId, RoomName: roomName, MaxMember: 0}}
}

func NewRoomEntity(roomId string, roomName string) *RoomEntity {
	return &RoomEntity{RoomId: roomId, RoomName: roomName, MaxMember: 0}
}

func NewAOBRoomEntity(roomId string, roomName string) *AOBRoomEntity {
	return &AOBRoomEntity{RoomEntity: RoomEntity{RoomId: roomId, RoomName: roomName, MaxMember: 0}}
}

type RoomConfig struct {
	Id        string
	Name      string
	Private   bool
	MaxMember int
}

//范围广播房间，适用于mmo大型场景
type AOBRoomEntity struct {
	RoomEntity
}

func (e *AOBRoomEntity) Broadcast(speaker string, handler func(receiver string)) int {
	panic("+++++++++++++++++++")
}

//常规房间
type RoomEntity struct {
	RoomId    string
	RoomName  string
	MaxMember int
	UserGroup *EntityListGroup

	EntityOwnerSupport
	ChannelEntity   *ChannelEntity
	VariableSupport *VariableSupport

	mutex sync.Mutex
}

func (e *RoomEntity) UID() string {
	return e.RoomId
}

func (e *RoomEntity) NickName() string {
	return e.RoomName
}

func (e *RoomEntity) EntityType() EntityType {
	return EntityRoom
}

func (e *RoomEntity) InitEntity() {
	e.UserGroup = NewEntityListGroup(e.RoomId, e.RoomName, EntityUser)
	e.ChannelEntity = NewChannelEntity(e.RoomId, e.RoomName)
	e.VariableSupport = NewVariableSupport()
	e.ChannelEntity.InitEntity()
}

func (e *RoomEntity) EnterRoom(userId string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.UserGroup.AppendEntity(userId)
	e.ChannelEntity.TouchChannel(userId)
	return nil
}

func (e *RoomEntity) LeaveRoom(userId string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.ChannelEntity.UnTouchChannel(userId)
	e.UserGroup.RemoveEntity(userId)
	return nil
}

func (e *RoomEntity) ContainUser(userId string) bool {
	return e.UserGroup.CheckEntity(userId)
}

func (e *RoomEntity) UserList() []string {
	return e.UserGroup.Entities()
}

func (e *RoomEntity) MyChannel() IChannelEntity {
	return e.ChannelEntity
}

func (e *RoomEntity) TouchChannel(subscriber string) {
	e.ChannelEntity.TouchChannel(subscriber)
}

func (e *RoomEntity) UnTouchChannel(subscriber string) {
	e.ChannelEntity.UnTouchChannel(subscriber)
}

func (e *RoomEntity) Broadcast(speaker string, handler func(receiver string)) int {
	return e.ChannelEntity.Broadcast(speaker, handler)
}

func (e *RoomEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	return e.ChannelEntity.BroadcastSome(speaker, receiver, handler)
}

func (e *RoomEntity) SetVar(key string, value interface{}) {
	e.VariableSupport.SetVar(key, value)
}

func (e *RoomEntity) GetVar(key string) interface{} {
	return e.VariableSupport.GetVar(key)
}

func (e *RoomEntity) CheckVar(key string) bool {
	return e.VariableSupport.CheckVar(key)
}

func (e *RoomEntity) RemoveVar(key string) {
	e.VariableSupport.RemoveVar(key)
}

//-----------------------------------------------

func NewIRoomIndex() IRoomIndex {
	return &RoomIndex{roomMap: make(map[string]IRoomEntity)}
}

func NewRoomIndex() RoomIndex {
	return RoomIndex{roomMap: make(map[string]IRoomEntity)}
}

type RoomIndex struct {
	roomMap map[string]IRoomEntity
	mu      sync.RWMutex
}

func (i *RoomIndex) CheckRoom(roomId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.roomMap[roomId]
	return ok
}

func (i *RoomIndex) GetRoom(roomId string) IRoomEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.roomMap[roomId]
}

func (i *RoomIndex) AddRoom(room IRoomEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == room {
		return errors.New("RoomIndex.AddRoom Error: room is nil")
	}
	roomId := room.UID()
	if i.CheckRoom(roomId) {
		return errors.New("RoomIndex.AddRoom Error: Room(" + roomId + ") Duplicate")
	}
	i.roomMap[roomId] = room
	return nil
}

func (i *RoomIndex) RemoveRoom(roomId string) (IRoomEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.roomMap[roomId]
	if ok {
		delete(i.roomMap, roomId)
		return e, nil
	}
	return nil, errors.New("RoomIndex.RemoveRoom Error: No Room(" + roomId + ")")
}

func (i *RoomIndex) UpdateRoom(room IRoomEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == room {
		return errors.New("RoomIndex.UpdateRoom Error: room is nil")
	}
	i.roomMap[room.UID()] = room
	return nil
}
