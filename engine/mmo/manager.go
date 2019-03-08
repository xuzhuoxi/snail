//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

type IEntityCreator interface {
	//构造世界
	CreateWorld(worldId string, worldName string)
	//构造区域
	CreateZone(zoneId string, zoneName string) (IZoneEntity, error)
	//构造房间
	CreateRoomAt(roomId string, roomName string, ownerId string) (IRoomEntity, error)

	//创建队伍
	CreateTeam(userId string) (ITeamEntity, error)
	//创建团队
	CreateCorps(teamId string) (ITeamCorpsEntity, error)
	//构造频道
	CreateChannel(chanId string, chanName string) (IChannelEntity, error)
}

type IEntityGetter interface {
	//获取区域实例
	GetZone(zoneId string) (IZoneEntity, bool)
	//获取房间实例
	GetRoom(roomId string) (IRoomEntity, bool)
	//获取用户实例
	GetUser(userId string) (IUserEntity, bool)
	//获取队伍实例
	GetTeam(teamId string) (ITeamEntity, bool)
	//获取队伍实例
	GetCorps(corpsId string) (ITeamCorpsEntity, bool)
	//获取频道实例
	GetChannel(chanId string) (IChannelEntity, bool)
}

//type IChannelManager interface {
//	//订阅频道
//	TouchChannel(chanId string, subscriber string)
//	//取消频道订阅
//	UnTouchChannel(chanId string, subscriber string)
//}

type IUserBehavior interface {
	//加入世界
	EnterWorld(user IUserEntity, roomId string) error
	//离开世界
	ExitWorld(userId string) error
	//在世界转移
	Transfer(userId string, toRoomId string) error
}

type IWorldManager interface {
	IEntityCreator
	//IChannelManager
	IEntityGetter
	IUserBehavior
}

type WorldManager struct {
	ZoneIndex    IZoneIndex
	RoomIndex    IRoomIndex
	UserIndex    IUserIndex
	TeamIndex    ITeamIndex
	CorpsIndex   ITeamCorpsIndex
	ChannelIndex IChannelIndex

	world      IWorldEntity
	createMu   sync.RWMutex
	transferMu sync.Mutex
}

func (w *WorldManager) CreateWorld(worldId string, worldName string) {
	if nil != w.world {
		return
	}
	w.world = CreateWorldEntity(worldId, worldName)
	w.world.InitEntity()
	w.ZoneIndex = NewIZoneIndex()
	w.RoomIndex = NewIRoomIndex()
	w.UserIndex = NewIUserIndex()
	w.TeamIndex = NewITeamIndex()
	w.CorpsIndex = NewITeamCorpsIndex()
	w.ChannelIndex = NewIChannelIndex()
}

