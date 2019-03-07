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
	IInitEntity
	IEntityOwner

	IUserGroup
	IVariableSupport
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
	EntityOwnerSupport

	UserGroup *EntityListGroup
	VariableSupport
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
	e.UserGroup = NewEntityListGroup(EntityUser)
	e.VariableSupport = *NewVariableSupport()
}

func (e *RoomEntity) UserList() []string {
	return e.UserGroup.Entities()
}

func (e *RoomEntity) ContainUser(userId string) bool {
	return e.UserGroup.ContainEntity(userId)
}

func (e *RoomEntity) AcceptUser(userId string) error {
	return e.UserGroup.Accept(userId)
}

func (e *RoomEntity) DropUser(userId string) error {
	return e.UserGroup.Drop(userId)
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
