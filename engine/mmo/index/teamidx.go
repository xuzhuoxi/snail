//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package index

import (
	"github.com/pkg/errors"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewITeamIndex() basis.ITeamIndex {
	return &TeamIndex{teamMap: make(map[string]basis.ITeamEntity)}
}

func NewTeamIndex() TeamIndex {
	return TeamIndex{teamMap: make(map[string]basis.ITeamEntity)}
}

type TeamIndex struct {
	teamMap map[string]basis.ITeamEntity
	mu      sync.RWMutex
}

func (i *TeamIndex) CheckTeam(teamId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.teamMap[teamId]
	return ok
}

func (i *TeamIndex) GetTeam(teamId string) basis.ITeamEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.teamMap[teamId]
}

func (i *TeamIndex) AddTeam(team basis.ITeamEntity) error {
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

func (i *TeamIndex) RemoveTeam(teamId string) (basis.ITeamEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.teamMap[teamId]
	if ok {
		delete(i.teamMap, teamId)
		return e, nil
	}
	return nil, errors.New("TeamIndex.RemoveTeam Error: No Team(" + teamId + ")")
}

func (i *TeamIndex) UpdateTeam(team basis.ITeamEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == team {
		return errors.New("TeamIndex.UpdateTeam Error: team is nil")
	}
	i.teamMap[team.UID()] = team
	return nil
}
