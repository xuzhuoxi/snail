//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

type IWorldEntity interface {
	IEntity
	IInitEntity
	IChannelBehavior
	IVariableSupport

	CreateZone(zoneId string, zoneName string) (IZoneEntity, error)
	CreateRoomAt(roomId string, roomName string, ownerId string) (IRoomEntity, error)
	CreateChannel(chanId string, chanName string) (IChannelEntity, error)

	Zones() []string
	CopyZones() []string
	GetZone(zoneId string) IZoneEntity
	GetRoom(roomId string) IRoomEntity
	GetUser(userId string) IUserEntity
	GetChannel(chanId string) IChannelEntity

	JoinWorld(user IUserEntity, roomId string) error
	Transfer(userId string, toRoomId string) error
}

//-----------------------------------------------

func CreateWorldEntity() IWorldEntity {
	return &WorldEntity{}
}

type WorldEntity struct {
	WorldId         string
	WorldName       string
	VariableSupport *VariableSupport
	ChannelEntity   *ChannelEntity

	ZoneIndex    IZoneIndex
	RoomIndex    IRoomIndex
	UserIndex    IUserIndex
	ChannelIndex IChannelIndex

	zoneList []string
	areaMu   sync.RWMutex

	transferMu sync.Mutex
}

func (w *WorldEntity) UID() string {
	return w.WorldId
}

func (w *WorldEntity) NickName() string {
	return w.WorldName
}

func (w *WorldEntity) InitEntity() {
	w.ZoneIndex = NewIZoneIndex()
	w.RoomIndex = NewIRoomIndex()
	w.UserIndex = NewIUserIndex()
	w.ChannelIndex = NewIChannelIndex()

	w.VariableSupport = NewVariableSupport()
	w.ChannelEntity = NewChannelEntity(w.WorldId, w.WorldName)
	w.ChannelEntity.InitEntity()
}

func (w *WorldEntity) ChannelId() string {
	return w.ChannelEntity.ChannelId()
}

func (w *WorldEntity) MyChannel() IChannelEntity {
	return w.ChannelEntity
}

func (w *WorldEntity) TouchChannel(subscriber string) {
	w.ChannelEntity.TouchChannel(subscriber)
}

func (w *WorldEntity) UnTouchChannel(subscriber string) {
	w.ChannelEntity.UnTouchChannel(subscriber)
}

func (w *WorldEntity) Broadcast(speaker string, handler func(receiver string)) int {
	return w.ChannelEntity.Broadcast(speaker, handler)
}

func (w *WorldEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	return w.ChannelEntity.BroadcastSome(speaker, receiver, handler)
}

func (w *WorldEntity) SetVar(key string, value interface{}) {
	w.VariableSupport.SetVar(key, value)
}

func (w *WorldEntity) GetVar(key string) interface{} {
	return w.VariableSupport.GetVar(key)
}

func (w *WorldEntity) CheckVar(key string) bool {
	return w.VariableSupport.CheckVar(key)
}

func (w *WorldEntity) RemoveVar(key string) {
	w.VariableSupport.RemoveVar(key)
}

func (w *WorldEntity) CreateZone(zoneId string, zoneName string) (IZoneEntity, error) {
	w.areaMu.Lock()
	defer w.areaMu.Unlock()
	if w.ZoneIndex.CheckZone(zoneId) {
		return nil, errors.New("WorldEntity.CreateZone Error: ZoneId(" + zoneId + ") Duplicate!")
	}
	zone := NewIZoneEntity(zoneId, zoneName)
	w.ZoneIndex.AddZone(zone)
	w.zoneList = append(w.zoneList, zoneId)
	return zone, nil
}

func (w *WorldEntity) CreateRoomAt(roomId string, roomName string, ownerId string) (IRoomEntity, error) {
	w.areaMu.Lock()
	defer w.areaMu.Unlock()
	if w.RoomIndex.CheckRoom(roomId) {
		return nil, errors.New("WorldEntity.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate!")
	}
	if "" != ownerId && !w.ZoneIndex.CheckZone(ownerId) {
		return nil, errors.New("WorldEntity.CreateRoomAt Error: OwnerId(" + ownerId + ") does net exist!")
	}
	room := NewIRoomEntity(roomId, roomName)
	w.RoomIndex.AddRoom(room)
	room.SetOwner(ownerId)
	if "" != ownerId {
		zone := w.ZoneIndex.GetZone(ownerId)
		zone.AddRoom(roomId)
	}
	return room, nil
}

