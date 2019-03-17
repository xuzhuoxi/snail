//
//Created by xuzhuoxi
//on 2019-03-17.
//@author xuzhuoxi
//
package index

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIEntityIndex(indexName string, entityType basis.EntityType) basis.IEntityIndex {
	return NewEntityIndex(indexName, entityType)
}

func NewEntityIndex(indexName string, entityType basis.EntityType) *EntityIndex {
	return &EntityIndex{indexName: indexName, entityType: entityType, entityMap: make(map[string]basis.IEntity)}
}

type EntityIndex struct {
	indexName  string
	entityType basis.EntityType
	entityMap  map[string]basis.IEntity
	entityMu   sync.RWMutex
}

func (i *EntityIndex) EntityType() basis.EntityType {
	return i.entityType
}

func (i *EntityIndex) Check(id string) bool {
	i.entityMu.RLock()
	defer i.entityMu.RUnlock()
	return i.check(id)
}

func (i *EntityIndex) check(id string) bool {
	_, ok := i.entityMap[id]
	return ok
}

func (i *EntityIndex) Get(id string) basis.IEntity {
	i.entityMu.RLock()
	defer i.entityMu.RUnlock()
	return i.entityMap[id]
}

func (i *EntityIndex) Add(entity basis.IEntity) error {
	i.entityMu.Lock()
	defer i.entityMu.Unlock()
	if nil == entity {
		//return errors.New(i.indexName + ".Add Error: entity is nil")
		return errors.New(fmt.Sprintf("%s.Add Error: entity is nil", i.indexName))
	}
	if !i.entityType.Include(entity.EntityType()) {
		//return errors.New(i.indexName + ".Add Error: Type is not included")
		return errors.New(fmt.Sprintf("%s.Add Error: Type is not included", i.indexName))
	}
	id := entity.UID()
	if i.check(id) {
		//return errors.New(i.indexName + ".Add Error: Id(" + id + ") Duplicate")
		return errors.New(fmt.Sprintf("%s.Add Error: Id(%s) Duplicate", i.indexName, id))
	}
	i.entityMap[id] = entity
	return nil
}

func (i *EntityIndex) Remove(id string) (basis.IEntity, error) {
	i.entityMu.Lock()
	defer i.entityMu.Unlock()
	e, ok := i.entityMap[id]
	if ok {
		delete(i.entityMap, id)
		return e, nil
	}
	//return nil, errors.New(i.indexName + ".Remove Error: No Entity(" + id + ")")
	return nil, errors.New(fmt.Sprintf("%s.Remove Error: No Entity(%s)", i.indexName, id))
}

func (i *EntityIndex) Update(entity basis.IEntity) error {
	i.entityMu.Lock()
	defer i.entityMu.Unlock()
	if nil == entity {
		//return errors.New(i.indexName + ".Update Error: entity is nil")
		return errors.New(fmt.Sprintf("%s.Update Error: entity is nil", i.indexName))
	}
	if !i.entityType.Include(entity.EntityType()) {
		//return errors.New(i.indexName + ".Update Error: Type is not included")
		return errors.New(fmt.Sprintf("%s.Update Error: Type is not included", i.indexName))
	}
	i.entityMap[entity.UID()] = entity
	return nil
}
