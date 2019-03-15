//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"errors"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIZoneIndex() basis.IZoneIndex {
	return &ZoneIndex{zoneMap: make(map[string]basis.IZoneEntity)}
}

func NewZoneIndex() ZoneIndex {
	return ZoneIndex{zoneMap: make(map[string]basis.IZoneEntity)}
}

type ZoneIndex struct {
	zoneMap map[string]basis.IZoneEntity
	mu      sync.RWMutex
}

func (i *ZoneIndex) CheckZone(zoneId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.checkZone(zoneId)
}

func (i *ZoneIndex) checkZone(zoneId string) bool {
	_, ok := i.zoneMap[zoneId]
	return ok
}

func (i *ZoneIndex) GetZone(zoneId string) basis.IZoneEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.zoneMap[zoneId]
}

func (i *ZoneIndex) AddZone(zone basis.IZoneEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == zone {
		return errors.New("ZoneIndex.AddZone Error: zone is nil")
	}
	zoneId := zone.UID()
	if i.checkZone(zoneId) {
		return errors.New("ZoneIndex.AddZone Error: Zone(" + zoneId + ") Duplicate")
	}
	i.zoneMap[zoneId] = zone
	return nil
}

func (i *ZoneIndex) RemoveZone(zoneId string) (basis.IZoneEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.zoneMap[zoneId]
	if ok {
		delete(i.zoneMap, zoneId)
		return e, nil
	}
	return nil, errors.New("ZoneIndex.RemoveZone Error: No Zone(" + zoneId + ")")
}

func (i *ZoneIndex) UpdateZone(zone basis.IZoneEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == zone {
		return errors.New("ZoneIndex.UpdateZone Error: zone is nil")
	}
	i.zoneMap[zone.UID()] = zone
	return nil
}
