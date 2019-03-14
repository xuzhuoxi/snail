//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package manager

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"github.com/xuzhuoxi/snail/engine/mmo/entity"
	"github.com/xuzhuoxi/snail/engine/mmo/index"
	"sync"
)

type IEntityCreator interface {
	//构造世界
	CreateWorld(worldId string, worldName string)
	//构造区域
	CreateZone(zoneId string, zoneName string) (basis.IZoneEntity, error)
	//构造房间
	CreateRoomAt(roomId string, roomName string, ownerId string) (basis.IRoomEntity, error)

	//创建队伍
	CreateTeam(userId string) (basis.ITeamEntity, error)
	//创建团队
	CreateCorps(teamId string) (basis.ITeamCorpsEntity, error)
	//构造频道
	CreateChannel(chanId string, chanName string) (basis.IChannelEntity, error)
}

type IEntityGetter interface {
	//获取区域实例
	GetZone(zoneId string) (basis.IZoneEntity, bool)
	//获取房间实例
	GetRoom(roomId string) (basis.IRoomEntity, bool)
	//获取用户实例
	GetUser(userId string) (basis.IUserEntity, bool)
	//获取队伍实例
	GetTeam(teamId string) (basis.ITeamEntity, bool)
	//获取队伍实例
	GetCorps(corpsId string) (basis.ITeamCorpsEntity, bool)
	//获取频道实例
	GetChannel(chanId string) (basis.IChannelEntity, bool)
}

//type IChannelManager interface {
//	//订阅频道
//	TouchChannel(chanId string, subscriber string)
//	//取消频道订阅
//	UnTouchChannel(chanId string, subscriber string)
//}

type IUserBehavior interface {
	//加入世界
	EnterWorld(user basis.IUserEntity, roomId string) error
	//离开世界
	ExitWorld(userId string) error
	//在世界转移
	Transfer(userId string, toRoomId string) error
}

type IBroadcastWorld interface {
	BroadcastWorld()
}

type IWorldManager interface {
	netx.ISockServerSetter
	IEntityCreator
	//IChannelManager
	IEntityGetter
	IUserBehavior
}

func NewIWorldManager() IWorldManager {
	return NewWorldManager()
}

func NewWorldManager() *WorldManager {
	return &WorldManager{}
}

//------------------------------------

type WorldManager struct {
	ZoneIndex      basis.IZoneIndex
	RoomIndex      basis.IRoomIndex
	UserIndex      basis.IUserIndex
	TeamIndex      basis.ITeamIndex
	TeamCorpsIndex basis.ITeamCorpsIndex
	ChannelIndex   basis.IChannelIndex
	indexMu        sync.RWMutex

	world      basis.IWorldEntity
	transferMu sync.Mutex

	server netx.ISockServer
}

func (w *WorldManager) SetServer(server netx.ISockServer) {
	w.server = server
}

func (w *WorldManager) CreateWorld(worldId string, worldName string) {
	if nil != w.world {
		return
	}
	w.world = entity.CreateWorldEntity(worldId, worldName)
	w.world.InitEntity()
	w.ZoneIndex = index.NewIZoneIndex()
	w.RoomIndex = index.NewIRoomIndex()
	w.UserIndex = index.NewIUserIndex()
	w.TeamIndex = index.NewITeamIndex()
	w.TeamCorpsIndex = index.NewITeamCorpsIndex()
	w.ChannelIndex = index.NewIChannelIndex()
}

func (w *WorldManager) CreateZone(zoneId string, zoneName string) (basis.IZoneEntity, error) {
	if w.ZoneIndex.CheckZone(zoneId) {
		return nil, errors.New("WorldManager.CreateZone Error: ZoneId(" + zoneId + ") Duplicate!")
	}
	zone := entity.NewIZoneEntity(zoneId, zoneName)
	zone.InitEntity()
	w.ZoneIndex.AddZone(zone)
	w.world.AddChild(zone)
	return zone, nil
}

func (w *WorldManager) CreateRoomAt(roomId string, roomName string, ownerId string) (basis.IRoomEntity, error) {
	if w.RoomIndex.CheckRoom(roomId) {
		return nil, errors.New("WorldManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate")
	}
	if "" != ownerId && !w.ZoneIndex.CheckZone(ownerId) {
		return nil, errors.New("WorldManager.CreateRoomAt Error: OwnerId(" + ownerId + ") does not exist")
	}
	room := entity.NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	w.RoomIndex.AddRoom(room)
	room.SetParent(ownerId)
	if "" != ownerId {
		zone := w.ZoneIndex.GetZone(ownerId)
		zone.AddChild(room)
	}
	return room, nil
}

