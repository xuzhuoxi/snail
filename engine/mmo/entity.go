//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

import "sync"

type EntityType int

const (
	EntityNone EntityType = iota
	EntityChannel
	EntityUser
	EntityRoom
	EntityZone
	EntityWorld

	EntityMax
)

type IEntity interface {
	//唯一标识
	UID() string
	//昵称，显示使用
	NickName() string
	//实体类型
	EntityType() EntityType
}

type IInitEntity interface {
	//初始化实体
	InitEntity()
}

type IEntityOwner interface {
	GetOwner() string
	NoneOwner() bool

	SetOwner(ownerId string)
	ClearOwner()
}

func NewEntityOwnerSupport() *EntityOwnerSupport {
	return &EntityOwnerSupport{}
}

type EntityOwnerSupport struct {
	Owner string
	oMu   sync.RWMutex
}

func (s *EntityOwnerSupport) GetOwner() string {
	s.oMu.RLock()
	defer s.oMu.RUnlock()
	return s.Owner
}

func (s *EntityOwnerSupport) NoneOwner() bool {
	s.oMu.RLock()
	defer s.oMu.RUnlock()
	return s.Owner == ""
}

func (s *EntityOwnerSupport) SetOwner(ownerId string) {
	s.oMu.Lock()
	defer s.oMu.Unlock()
	s.Owner = ownerId
}

func (s *EntityOwnerSupport) ClearOwner() {
	s.oMu.Lock()
	defer s.oMu.Unlock()
	s.Owner = ""
}
