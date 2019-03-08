//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package mmo

import (
	"github.com/pkg/errors"
	"sync"
)

//队伍索引
type ITeamIndex interface {
	//检查Team是否存在
	CheckTeam(teamId string) bool
	//获取Team
	GetTeam(teamId string) ITeamEntity
	//添加一个新Team到索引中
	AddTeam(team ITeamEntity) error
	//从索引中移除一个Team
	RemoveTeam(teamId string) (ITeamEntity, error)
	//从索引中更新一个Team
	UpdateTeam(team ITeamEntity) error
}

func NewITeamIndex() ITeamIndex {
	return &TeamIndex{teamMap: make(map[string]ITeamEntity)}
}

func NewTeamIndex() TeamIndex {
	return TeamIndex{teamMap: make(map[string]ITeamEntity)}
}

type TeamIndex struct {
	teamMap map[string]ITeamEntity
	mu      sync.RWMutex
}

func (i *TeamIndex) CheckTeam(teamId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.teamMap[teamId]
	return ok
}

func (i *TeamIndex) GetTeam(teamId string) ITeamEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.teamMap[teamId]
}

func (i *TeamIndex) AddTeam(team ITeamEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == team {
		return errors.New("TeamIndex.AddTeam Error: team is nil")
	}
	teamId := team.UID()
	if i.CheckTeam(teamId) {
		return errors.New("TeamIndex.AddTeam Error: Team(" + teamId + ") Duplicate")
	}
	i.teamMap[teamId] = team
	return nil
}

func (i *TeamIndex) RemoveTeam(teamId string) (ITeamEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.teamMap[teamId]
	if ok {
		delete(i.teamMap, teamId)
		return e, nil
	}
	return nil, errors.New("TeamIndex.RemoveTeam Error: No Team(" + teamId + ")")
}

func (i *TeamIndex) UpdateTeam(team ITeamEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == team {
		return errors.New("TeamIndex.UpdateTeam Error: team is nil")
	}
	i.teamMap[team.UID()] = team
	return nil
}