func (w *WorldManager) CreateZone(zoneId string, zoneName string) (IZoneEntity, error) {
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

func (w *WorldManager) CreateTeam(userId string) (ITeamEntity, error) {
	w.createMu.Lock()
	defer w.createMu.Unlock()
	panic("implement me")
}

func (w *WorldManager) CreateCorps(teamId string) (ITeamCorpsEntity, error) {
	panic("implement me")
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

//----------------------------

func (w *WorldManager) GetZone(zoneId string) (IZoneEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	if zone := w.ZoneIndex.GetZone(zoneId); nil != zone {
		return zone, true
	}
	return nil, false
}

func (w *WorldManager) GetRoom(roomId string) (IRoomEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	if room := w.RoomIndex.GetRoom(roomId); nil != room {
		return room, true
	}
	return nil, false
}

func (w *WorldManager) GetUser(userId string) (IUserEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	if user := w.UserIndex.GetUser(userId); nil != user {
		return user, true
	}
	return nil, false
}

func (w *WorldManager) GetTeam(teamId string) (ITeamEntity, bool) {
	panic("implement me")
}

func (w *WorldManager) GetCorps(corpsId string) (ITeamCorpsEntity, bool) {
	panic("implement me")
}

func (w *WorldManager) GetChannel(chanId string) (IChannelEntity, bool) {
	w.createMu.RLock()
	defer w.createMu.RUnlock()
	if channel := w.ChannelIndex.GetChannel(chanId); nil != channel {
		return channel, true
	}
	return nil, false
}

//----------------------------

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
	room.AcceptUser(userId)
	user.SetZone(room.GetOwner(), roomId)
	return nil
}

func (w *WorldManager) ExitWorld(userId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if "" == userId || w.UserIndex.CheckUser(userId) {
		return errors.New("WorldManager.ExitWorld Error: User() does not exist")
	}
	user := w.UserIndex.GetUser(userId)
	_, roomId := user.GetLocation()
	if room := w.RoomIndex.GetRoom(roomId); nil != room {
		room.DropUser(userId)
	}
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
	_, roomId := user.GetLocation()
	if roomId == toRoomId {
		return errors.New(fmt.Sprintf("EnterWorld Error: user(%s) already in the room(%s)", userId, toRoomId))
	}
	w.exitCurrentRoom(user)
	toRoom := w.RoomIndex.GetRoom(toRoomId)
	toRoom.AcceptUser(userId)
	user.SetZone(toRoom.GetOwner(), toRoomId)
	return nil
}

func (w *WorldManager) exitCurrentRoom(user IUserEntity) error {
	_, roomId := user.GetLocation()
	if "" == roomId || !w.RoomIndex.CheckRoom(roomId) {
		return errors.New("WorldManager.exitCurrentRoom Error: room is nil")
	}
	room := w.RoomIndex.GetRoom(roomId)
	userId := user.UID()
	if room.ContainUser(userId) {
		room.DropUser(userId)
	}
	user.SetRoom("")
	return nil
}

//----------------------------------

func (w *WorldManager) Broadcast(speaker string, broadcastType BroadcastType, handler func(entity IUserEntity)) error {
	if !w.UserIndex.CheckUser(speaker) {
		return errors.New(fmt.Sprintf("Speaker(%s) does not exist", speaker))
	}
	userEntity := w.UserIndex.GetUser(speaker)
	zoneId, roomId := userEntity.GetLocation()
	switch broadcastType {
	case BroadcastWorld:
		w.broadcastWorld(userEntity, handler)
	case BroadcastZone:
		w.broadcastZone(userEntity, zoneId, handler)
	case BroadcastRoom:
		w.broadcastRoom(userEntity, roomId, handler)
	}
	return nil
}

func (w *WorldManager) broadcastWorld(speaker IUserEntity, handler func(entity IUserEntity)) {
	if speaker.OnBlackChannel(w.world.UID()) {
		return
	}
	zones := w.world.ZoneList()
	for _, zoneId := range zones {
		if speaker.OnBlackChannel(zoneId) {
			continue
		}
		w.broadcastZone(speaker, zoneId, handler)
	}
}

func (w *WorldManager) broadcastZone(speaker IUserEntity, zoneId string, handler func(entity IUserEntity)) {
	if "" == zoneId || !w.ZoneIndex.CheckZone(zoneId) || !speaker.OnBlackChannel(zoneId) {
		return
	}
	rooms := w.ZoneIndex.GetZone(zoneId).RoomList()
	for _, roomId := range rooms {
		if speaker.OnBlackChannel(roomId) {
			continue
		}
		w.broadcastRoom(speaker, zoneId, handler)
	}
}

func (w *WorldManager) broadcastRoom(speaker IUserEntity, roomId string, handler func(entity IUserEntity)) {
	if "" == roomId || !w.RoomIndex.CheckRoom(roomId) || speaker.OnBlackChannel(roomId) {
		return
	}
	copyUsers := slicex.CopyString(w.RoomIndex.GetRoom(roomId).UserList())
	for _, userId := range copyUsers {
		if speaker.OnBlackChannel(userId) || !w.UserIndex.CheckUser(userId) {
			continue
		}
		handler(w.UserIndex.GetUser(userId))
	}
}

//--------------------------------------

//func (w *WorldManager) TouchChannel(chanId string, subscriber string) {
//	if channel := w.ChannelIndex.GetChannel(chanId); nil != channel {
//		channel.TouchChannel(subscriber)
//	}
//}
//
//func (w *WorldManager) UnTouchChannel(chanId string, subscriber string) {
//	if channel := w.ChannelIndex.GetChannel(chanId); nil != channel {
//		channel.UnTouchChannel(subscriber)
//	}
//}
