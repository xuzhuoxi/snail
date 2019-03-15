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

func NewITeamCorpsIndex() basis.ITeamCorpsIndex {
	return &TeamCorpsIndex{corpsMap: make(map[string]basis.ITeamCorpsEntity)}
}

func NewTeamCorpsIndex() TeamCorpsIndex {
	return TeamCorpsIndex{corpsMap: make(map[string]basis.ITeamCorpsEntity)}
}

type TeamCorpsIndex struct {
	corpsMap map[string]basis.ITeamCorpsEntity
	mu       sync.RWMutex
}

func (i *TeamCorpsIndex) CheckCorps(corpsId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.checkCorps(corpsId)
}

func (i *TeamCorpsIndex) checkCorps(corpsId string) bool {
	_, ok := i.corpsMap[corpsId]
	return ok
}

func (i *TeamCorpsIndex) GetCorps(corpsId string) basis.ITeamCorpsEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.corpsMap[corpsId]
}

func (i *TeamCorpsIndex) AddCorps(corps basis.ITeamCorpsEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == corps {
		return errors.New("TeamCorpsIndex.AddCorps Error: corps is nil")
	}
	corpsId := corps.UID()
	if i.checkCorps(corpsId) {
		return errors.New("TeamCorpsIndex.AddCorps Error: Corps(" + corpsId + ") Duplicate")
	}
	i.corpsMap[corpsId] = corps
	return nil
}

func (i *TeamCorpsIndex) RemoveCorps(corpsId string) (basis.ITeamCorpsEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.corpsMap[corpsId]
	if ok {
		delete(i.corpsMap, corpsId)
		return e, nil
	}
	return nil, errors.New("TeamCorpsIndex.RemoveCorps Error: No Corps(" + corpsId + ")")
}

func (i *TeamCorpsIndex) UpdateCorps(corps basis.ITeamCorpsEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == corps {
		return errors.New("TeamCorpsIndex.UpdateCorps Error: corps is nil")
	}
	i.corpsMap[corps.UID()] = corps
	return nil
}
