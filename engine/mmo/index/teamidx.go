//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewITeamIndex() basis.ITeamIndex {
	return NewTeamIndex()
}

func NewTeamIndex() *TeamIndex {
	return &TeamIndex{EntityIndex: *NewEntityIndex("TeamIndex", basis.EntityTeam)}
}

type TeamIndex struct {
	EntityIndex
}

func (i *TeamIndex) CheckTeam(teamId string) bool {
	return i.EntityIndex.Check(teamId)
}

func (i *TeamIndex) GetTeam(teamId string) basis.ITeamEntity {
	entity := i.EntityIndex.Get(teamId)
	if nil != entity {
		return entity.(basis.ITeamEntity)
	}
	return nil
}

func (i *TeamIndex) AddTeam(team basis.ITeamEntity) error {
	return i.EntityIndex.Add(team)
}

func (i *TeamIndex) RemoveTeam(teamId string) (basis.ITeamEntity, error) {
	c, err := i.EntityIndex.Remove(teamId)
	if nil != c {
		return c.(basis.ITeamEntity), err
	}
	return nil, err
}

func (i *TeamIndex) UpdateTeam(team basis.ITeamEntity) error {
	return i.EntityIndex.Update(team)
}