func (w *WorldManager) CreateTeam(userId string) (basis.ITeamEntity, error) {
	w.indexMu.Lock()
	defer w.indexMu.Unlock()
	if userId == "" || !w.UserIndex.CheckUser(userId) {
		return nil, errors.New(fmt.Sprintf("WorldManager.CreateTeam Error: User(%s) does not exist", userId))
	}
	team := entity.NewITeamEntity(basis.GetTeamId(), basis.TeamName, basis.MaxTeamMember)
	w.TeamIndex.AddTeam(team)
	team.AddChild(w.UserIndex.GetUser(userId))
	team.SetParent(userId)
	return team, nil
}

func (w *WorldManager) CreateCorps(teamId string) (basis.ITeamCorpsEntity, error) {
	w.indexMu.Lock()
	defer w.indexMu.Unlock()
	if teamId == "" || !w.TeamIndex.CheckTeam(teamId) {
		return nil, errors.New(fmt.Sprintf("WorldManager.CreateCorps Error: Team(%s) does not exist", teamId))
	}
	teamCorps := entity.NewITeamCorpsEntity(basis.GetTeamCorpsId(), basis.TeamCorpsName)
	w.TeamCorpsIndex.AddCorps(teamCorps)
	teamCorps.AddChild(w.TeamIndex.GetTeam(teamId))
	teamCorps.SetParent(teamId)
	return teamCorps, nil
}

func (w *WorldManager) CreateChannel(chanId string, chanName string) (basis.IChannelEntity, error) {
	w.indexMu.Lock()
	defer w.indexMu.Unlock()
	if w.ChannelIndex.CheckChannel(chanId) {
		return nil, errors.New("WorldEntity.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel := entity.NewIChannelEntity(chanId, chanName)
	w.ChannelIndex.AddChannel(channel)
	return channel, nil
}

//----------------------------

func (w *WorldManager) GetZone(zoneId string) (basis.IZoneEntity, bool) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	if zone := w.ZoneIndex.GetZone(zoneId); nil != zone {
		return zone, true
	}
	return nil, false
}

func (w *WorldManager) GetRoom(roomId string) (basis.IRoomEntity, bool) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	if room := w.RoomIndex.GetRoom(roomId); nil != room {
		return room, true
	}
	return nil, false
}

func (w *WorldManager) GetUser(userId string) (basis.IUserEntity, bool) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	if user := w.UserIndex.GetUser(userId); nil != user {
		return user, true
	}
	return nil, false
}

func (w *WorldManager) GetTeam(teamId string) (basis.ITeamEntity, bool) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	if team := w.TeamIndex.GetTeam(teamId); nil != team {
		return team, true
	}
	return nil, false
}

func (w *WorldManager) GetCorps(corpsId string) (basis.ITeamCorpsEntity, bool) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	if teamCorps := w.TeamCorpsIndex.GetCorps(corpsId); nil != teamCorps {
		return teamCorps, true
	}
	return nil, false
}

func (w *WorldManager) GetChannel(chanId string) (basis.IChannelEntity, bool) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	if channel := w.ChannelIndex.GetChannel(chanId); nil != channel {
		return channel, true
	}
	return nil, false
}

//----------------------------

func (w *WorldManager) EnterWorld(user basis.IUserEntity, roomId string) error {
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
	room.AddChild(user)
	user.SetZone(room.GetParent(), roomId)
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
		user.DestroyEntity()
		room.RemoveChild(user)
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
	toRoom.AddChild(user)
	user.SetZone(toRoom.GetParent(), toRoomId)
	return nil
}

func (w *WorldManager) exitCurrentRoom(user basis.IUserEntity) error {
	_, roomId := user.GetLocation()
	if "" == roomId || !w.RoomIndex.CheckRoom(roomId) {
		return errors.New("WorldManager.exitCurrentRoom Error: room is nil")
	}
	room := w.RoomIndex.GetRoom(roomId)
	err := room.RemoveChild(user)
	if nil != err {
		return err
	}
	user.SetRoom("")
	return nil
}

//----------------------------------

func (w *WorldManager) Broadcast(speaker string, broadcastType entity.BroadcastType, handler func(entity basis.IUserEntity)) error {
	//if !w.UserIndex.CheckUser(speaker) {
	//	return errors.New(fmt.Sprintf("Speaker(%s) does not exist", speaker))
	//}
	//userEntity := w.UserIndex.GetUser(speaker)
	//zoneId, roomId := userEntity.GetLocation()
	//switch broadcastType {
	//case BroadcastWorld:
	//	w.broadcastWorld(userEntity, handler)
	//case BroadcastZone:
	//	w.broadcastZone(userEntity, zoneId, handler)
	//case BroadcastRoom:
	//	w.broadcastRoom(userEntity, roomId, handler)
	//}
	return nil
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
