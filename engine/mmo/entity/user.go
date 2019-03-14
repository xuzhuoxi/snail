//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package entity

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

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

	Pos   basis.XYZ
	posMu sync.RWMutex

	UserSubscriber
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

func (e *UserEntity) EntityType() basis.EntityType {
	return basis.EntityUser
}

func (e *UserEntity) InitEntity() {
	e.UserSubscriber = *NewUserSubscriber()
	e.VariableSupport = *NewVariableSupport()
}

func (e *UserEntity) DestroyEntity() {
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
		e.UserSubscriber.AddWhite(zoneId)
	}
	if roomId != e.RoomId {
		e.RoomId = roomId
		e.UserSubscriber.AddWhite(roomId)
	}
}

func (e *UserEntity) SetRoom(roomId string) {
	e.locMu.Lock()
	defer e.locMu.Unlock()
	if roomId == e.RoomId {
		return
	}
	e.RoomId = roomId
	e.UserSubscriber.AddWhite(roomId)
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

func (e *UserEntity) GetPosition() basis.XYZ {
	e.posMu.RLock()
	defer e.posMu.RUnlock()
	return e.Pos
}

func (e *UserEntity) SetPosition(pos basis.XYZ) {
	e.posMu.Lock()
	defer e.posMu.Unlock()
	e.Pos = pos
}
