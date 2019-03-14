//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package entity

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIMapEntityContainer(maxCount int) basis.IEntityContainer {
	return NewMapEntityContainer(maxCount)
}

func NewIListEntityContainer(maxCount int) basis.IEntityContainer {
	return NewListEntityContainer(maxCount)
}

func NewMapEntityContainer(maxCount int) *MapEntityContainer {
	return &MapEntityContainer{maxCount: maxCount, entities: make(map[string]basis.IEntity)}
}

func NewListEntityContainer(maxCount int) *ListEntityContainer {
	return &ListEntityContainer{maxCount: maxCount}
}

//--------------------

type MapEntityContainer struct {
	maxCount    int
	entities    map[string]basis.IEntity
	containerMu sync.RWMutex
}

func (c *MapEntityContainer) NumChildren() int {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	return len(c.entities)
}

func (c *MapEntityContainer) Full() bool {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	return c.isFull()
}

func (c *MapEntityContainer) Contains(entity basis.IEntity) (isContains bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	_, isContains = c.entities[entity.UID()]
	return
}

func (c *MapEntityContainer) ContainsById(entityId string) (isContains bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	_, isContains = c.entities[entityId]
	return
}

func (c *MapEntityContainer) GetChildById(entityId string) (entity basis.IEntity, ok bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	entity, ok = c.entities[entityId]
	return
}

func (c *MapEntityContainer) ReplaceChildInto(entity basis.IEntity) error {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	if c.isFull() {
		return errors.New("Container is full ")
	}
	c.entities[entity.UID()] = entity
	return nil
}

func (c *MapEntityContainer) AddChild(entity basis.IEntity) error {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, isContains := c.entities[id]
	if isContains {
		return errors.New(fmt.Sprintf("Entity(%s) is already in the container", id))
	}
	if c.isFull() {
		return errors.New("Container is full ")
	}
	c.entities[id] = entity
	return nil
}

func (c *MapEntityContainer) RemoveChild(entity basis.IEntity) error {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, isContains := c.entities[id]
	if !isContains {
		return errors.New(fmt.Sprintf("Entity(%s) does not exist in the container", id))
	}
	delete(c.entities, id)
	return nil
}

func (c *MapEntityContainer) RemoveChildById(entityId string) (entity basis.IEntity, ok bool) {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	entity, ok = c.entities[entityId]
	if ok {
		delete(c.entities, entityId)
	}
	return
}

func (c *MapEntityContainer) ForEachChildren(each func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool)) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	if 0 == len(c.entities) {
		return
	}
	for _, entity := range c.entities {
		child := entity
		interruptCurrent, interruptRecurse := each(child)
		if interruptCurrent {
			return
		}
		if interruptRecurse {
			continue
		}
		if container, ok := entity.(basis.IEntityContainer); ok {
			container.ForEachChildren(each)
		}
	}
}

func (c *MapEntityContainer) ForEachChildrenByType(entityType basis.EntityType, each func(child basis.IEntity), recurse bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	if 0 == len(c.entities) {
		return
	}
	if recurse {
		for _, entity := range c.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
			if container, ok := entity.(basis.IEntityContainer); ok {
				container.ForEachChildrenByType(entityType, each, true)
			}
		}
	} else {
		for _, entity := range c.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
		}
	}
}

func (c *MapEntityContainer) isFull() bool {
	return c.maxCount > 0 && c.maxCount <= len(c.entities)
}

//--------------------

type ListEntityContainer struct {
	maxCount    int
	entities    []basis.IEntity
	containerMu sync.RWMutex
}

func (c *ListEntityContainer) NumChildren() int {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	return len(c.entities)
}

func (c *ListEntityContainer) Full() bool {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	return c.isFull()
}

func (c *ListEntityContainer) Contains(entity basis.IEntity) (isContains bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	if nil == entity {
		return false
	}
	e, _, ok := c.contains(entity.UID())
	if ok {
		return basis.EntityEqual(e, entity)
	}
	return false
}

func (c *ListEntityContainer) ContainsById(entityId string) (isContains bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	_, _, isContains = c.contains(entityId)
	return
}

func (c *ListEntityContainer) GetChildById(entityId string) (entity basis.IEntity, ok bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	entity, _, ok = c.contains(entityId)
	return
}

func (c *ListEntityContainer) ReplaceChildInto(entity basis.IEntity) error {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, index, isContains := c.contains(id)
	if isContains {
		c.entities[index] = entity
	} else {
		if c.isFull() {
			return errors.New("Container is full ")
		}
		c.entities = append(c.entities, entity)
	}
	return nil
}

func (c *ListEntityContainer) AddChild(entity basis.IEntity) error {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, _, isContains := c.contains(id)
	if isContains {
		return errors.New(fmt.Sprintf("Entity(%s) is already in the container", id))
	}
	if c.isFull() {
		return errors.New("Container is full ")
	}
	c.entities = append(c.entities, entity)
	return nil
}

func (c *ListEntityContainer) RemoveChild(entity basis.IEntity) error {
	c.containerMu.Lock()
	defer c.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, index, isContains := c.contains(id)
	if !isContains {
		return errors.New(fmt.Sprintf("Entity(%s) does not exist in the container", id))
	}
	c.entities = append(c.entities[:index], c.entities[index+1:]...)
	return nil
}

func (c *ListEntityContainer) RemoveChildById(entityId string) (entity basis.IEntity, ok bool) {
	var index int
	entity, index, ok = c.contains(entityId)
	if ok {
		c.entities = append(c.entities[:index], c.entities[index+1:]...)
	}
	return
}

func (c *ListEntityContainer) ForEachChildren(each func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool)) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	if 0 == len(c.entities) {
		return
	}
	for _, entity := range c.entities {
		child := entity
		interruptCurrent, interruptRecurse := each(child)
		if interruptCurrent {
			return
		}
		if interruptRecurse {
			continue
		}
		if container, ok := entity.(basis.IEntityContainer); ok {
			container.ForEachChildren(each)
		}
	}
}

func (c *ListEntityContainer) ForEachChildrenByType(entityType basis.EntityType, each func(child basis.IEntity), recurse bool) {
	c.containerMu.RLock()
	defer c.containerMu.RUnlock()
	if 0 == len(c.entities) {
		return
	}
	if recurse {
		for _, entity := range c.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
			if container, ok := entity.(basis.IEntityContainer); ok {
				container.ForEachChildrenByType(entityType, each, true)
			}
		}
	} else {
		for _, entity := range c.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
		}
	}
}

func (c *ListEntityContainer) contains(entityId string) (entity basis.IEntity, index int, isContains bool) {
	for index = 0; index < len(c.entities); index++ {
		if c.entities[index].UID() == entityId {
			entity = c.entities[index]
			isContains = true
			return
		}
	}
	return nil, -1, false
}

func (c *ListEntityContainer) isFull() bool {
	return c.maxCount > 0 && c.maxCount <= len(c.entities)
}
