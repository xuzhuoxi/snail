//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"fmt"
	"sync"
)

type IWorldManager interface {
	//构造世界
	CreateWorld()
	//构造区域
	CreateZone(zoneId string, zoneName string) (IZoneEntity, error)
	//构造房间
	CreateRoomAt(roomId string, roomName string, ownerId string) (IRoomEntity, error)
	//构造频道
	CreateChannel(chanId string, chanName string) (IChannelEntity, error)
	//获取区域实例
	GetZone(zoneId string) (IZoneEntity, bool)
	//获取房间实例
	GetRoom(roomId string) (IRoomEntity, bool)
	//获取用户实例
	GetUser(userId string) (IUserEntity, bool)
	//获取频道实例
	GetChannel(chanId string) (IChannelEntity, bool)

	//加入世界
	EnterWorld(user IUserEntity, roomId string) error
	//离开世界
	ExitWorld(userId string) error
	//在世界转移
	Transfer(userId string, toRoomId string) error

	//订阅频道
	TouchChannel(chanId string, subscriber string)
	//取消频道订阅
	UnTouchChannel(chanId string, subscriber string)
}

type WorldManager struct {
	ZoneIndex    IZoneIndex
	RoomIndex    IRoomIndex
	UserIndex    IUserIndex
	ChannelIndex IChannelIndex

	world      IWorldEntity
	createMu   sync.RWMutex
	transferMu sync.Mutex
}

func (w *WorldManager) CreateWorld() {
	w.world = CreateWorldEntity()
	w.world.InitEntity()
	w.ZoneIndex = NewIZoneIndex()
	w.RoomIndex = NewIRoomIndex()
	w.UserIndex = NewIUserIndex()
	w.ChannelIndex = NewIChannelIndex()
}

func (w *WorldManager) CreateZone(zoneId string, zoneName string) (IZoneEntity, error) {
	w.createMu.Lock()
	defer w.createMu.Unlock()
	if w.ZoneIndex.CheckZone(zoneId) {
		return nil, errors.New("WorldManager.CreateZone Error: ZoneId(" + zoneId + ") Duplicate!")
	}
	zone := NewIZoneEntity(zoneId, zoneName)
	zone.InitEntity()
	w.ZoneIndex.AddZone(zone)
	w.world.AddZone(zoneId)
	return zone, nil
}

func (w *WorldManager) CreateRoomAt(roomId string, roomName string, ownerId string) (IRoomEntity, error) {
	w.createMu.Lock()
	defer w.createMu.Unlock()
	if w.RoomIndex.CheckRoom(roomId) {
		return nil, errors.New("WorldManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate!")
	}
	if "" != ownerId && !w.ZoneIndex.CheckZone(ownerId) {
		return nil, errors.New("WorldManager.CreateRoomAt Error: OwnerId(" + ownerId + ") does net exist!")
	}
	room := NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	w.RoomIndex.AddRoom(room)
	room.SetOwner(ownerId)
	if "" != ownerId {
		zone := w.ZoneIndex.GetZone(ownerId)
		zone.AddRoom(roomId)
	}
	return room, nil
}

func (w *WorldManager) CreateChannel(chanId string, chanName string) (IChannelEntity, error) {
	w.createMu.Lock()
	defer w.createMu.Unlock()
	if w.ChannelIndex.CheckChannel(chanId) {
		return nil, errors.New("WorldEntity.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel := NewIChannelEntity(chanId, chanName)
	w.ChannelIndex.AddChannel(channel)
	return channel, nil
}

func (w *WorldManager) GetZone(zoneId string) (IZoneEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	zone := w.ZoneIndex.GetZone(zoneId)
	if nil == zone {
		return nil, false
	}
	return zone, true
}

func (w *WorldManager) GetRoom(roomId string) (IRoomEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	room := w.RoomIndex.GetRoom(roomId)
	if nil == room {
		return nil, false
	}
	return room, true
}

func (w *WorldManager) GetUser(userId string) (IUserEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	user := w.UserIndex.GetUser(userId)
	if nil == user {
		return nil, false
	}
	return user, true
}

func (w *WorldManager) GetChannel(chanId string) (IChannelEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	channel := w.ChannelIndex.GetChannel(chanId)
	if nil == channel {
		return nil, false
	}
	return channel, true
}

func (w *WorldManager) EnterWorld(user IUserEntity, roomId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if nil == user {
		return errors.New("WorldManager.EnterWorld Error: user is nil")
	}
	if !w.RoomIndex.CheckRoom(roomId) {
		return errors.New("WorldManager.EnterWorld Error: Room(" + roomId + ") does not exist")
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
	if zone, ok := w.GetZone(room.GetOwner()); ok { //加入区频道
		zone.TouchChannel(userId)
	}
	w.world.TouchChannel(userId) //加入世界频道
	return nil
}

func (w *WorldManager) ExitWorld(userId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if "" == userId || w.UserIndex.CheckUser(userId) {
		return errors.New("WorldManager.ExitWorld Error: User() does not exist")
	}
	user := w.UserIndex.GetUser(userId)
	roomId := user.CurrentRoom()
	if room := w.RoomIndex.GetRoom(roomId); nil != room {
		room.LeaveRoom(userId)
	}
	zoneId := user.CurrentZone()
	if zone := w.ZoneIndex.GetZone(zoneId); nil != zone {
		zone.UnTouchChannel(userId)
	}
	w.world.UnTouchChannel(userId)
	return nil
}

func (w *WorldManager) Transfer(userId string, toRoomId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if "" == userId || !w.UserIndex.CheckUser(userId) {
		return errors.New(fmt.Sprintf("EnterWorld Error: user(%s) does not exist", userId))
	}
	if "" == toRoomId || !w.RoomIndex.CheckRoom(toRoomId) {
		return errors.New(fmt.Sprintf("EnterWorld Error: Target room(%s) does not exist", toRoomId))
	}
	user := w.UserIndex.GetUser(userId)
	room := w.RoomIndex.GetRoom(toRoomId)
	if user.CurrentRoom() == toRoomId || room.ContainUser(userId) {
		return errors.New(fmt.Sprintf("EnterWorld Error: user(%s) already in the room(%s)", userId, toRoomId))
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

func (w *WorldManager) exitCurrentRoom(user IUserEntity) error {
	if nil == user {
		return errors.New("WorldManager.exitCurrentRoom Error: user is nil")
	}
	roomId := user.CurrentRoom()
	if "" == roomId || !w.RoomIndex.CheckRoom(roomId) {
		return errors.New("WorldManager.exitCurrentRoom Error: room is nil")
	}
	room := w.RoomIndex.GetRoom(roomId)
	userId := user.UID()
	if room.ContainUser(userId) {
		room.LeaveRoom(userId)
	}
	user.SetWorldLocation(user.CurrentZone(), "")
	return nil
}

func (w *WorldManager) TouchChannel(chanId string, subscriber string) {
	if channel := w.ChannelIndex.GetChannel(chanId); nil != channel {
		channel.TouchChannel(subscriber)
	}
}

func (w *WorldManager) UnTouchChannel(chanId string, subscriber string) {
	if channel := w.ChannelIndex.GetChannel(chanId); nil != channel {
		channel.UnTouchChannel(subscriber)
	}
}
