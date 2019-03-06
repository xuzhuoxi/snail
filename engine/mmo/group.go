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

//组
type IEntityGroup interface {
	//组id
	GroupId() string
	//组名称
	GroupName() string
	//接纳实体的类型
	EntityType() EntityType
	//包含实体id
	Entities() []string
	//包含实体id
	CopyEntities() []string
	//检查实体是否属于当前组
	CheckEntity(entityId string) bool
	//加入实体到组,进行唯一性检查
	AppendEntity(entityId string) error
	//从组中移除实体
	RemoveEntity(entityId string) error
}

func NewIEntityGroup(groupId string, groupName string, entityType EntityType, userMap bool) IEntityGroup {
	if userMap {
		return NewEntityMapGroup(groupId, groupName, entityType)
	} else {
		return NewEntityListGroup(groupId, groupName, entityType)
	}
}

func NewEntityListGroup(groupId string, groupName string, entityType EntityType) *EntityListGroup {
	return &EntityListGroup{groupId: groupId, groupName: groupName, entityType: entityType}
}

func NewEntityMapGroup(groupId string, groupName string, entityType EntityType) *EntityMapGroup {
	return &EntityMapGroup{groupId: groupId, groupName: groupName, entityType: entityType}
}

//------------------------------

var defaultEntityMapGroupValue = struct{}{}

type EntityMapGroup struct {
	groupId    string
	groupName  string
	entityType EntityType
	entityMap  map[string]struct{}
	entityMu   sync.RWMutex
}

func (g *EntityMapGroup) GroupId() string {
	return g.groupId
}

func (g *EntityMapGroup) GroupName() string {
	return g.groupName
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

func (g *EntityMapGroup) CheckEntity(entityId string) bool {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	_, ok := g.entityMap[entityId]
	return ok
}

func (g *EntityMapGroup) AppendEntity(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	_, ok := g.entityMap[entityId]
	if ok {
		return errors.New("EntityMapGroup.AppendEntity Error: Entity(" + entityId + ") Duplicate")
	}
	g.entityMap[entityId] = defaultEntityMapGroupValue
	return nil
}

func (g *EntityMapGroup) RemoveEntity(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	_, ok := g.entityMap[entityId]
	if !ok {
		return errors.New("EntityMapGroup.RemoveEntity Error: No Entity(" + entityId + ")")
	}
	delete(g.entityMap, entityId)
	return nil
}

//---------------------------------------

type EntityListGroup struct {
	groupId    string
	groupName  string
	entityType EntityType
	entityList []string
	entityMu   sync.RWMutex
}

func (g *EntityListGroup) EntityType() EntityType {
	return g.entityType
}

func (g *EntityListGroup) GroupId() string {
	return g.groupId
}

func (g *EntityListGroup) GroupName() string {
	return g.groupName
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

func (g *EntityListGroup) CheckEntity(entityId string) bool {
	g.entityMu.RLock()
	defer g.entityMu.RUnlock()
	_, ok := slicex.IndexString(g.entityList, entityId)
	return ok
}

func (g *EntityListGroup) AppendEntity(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	_, ok := slicex.IndexString(g.entityList, entityId)
	if ok {
		return errors.New("EntityListGroup.AppendEntity Error: Entity(" + entityId + ") Duplicate")
	}
	g.entityList = append(g.entityList, entityId)
	return nil
}

func (g *EntityListGroup) RemoveEntity(entityId string) error {
	g.entityMu.Lock()
	defer g.entityMu.Unlock()
	index, ok := slicex.IndexString(g.entityList, entityId)
	if !ok {
		return errors.New("EntityListGroup.RemoveEntity Error: No Entity(" + entityId + ")")
	}
	g.entityList = append(g.entityList[:index], g.entityList[index+1:]...)
	return nil
}
