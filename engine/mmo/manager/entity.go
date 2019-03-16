//
//Created by xuzhuoxi
//on 2019-03-15.
//@author xuzhuoxi
//
package manager

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"github.com/xuzhuoxi/snail/engine/mmo/entity"
	"github.com/xuzhuoxi/snail/engine/mmo/index"
	"sync"
)

type IEntityCreator interface {
	//构造世界
	InitWorld(worldId string, worldName string) (basis.IWorldEntity, error)
	//构造区域
	CreateZone(zoneId string, zoneName string) (basis.IZoneEntity, error)
	//构造房间
	CreateRoomAt(roomId string, roomName string, ownerId string) (basis.IRoomEntity, error)

	//创建队伍
	CreateTeam(userId string) (basis.ITeamEntity, error)
	//创建团队
	CreateTeamCorps(teamId string) (basis.ITeamCorpsEntity, error)
	//构造频道
	CreateChannel(chanId string, chanName string) (basis.IChannelEntity, error)
}

type IEntityIndexSet interface {
	ZoneIndex() basis.IZoneIndex
	RoomIndex() basis.IRoomIndex
	UserIndex() basis.IUserIndex
	TeamIndex() basis.ITeamIndex
	TeamCorpsIndex() basis.ITeamCorpsIndex
	ChannelIndex() basis.IChannelIndex
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
	GetTeamCorps(corpsId string) (basis.ITeamCorpsEntity, bool)
	//获取频道实例
	GetChannel(chanId string) (basis.IChannelEntity, bool)
}

type IEntityManager interface {
	eventx.IEventDispatcher
	IEntityCreator
	IEntityGetter
	IEntityIndexSet
	basis.IManagerBase
	World() basis.IWorldEntity
}

func NewIEntityManager() IEntityManager {
	return NewEntityManager()
}

func NewEntityManager() IEntityManager {
	return &EntityManager{}
}

//----------------------------

type EntityManager struct {
	zoneIndex        basis.IZoneIndex
	zoneIndexMu      sync.RWMutex
	roomIndex        basis.IRoomIndex
	roomIndexMu      sync.RWMutex
	userIndex        basis.IUserIndex
	userIndexMu      sync.RWMutex
	teamIndex        basis.ITeamIndex
	teamIndexMu      sync.RWMutex
	teamCorpsIndex   basis.ITeamCorpsIndex
	teamCorpsIndexMu sync.RWMutex
	channelIndex     basis.IChannelIndex
	chanIndexMu      sync.RWMutex

	rootWorld basis.IWorldEntity
	logger    logx.ILogger
	eventx.EventDispatcher
}

func (m *EntityManager) InitManager() {
	return
}

func (m *EntityManager) DisposeManager() {
	return
}

func (m *EntityManager) SetLogger(logger logx.ILogger) {
	m.logger = logger
}

func (m *EntityManager) InitWorld(worldId string, worldName string) (basis.IWorldEntity, error) {
	if nil != m.rootWorld {
		return nil, errors.New("World is exist. ")
	}
	m.rootWorld = entity.CreateWorldEntity(worldId, worldName)
	m.rootWorld.InitEntity()
	m.addEntityEventListener(m.rootWorld)
	m.zoneIndex = index.NewIZoneIndex()
	m.roomIndex = index.NewIRoomIndex()
	m.userIndex = index.NewIUserIndex()
	m.teamIndex = index.NewITeamIndex()
	m.teamCorpsIndex = index.NewITeamCorpsIndex()
	m.channelIndex = index.NewIChannelIndex()
	return m.rootWorld, nil
}

func (m *EntityManager) CreateZone(zoneId string, zoneName string) (basis.IZoneEntity, error) {
	m.zoneIndexMu.Lock()
	defer m.zoneIndexMu.Unlock()
	if m.zoneIndex.CheckZone(zoneId) {
		return nil, errors.New("EntityManager.CreateZone Error: ZoneId(" + zoneId + ") Duplicate!")
	}
	zone := entity.NewIZoneEntity(zoneId, zoneName)
	zone.InitEntity()
	m.addEntityEventListener(zone)
	m.zoneIndex.AddZone(zone)
	zone.SetParent(m.rootWorld.UID())
	m.rootWorld.AddChild(zone)
	return zone, nil
}

func (m *EntityManager) CreateRoomAt(roomId string, roomName string, ownerId string) (basis.IRoomEntity, error) {
	m.roomIndexMu.Lock()
	defer m.roomIndexMu.Unlock()
	if m.roomIndex.CheckRoom(roomId) {
		return nil, errors.New("EntityManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate")
	}
	if "" != ownerId && !m.zoneIndex.CheckZone(ownerId) {
		return nil, errors.New("EntityManager.CreateRoomAt Error: OwnerId(" + ownerId + ") does not exist")
	}
	room := entity.NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	m.addEntityEventListener(room)
	m.roomIndex.AddRoom(room)
	room.SetParent(ownerId)
	if "" != ownerId {
		zone := m.zoneIndex.GetZone(ownerId)
		zone.AddChild(room)
	}
	return room, nil
}

