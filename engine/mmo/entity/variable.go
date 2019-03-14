//
//Created by xuzhuoxi
//on 2019-03-03.
//@author xuzhuoxi
//
package entity

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIVariableSupport() basis.IVariableSupport {
	return &VariableSupport{set: make(map[string]interface{})}
}

func NewVariableSupport() *VariableSupport {
	return &VariableSupport{set: make(map[string]interface{})}
}

//---------------------------------------------

type VariableSupport struct {
	eventx.EventDispatcher
	set basis.VarSet
	mu  sync.RWMutex
}

func (s *VariableSupport) Vars() basis.VarSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set
}

func (s *VariableSupport) SetVar(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[key] = value
	s.DispatchEvent(basis.EventSetVariable, []interface{}{key, value})
}

func (s *VariableSupport) SetVars(kv basis.VarSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range kv {
		s.set[k] = v
	}
	s.DispatchEvent(basis.EventSetMultiVariable, kv)
}

func (s *VariableSupport) GetVar(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set[key]
}

func (s *VariableSupport) CheckVar(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[key]
	return ok
}

func (s *VariableSupport) RemoveVar(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.set, key)
}
