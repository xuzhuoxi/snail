//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewIWorldIndex() basis.IWorldIndex {
	return NewWorldIndex()
}

func NewWorldIndex() *WorldIndex {
	return &WorldIndex{EntityIndex: *NewEntityIndex("WorldIndex", basis.EntityWorld)}
}

type WorldIndex struct {
	EntityIndex
}

func (i *WorldIndex) CheckWorld(worldId string) bool {
	return i.EntityIndex.Check(worldId)
}

func (i *WorldIndex) GetWorld(worldId string) basis.IWorldEntity {
	entity := i.EntityIndex.Get(worldId)
	if nil != entity {
		return entity.(basis.IWorldEntity)
	}
	return nil
}

func (i *WorldIndex) AddWorld(world basis.IWorldEntity) error {
	return i.EntityIndex.Add(world)
}

func (i *WorldIndex) RemoveWorld(world string) (basis.IWorldEntity, error) {
	c, err := i.EntityIndex.Remove(world)
	if nil != c {
		return c.(basis.IWorldEntity), err
	}
	return nil, err
}

func (i *WorldIndex) UpdateWorld(world basis.IWorldEntity) error {
	return i.EntityIndex.Update(world)
}
