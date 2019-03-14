//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package entity

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIEntityChildSupport() basis.IEntityChild {
	return &EntityChildSupport{}
}

func NewEntityChildSupport() *EntityChildSupport {
	return &EntityChildSupport{}
}

type EntityChildSupport struct {
	Owner string
	oMu   sync.RWMutex
}

func (s *EntityChildSupport) GetParent() string {
	s.oMu.RLock()
	defer s.oMu.RUnlock()
	return s.Owner
}

func (s *EntityChildSupport) NoneParent() bool {
	s.oMu.RLock()
	defer s.oMu.RUnlock()
	return s.Owner == ""
}

func (s *EntityChildSupport) SetParent(parentId string) {
	s.oMu.Lock()
	defer s.oMu.Unlock()
	s.Owner = parentId
}

func (s *EntityChildSupport) ClearParent() {
	s.oMu.Lock()
	defer s.oMu.Unlock()
	s.Owner = ""
}
