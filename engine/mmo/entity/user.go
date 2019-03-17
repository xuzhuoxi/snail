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

func NewIUserEntity(userId string, userName string) basis.IUserEntity {
	return NewUserEntity(userId, userName)
}

func NewUserEntity(userId string, userName string) *UserEntity {
	return &UserEntity{Uid: userId, Nick: userName}
}

type UserEntity struct {
	Uid  string //用户标识，唯一，内部使用
	Name string //用户名，唯一
	Nick string //用户昵称
	Addr string //用户历史或当前连接地址

	LocType basis.EntityType
	LocId   string
	locMu   sync.RWMutex

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
	e.VariableSupport = *NewVariableSupport(e)
}

func (e *UserEntity) DestroyEntity() {
}

func (e *UserEntity) GetLocation() (idType basis.EntityType, id string) {
	e.locMu.RLock()
	defer e.locMu.RUnlock()
	return e.LocType, e.LocId
}

func (e *UserEntity) SetLocation(idType basis.EntityType, id string) {
	e.locMu.Lock()
	defer e.locMu.Unlock()
	if idType != e.LocType {
		e.LocType = idType
	}
	if id != e.LocId {
		e.UserSubscriber.RemoveWhite(e.LocId)
		e.LocId = id
		e.UserSubscriber.AddWhite(e.LocId)
	}
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