func (m *EntityManager) CreateTeam(userId string) (basis.ITeamEntity, error) {
	m.teamIndexMu.Lock()
	defer m.teamIndexMu.Unlock()
	if userId == "" || !m.userIndex.CheckUser(userId) {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateTeam Error: User(%s) does not exist", userId))
	}
	team := entity.NewITeamEntity(basis.GetTeamId(), basis.TeamName, basis.MaxTeamMember)
	team.InitEntity()
	m.addEntityEventListener(team)
	m.teamIndex.AddTeam(team)
	team.AddChild(m.userIndex.GetUser(userId))
	team.SetParent(userId)
	return team, nil
}

func (m *EntityManager) CreateTeamCorps(teamId string) (basis.ITeamCorpsEntity, error) {
	m.teamCorpsIndexMu.Lock()
	defer m.teamCorpsIndexMu.Unlock()
	if teamId == "" || !m.teamIndex.CheckTeam(teamId) {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateTeamCorps Error: Team(%s) does not exist", teamId))
	}
	teamCorps := entity.NewITeamCorpsEntity(basis.GetTeamCorpsId(), basis.TeamCorpsName)
	teamCorps.InitEntity()
	m.addEntityEventListener(teamCorps)
	m.teamCorpsIndex.AddCorps(teamCorps)
	teamCorps.AddChild(m.teamIndex.GetTeam(teamId))
	teamCorps.SetParent(teamId)
	return teamCorps, nil
}

func (m *EntityManager) CreateChannel(chanId string, chanName string) (basis.IChannelEntity, error) {
	m.chanIndexMu.Lock()
	defer m.chanIndexMu.Unlock()
	if m.channelIndex.CheckChannel(chanId) {
		return nil, errors.New("EntityManager.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel := entity.NewIChannelEntity(chanId, chanName)
	channel.InitEntity()
	m.addEntityEventListener(channel)
	m.channelIndex.AddChannel(channel)
	return channel, nil
}

func (m *EntityManager) addEntityEventListener(entity basis.IEntity) {
	if dispatcher, ok := entity.(basis.IVariableSupport); ok {
		dispatcher.AddEventListener(basis.EventVariableChanged, m.onEntityVar)
	}
}

func (m *EntityManager) removeEntityEventListener(entity basis.IEntity) {
	if dispatcher, ok := entity.(basis.IVariableSupport); ok {
		dispatcher.RemoveEventListener(basis.EventVariableChanged, m.onEntityVar)
	}
}

//事件转发
func (m *EntityManager) onEntityVar(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	m.DispatchEvent(evd.EventType, m, []interface{}{evd.CurrentTarget(), evd.Data})
}

//----------------------------

func (m *EntityManager) World() basis.IWorldEntity {
	return m.rootWorld
}

func (m *EntityManager) GetZone(zoneId string) (basis.IZoneEntity, bool) {
	m.zoneIndexMu.RLock()
	defer m.zoneIndexMu.RUnlock()
	if zone := m.zoneIndex.GetZone(zoneId); nil != zone {
		return zone, true
	}
	return nil, false
}

func (m *EntityManager) GetRoom(roomId string) (basis.IRoomEntity, bool) {
	m.roomIndexMu.RLock()
	defer m.roomIndexMu.RUnlock()
	if room := m.roomIndex.GetRoom(roomId); nil != room {
		return room, true
	}
	return nil, false
}

func (m *EntityManager) GetUser(userId string) (basis.IUserEntity, bool) {
	m.userIndexMu.RLock()
	defer m.userIndexMu.RUnlock()
	if user := m.userIndex.GetUser(userId); nil != user {
		return user, true
	}
	return nil, false
}

func (m *EntityManager) GetTeam(teamId string) (basis.ITeamEntity, bool) {
	m.teamIndexMu.RLock()
	defer m.teamIndexMu.RUnlock()
	if team := m.teamIndex.GetTeam(teamId); nil != team {
		return team, true
	}
	return nil, false
}

func (m *EntityManager) GetTeamCorps(corpsId string) (basis.ITeamCorpsEntity, bool) {
	m.teamCorpsIndexMu.RLock()
	defer m.teamCorpsIndexMu.RUnlock()
	if teamCorps := m.teamCorpsIndex.GetCorps(corpsId); nil != teamCorps {
		return teamCorps, true
	}
	return nil, false
}

func (m *EntityManager) GetChannel(chanId string) (basis.IChannelEntity, bool) {
	m.chanIndexMu.RLock()
	defer m.chanIndexMu.RUnlock()
	if channel := m.channelIndex.GetChannel(chanId); nil != channel {
		return channel, true
	}
	return nil, false
}

//-----------------------

func (m *EntityManager) ZoneIndex() basis.IZoneIndex {
	return m.zoneIndex
}

func (m *EntityManager) RoomIndex() basis.IRoomIndex {
	return m.roomIndex
}

func (m *EntityManager) UserIndex() basis.IUserIndex {
	return m.userIndex
}

func (m *EntityManager) TeamIndex() basis.ITeamIndex {
	return m.teamIndex
}

func (m *EntityManager) TeamCorpsIndex() basis.ITeamCorpsIndex {
	return m.teamCorpsIndex
}

func (m *EntityManager) ChannelIndex() basis.IChannelIndex {
	return m.channelIndex
}
