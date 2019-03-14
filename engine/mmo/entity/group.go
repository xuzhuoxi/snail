//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package entity

import (
	"errors"
	"github.com/xuzhuoxi/infra-go/slicex"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIEntityGroup(entityType basis.EntityType, userMap bool) basis.IEntityGroup {
	if userMap {
		return NewEntityMapGroup(entityType)
	} else {
		return NewEntityListGroup(entityType)
	}
}

func NewEntityListGroup(entityType basis.EntityType) *EntityListGroup {
	return &EntityListGroup{entityType: entityType}
}

func NewEntityMapGroup(entityType basis.EntityType) *EntityMapGroup {
	return &EntityMapGroup{entityType: entityType}
}

//------------------------------

type EntityMapGroup struct {
	entityType basis.EntityType
	entityMap  map[string]*struct{}
	max        int
	entityMu   sync.RWMutex
}

func (g *EntityMapGroup) EntityType() basis.EntityType {
	return g.entityType
}

func (g *EntityMapGroup) MaxLen() int {
	return g.max
}

func (g *EntityMapGroup) Len() int {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	return len(g.entityMap)
}

func (g *EntityMapGroup) IsFull() bool {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	return g.isFull()
}

func (g *EntityMapGroup) Entities() []string {
	return g.CopyEntities()
}

func (g *EntityMapGroup) CopyEntities() []string {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	var rs []string
	for key, _ := range g.entityMap {
		rs = append(rs, key)
	}
	return rs
}

func (g *EntityMapGroup) ContainEntity(entityId string) bool {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	_, ok := g.entityMap[entityId]
	return ok
}

func (g *EntityMapGroup) Accept(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	_, ok := g.entityMap[entityId]
	if ok {
		return errors.New("EntityMapGroup.Accept Error: Entity(" + entityId + ") Duplicate")
	}
	if g.isFull() {
		return errors.New("EntityMapGroup.Accept Error: Group is Full")
	}
	g.entityMap[entityId] = nil
	return nil
}

func (g *EntityMapGroup) AcceptMulti(entityId []string) (count int, err error) {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityMapGroup.AcceptMulti Error: len = 0")
	}
	for _, id := range entityId {
		_, ok := g.entityMap[id]
		if ok && nil != err {
			err = errors.New("EntityMapGroup.AcceptMulti Error: Entity Duplicate")
			continue
		}
		count++
		g.entityMap[id] = nil
	}
	return
}

func (g *EntityMapGroup) Drop(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	_, ok := g.entityMap[entityId]
	if !ok {
		return errors.New("EntityMapGroup.Drop Error: No Entity(" + entityId + ")")
	}
	delete(g.entityMap, entityId)
	return nil
}

func (g *EntityMapGroup) DropMulti(entityId []string) (count int, err error) {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityMapGroup.DropMulti Error: len = 0")
	}
	if g.max > 0 && len(g.entityMap)+len(entityId) > g.max {
		return 0, errors.New("EntityMapGroup.AcceptMulti Error: Overcrowding")
	}
	for _, id := range entityId {
		_, ok := g.entityMap[id]
		if !ok && nil != err {
			err = errors.New("EntityMapGroup.DropMulti Error: No Entity")
			continue
		}
		count++
		delete(g.entityMap, id)
	}
	return
}

func (g *EntityMapGroup) isFull() bool {
	return g.max > 0 && len(g.entityMap) >= g.max
}

//---------------------------------------

type EntityListGroup struct {
	entityType basis.EntityType
	entityList []string
	max        int
	entityMu   sync.RWMutex
}

func (g *EntityListGroup) EntityType() basis.EntityType {
	return g.entityType
}

func (g *EntityListGroup) MaxLen() int {
	return g.max
}

func (g *EntityListGroup) Len() int {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	return len(g.entityList)
}

func (g *EntityListGroup) IsFull() bool {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	return g.isFull()
}

func (g *EntityListGroup) Entities() []string {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	return g.entityList
}

func (g *EntityListGroup) CopyEntities() []string {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	return slicex.CopyString(g.entityList)
}

func (g *EntityListGroup) ContainEntity(entityId string) bool {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	_, ok := slicex.IndexString(g.entityList, entityId)
	return ok
}

func (g *EntityListGroup) Accept(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	_, ok := slicex.IndexString(g.entityList, entityId)
	if ok {
		return errors.New("EntityListGroup.Accept Error: Entity(" + entityId + ") Duplicate")
	}
	if g.isFull() {
		return errors.New("EntityListGroup.Accept Error: Group is Full")
	}
	g.entityList = append(g.entityList, entityId)
	return nil
}

func (g *EntityListGroup) AcceptMulti(entityId []string) (count int, err error) {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityListGroup.AcceptMulti Error: len = 0")
	}
	if g.max > 0 && len(g.entityList)+len(entityId) > g.max {
		return 0, errors.New("EntityListGroup.AcceptMulti Error: Overcrowding")
	}
	for _, id := range entityId {
		_, ok := slicex.IndexString(g.entityList, id)
		if ok && nil != err {
			err = errors.New("EntityListGroup.AcceptMulti Error: Entity Duplicate")
			continue
		}
		count++
		g.entityList = append(g.entityList, id)
	}
	return
}

func (g *EntityListGroup) Drop(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	index, ok := slicex.IndexString(g.entityList, entityId)
	if !ok {
		return errors.New("EntityListGroup.Drop Error: No Entity(" + entityId + ")")
	}
	g.entityList = append(g.entityList[:index], g.entityList[index+1:]...)
	return nil
}

func (g *EntityListGroup) DropMulti(entityId []string) (count int, err error) {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityListGroup.DropMulti Error: len = 0")
	}
	for _, id := range entityId {
		index, ok := slicex.IndexString(g.entityList, id)
		if !ok && nil != err {
			err = errors.New("EntityListGroup.DropMulti Error: No Entity")
			continue
		}
		count++
		g.entityList = append(g.entityList[:index], g.entityList[index+1:]...)
	}
	return
}

func (g *EntityListGroup) isFull() bool {
	return g.max > 0 && len(g.entityList) >= g.max
}
