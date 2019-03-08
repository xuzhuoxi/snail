//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package mmo

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

//队伍实体
type ITeamEntity interface {
	IEntity
	IInitEntity
	IEntityOwner

	IVariableSupport

	//队长
	Leader() string
	//用户列表
	MemberList() []string
	//检查用户
	ContainMember(memberId string) bool
	//加入用户,进行唯一性检查
	AcceptMember(memberId string) error
	//从组中移除用户
	DropMember(memberId string) error
	//从组中移除用户
	RiseLeader(memberId string) error
	//解散队伍
	DisbandTeam() error
}

func NewITeamEntity(teamId string, teamName string, maxMember int) ITeamEntity {
	return &TeamEntity{TeamId: teamId, TeamName: teamName, MaxMember: maxMember}
}

func NewTeamEntity(teamId string, teamName string, maxMember int) *TeamEntity {
	return &TeamEntity{TeamId: teamId, TeamName: teamName, MaxMember: maxMember}
}

//常规房间
type TeamEntity struct {
	TeamId    string
	TeamName  string
	MaxMember int
	EntityOwnerSupport

	UserGroup *EntityListGroup
	VariableSupport

	teamMu sync.RWMutex
}

func (e *TeamEntity) UID() string {
	return e.TeamId
}

func (e *TeamEntity) NickName() string {
	return e.TeamName
}

func (e *TeamEntity) EntityType() EntityType {
	return EntityTeam
}

func (e *TeamEntity) InitEntity() {
	e.UserGroup = NewEntityListGroup(EntityUser)
	e.VariableSupport = *NewVariableSupport()
}

func (e *TeamEntity) Leader() string {
	e.teamMu.RLock()
	defer e.teamMu.RUnlock()
	return e.Owner
}

func (e *TeamEntity) MemberList() []string {
	return e.UserGroup.Entities()
}

func (e *TeamEntity) ContainMember(memberId string) bool {
	return e.UserGroup.ContainEntity(memberId)
}

func (e *TeamEntity) AcceptMember(memberId string) error {
	return e.UserGroup.Accept(memberId)
}

func (e *TeamEntity) DropMember(memberId string) error {
	e.teamMu.RLock()
	defer e.teamMu.RUnlock()
	err := e.UserGroup.Drop(memberId)
	if nil != err {
		return err
	}
	if memberId == e.Owner {
		if 0 == e.UserGroup.Len() {
			return e.disbandTeam()
		}
		e.SetOwner(e.UserGroup.Entities()[0])
	}
	return nil
}

func (e *TeamEntity) RiseLeader(memberId string) error {
	e.teamMu.Lock()
	defer e.teamMu.Unlock()
	if memberId == e.Owner {
		return errors.New(fmt.Sprintf("%s is already the leader", memberId))
	}
	if !e.UserGroup.ContainEntity(memberId) {
		return errors.New(fmt.Sprintf("%s is not a member", memberId))
	}
	e.SetOwner(memberId)
	return nil
}

func (e *TeamEntity) DisbandTeam() error {
	if e.UserGroup.Len() == 0 {
		return nil
	}
	return e.disbandTeam()
}

func (e *TeamEntity) disbandTeam() error {
	return nil
}
