//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

type IZoneGroup interface {
	//区域列表
	ZoneList() []string
	//检查区域存在性
	ContainZone(zoneId string) bool
	//添加区域
	AddZone(zoneId string) error
	//移除区域
	RemoveZone(zoneId string) error
}

type IRoomGroup interface {
	//房间列表
	RoomList() []string
	//检查房间存在性
	ContainRoom(roomId string) bool
	//添加房间
	AddRoom(roomId string) error
	//移除房间
	RemoveRoom(roomId string) error
}

type IUserGroup interface {
	//用户列表
	UserList() []string
	//检查用户
	ContainUser(userId string) bool
	//加入用户,进行唯一性检查
	AcceptUser(userId string) error
	//从组中移除用户
	DropUser(userId string) error
}

//组
type IEntityGroup interface {
	//接纳实体的类型
	EntityType() EntityType
	//包含实体id
	Entities() []string
	//包含实体id
	CopyEntities() []string
	//检查实体是否属于当前组
	ContainEntity(entityId string) bool

	//加入实体到组,进行唯一性检查
	Accept(entityId string) error
	//加入实体到组,进行唯一性检查
	AcceptMulti(entityId []string) (count int, err error)
	//从组中移除实体
	Drop(entityId string) error
	//从组中移除实体
	DropMulti(entityId []string) (count int, err error)
}

func NewIEntityGroup(entityType EntityType, userMap bool) IEntityGroup {
	if userMap {
		return NewEntityMapGroup(entityType)
	} else {
		return NewEntityListGroup(entityType)
	}
}

func NewEntityListGroup(entityType EntityType) *EntityListGroup {
	return &EntityListGroup{entityType: entityType}
}

func NewEntityMapGroup(entityType EntityType) *EntityMapGroup {
	return &EntityMapGroup{entityType: entityType}
}

//------------------------------

type EntityMapGroup struct {
	entityType EntityType
	entityMap  map[string]*struct{}
	entityMu   sync.RWMutex
}

func (g *EntityMapGroup) EntityType() EntityType {
	return g.entityType
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

//---------------------------------------

type EntityListGroup struct {
	entityType EntityType
	entityList []string
	entityMu   sync.RWMutex
}

func (g *EntityListGroup) EntityType() EntityType {
	return g.entityType
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
	g.entityList = append(g.entityList, entityId)
	return nil
}

func (g *EntityListGroup) AcceptMulti(entityId []string) (count int, err error) {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityListGroup.AcceptMulti Error: len = 0")
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
