//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewITeamCorpsIndex() basis.ITeamCorpsIndex {
	return NewTeamCorpsIndex()
}

func NewTeamCorpsIndex() *TeamCorpsIndex {
	return &TeamCorpsIndex{EntityIndex: *NewEntityIndex("TeamCorpsIndex", basis.EntityTeamCorps)}
}

type TeamCorpsIndex struct {
	EntityIndex
}

func (i *TeamCorpsIndex) CheckCorps(corpsId string) bool {
	return i.EntityIndex.Check(corpsId)
}

func (i *TeamCorpsIndex) GetCorps(corpsId string) basis.ITeamCorpsEntity {
	entity := i.EntityIndex.Get(corpsId)
	if nil != entity {
		return entity.(basis.ITeamCorpsEntity)
	}
	return nil
}

func (i *TeamCorpsIndex) AddCorps(corps basis.ITeamCorpsEntity) error {
	return i.EntityIndex.Add(corps)
}

func (i *TeamCorpsIndex) RemoveCorps(corpsId string) (basis.ITeamCorpsEntity, error) {
	c, err := i.EntityIndex.Remove(corpsId)
	if nil != c {
		return c.(basis.ITeamCorpsEntity), err
	}
	return nil, err
}

func (i *TeamCorpsIndex) UpdateCorps(corps basis.ITeamCorpsEntity) error {
	return i.EntityIndex.Update(corps)
}
