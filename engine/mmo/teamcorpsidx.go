//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"sync"
)

type ITeamCorpsIndex interface {
	//检查Corps是否存在
	CheckCorps(corpsId string) bool
	//获取Corps
	GetCorps(corpsId string) ITeamCorpsEntity
	//添加一个新Corps到索引中
	AddCorps(corps ITeamCorpsEntity) error
	//从索引中移除一个Corps
	RemoveCorps(corpsId string) (ITeamCorpsEntity, error)
	//更新一个新Corps到索引中
	UpdateCorps(corps ITeamCorpsEntity) error
}

func NewITeamCorpsIndex() ITeamCorpsIndex {
	return &TeamCorpsIndex{corpsMap: make(map[string]ITeamCorpsEntity)}
}

func NewTeamCorpsIndex() TeamCorpsIndex {
	return TeamCorpsIndex{corpsMap: make(map[string]ITeamCorpsEntity)}
}

type TeamCorpsIndex struct {
	corpsMap map[string]ITeamCorpsEntity
	mu       sync.RWMutex
}

func (i *TeamCorpsIndex) CheckCorps(corpsId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.corpsMap[corpsId]
	return ok
}

func (i *TeamCorpsIndex) GetCorps(corpsId string) ITeamCorpsEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.corpsMap[corpsId]
}

func (i *TeamCorpsIndex) AddCorps(corps ITeamCorpsEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == corps {
		return errors.New("TeamCorpsIndex.AddCorps Error: corps is nil")
	}
	corpsId := corps.UID()
	if i.CheckCorps(corpsId) {
		return errors.New("TeamCorpsIndex.AddCorps Error: Corps(" + corpsId + ") Duplicate")
	}
	i.corpsMap[corpsId] = corps
	return nil
}

func (i *TeamCorpsIndex) RemoveCorps(corpsId string) (ITeamCorpsEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.corpsMap[corpsId]
	if ok {
		delete(i.corpsMap, corpsId)
		return e, nil
	}
	return nil, errors.New("TeamCorpsIndex.RemoveCorps Error: No Corps(" + corpsId + ")")
}

func (i *TeamCorpsIndex) UpdateCorps(corps ITeamCorpsEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == corps {
		return errors.New("TeamCorpsIndex.UpdateCorps Error: corps is nil")
	}
	i.corpsMap[corps.UID()] = corps
	return nil
}