func (w *WorldEntity) CreateChannel(chanId string, chanName string) (IChannelEntity, error) {
	w.areaMu.Lock()
	defer w.areaMu.Unlock()
	if w.ChannelIndex.CheckChannel(chanId) {
		return nil, errors.New("WorldEntity.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel := NewIChannelEntity(chanId, chanName)
	w.ChannelIndex.AddChannel(channel)
	return channel, nil
}

func (w *WorldEntity) Zones() []string {
	return w.zoneList
}

func (w *WorldEntity) CopyZones() []string {
	return slicex.CopyString(w.zoneList)
}

func (w *WorldEntity) GetZone(zoneId string) IZoneEntity {
	return w.ZoneIndex.GetZone(zoneId)
}

func (w *WorldEntity) GetRoom(roomId string) IRoomEntity {
	return w.RoomIndex.GetRoom(roomId)
}

func (w *WorldEntity) GetUser(userId string) IUserEntity {
	return w.UserIndex.GetUser(userId)
}

func (w *WorldEntity) GetChannel(chanId string) IChannelEntity {
	return w.ChannelIndex.GetChannel(chanId)
}

func (w *WorldEntity) JoinWorld(user IUserEntity, roomId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if nil == user {
		return errors.New("WorldEntity.JoinWorld Error: user is nil")
	}
	if !w.RoomIndex.CheckRoom(roomId) {
		return errors.New("WorldEntity.JoinWorld Error: Room(" + roomId + ") does not exist")
	}
	userId := user.UID()
	if w.UserIndex.CheckUser(userId) {
		oldUser := w.UserIndex.GetUser(userId)
		w.exitCurrentRoom(oldUser)
	}
	w.UserIndex.UpdateUser(user)
	room := w.RoomIndex.GetRoom(roomId)
	room.EnterRoom(userId)
	user.SetWorldLocation(room.GetOwner(), roomId)
	return nil
}

func (w *WorldEntity) Transfer(userId string, toRoomId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if "" == userId || !w.UserIndex.CheckUser(userId) {
		return errors.New(fmt.Sprintf("JoinWorld Error: user(%s) does not exist", userId))
	}
	if "" == toRoomId || !w.RoomIndex.CheckRoom(toRoomId) {
		return errors.New(fmt.Sprintf("JoinWorld Error: Target room(%s) does not exist", toRoomId))
	}
	user := w.UserIndex.GetUser(userId)
	room := w.RoomIndex.GetRoom(toRoomId)
	if user.CurrentRoom() == toRoomId || room.ContainUser(userId) {
		return errors.New(fmt.Sprintf("JoinWorld Error: user(%s) already in the room(%s)", userId, toRoomId))
	}
	//离开当前
	err := w.exitCurrentRoom(user)
	if nil != err {
		return err
	}
	//进入新的
	if user.CurrentRoom() != toRoomId {
		user.SetWorldLocation(room.GetOwner(), toRoomId)
	}
	if !room.ContainUser(userId) {
		room.EnterRoom(userId)
	}
	return nil
}

func (w *WorldEntity) exitCurrentRoom(user IUserEntity) error {
	if nil == user {
		return errors.New("WorldEntity.exitCurrentRoom Error: user is nil")
	}
	roomId := user.CurrentRoom()
	if "" == roomId || !w.RoomIndex.CheckRoom(roomId) {
		return errors.New("WorldEntity.exitCurrentRoom Error: room is nil")
	}
	room := w.RoomIndex.GetRoom(roomId)
	userId := user.UID()
	if room.ContainUser(userId) {
		room.LeaveRoom(userId)
	}
	user.SetWorldLocation(user.CurrentZone(), "")
	return nil
}
