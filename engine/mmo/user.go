//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

import (
	"sync"
)

//用户实体
type IUserEntity interface {
	IEntity
	IInitEntity
	IChannelSubscriber
	IVariableSupport
	//用户名
	UserName() string

	GetLocation() (zoneId string, roomId string)
	SetZone(zoneId string, roomId string)
	SetRoom(roomId string)

	GetTeamInfo() (teamId string, corpsId string)
	SetTeam(teamId string)
	SetCorps(corpsId string)

	GetPosition() XYZ
	SetPosition(pos XYZ)
}

type UserEntity struct {
	Uid  string //用户标识，唯一，内部使用
	Name string //用户名，唯一
	Nick string //用户昵称
	Addr string //用户连接地址

	ZoneId string
	RoomId string
	locMu  sync.RWMutex

	CorpsId string
	TeamId  string
	teamMu  sync.RWMutex

	Pos   XYZ
	posMu sync.RWMutex

	ChannelSubscriber
	VariableSupport
}

func (e *UserEntity) UID() string {
	return e.Uid
}

func (e *UserEntity) UserName() string {
	return e.Name
}

func (e *UserEntity) NickName() string {
	return e.Nick
}

func (e *UserEntity) EntityType() EntityType {
	return EntityUser
}

func (e *UserEntity) InitEntity() {
	e.ChannelSubscriber = *NewChannelSubscriber()
	e.VariableSupport = *NewVariableSupport()
}

func (e *UserEntity) GetLocation() (zoneId string, roomId string) {
	e.locMu.RLock()
	defer e.locMu.RUnlock()
	return e.ZoneId, e.RoomId
}

func (e *UserEntity) SetZone(zoneId string, roomId string) {
	e.locMu.Lock()
	defer e.locMu.Unlock()
	if zoneId != e.ZoneId {
		e.ZoneId = zoneId
		e.ChannelSubscriber.AddWhiteChannel(zoneId)
	}
	if roomId != e.RoomId {
		e.RoomId = roomId
		e.ChannelSubscriber.AddWhiteChannel(roomId)
	}
}

func (e *UserEntity) SetRoom(roomId string) {
	e.locMu.Lock()
	defer e.locMu.Unlock()
	if roomId == e.RoomId {
		return
	}
	e.RoomId = roomId
	e.ChannelSubscriber.AddWhiteChannel(roomId)
}

//---------------------------------

func (e *UserEntity) GetTeamInfo() (teamId string, corpsId string) {
	e.teamMu.RLock()
	defer e.teamMu.RUnlock()
	return e.TeamId, e.CorpsId
}

func (e *UserEntity) SetTeamInfo(teamId string, corpsId string) {
	e.teamMu.Lock()
	defer e.teamMu.Unlock()
	if e.TeamId != teamId {
		e.TeamId = teamId
	}
	if e.CorpsId != corpsId {
		e.CorpsId = corpsId
	}
}

func (e *UserEntity) SetCorps(corpsId string) {
	e.teamMu.Lock()
	defer e.teamMu.Unlock()
	if e.CorpsId != corpsId || e.TeamId == "" {
		e.CorpsId = corpsId
	}
}

func (e *UserEntity) SetTeam(teamId string) {
	e.teamMu.Lock()
	defer e.teamMu.Unlock()
	if e.TeamId != teamId {
		e.TeamId = teamId
		if teamId == "" { //没有队伍，不能回防兵团中
			e.CorpsId = ""
		}
	}
}

//---------------------------------

func (e *UserEntity) GetPosition() XYZ {
	e.posMu.RLock()
	defer e.posMu.RUnlock()
	return e.Pos
}

func (e *UserEntity) SetPosition(pos XYZ) {
	e.posMu.Lock()
	defer e.posMu.Unlock()
	e.Pos = pos
}
