//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewIZoneIndex() basis.IZoneIndex {
	return NewZoneIndex()
}

func NewZoneIndex() *ZoneIndex {
	return &ZoneIndex{EntityIndex: *NewEntityIndex("ZoneIndex", basis.EntityZone)}
}

type ZoneIndex struct {
	EntityIndex
}

func (i *ZoneIndex) CheckZone(zoneId string) bool {
	return i.EntityIndex.Check(zoneId)
}

func (i *ZoneIndex) GetZone(zoneId string) basis.IZoneEntity {
	entity := i.EntityIndex.Get(zoneId)
	if nil != entity {
		return entity.(basis.IZoneEntity)
	}
	return nil
}

func (i *ZoneIndex) AddZone(zone basis.IZoneEntity) error {
	return i.EntityIndex.Add(zone)
}

func (i *ZoneIndex) RemoveZone(zoneId string) (basis.IZoneEntity, error) {
	c, err := i.EntityIndex.Remove(zoneId)
	if nil != c {
		return c.(basis.IZoneEntity), err
	}
	return nil, err
}

func (i *ZoneIndex) UpdateZone(zone basis.IZoneEntity) error {
	return i.EntityIndex.Update(zone)
}
