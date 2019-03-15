//
//Created by xuzhuoxi
//on 2019-03-03.
//@author xuzhuoxi
//
package entity

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/lang"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIVariableSupport(currentTarget basis.IEntity) basis.IVariableSupport {
	return NewVariableSupport(currentTarget)
}

func NewVariableSupport(currentTarget basis.IEntity) *VariableSupport {
	return &VariableSupport{currentTarget: currentTarget, vars: make(map[string]interface{})}
}

//---------------------------------------------

type VariableSupport struct {
	currentTarget basis.IEntity
	eventx.EventDispatcher
	vars basis.VarSet
	mu   sync.RWMutex
}

func (s *VariableSupport) Vars() basis.VarSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vars
}

func (s *VariableSupport) SetVar(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if val, ok := s.vars[key]; ok && lang.Equal(val, value) {
		return
	}
	s.vars[key] = value
	kv := basis.NewVarSet()
	kv[key] = value
	s.DispatchEvent(basis.EventSetVariable, s.currentTarget, kv)
}

func (s *VariableSupport) SetVars(kv basis.VarSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var rm []string
	for k, v := range kv {
		if val, ok := s.vars[k]; ok && lang.Equal(val, v) {
			rm = append(rm, k)
			continue
		}
		s.vars[k] = v
	}
	if len(rm) > 0 { //去重
		for _, k := range rm {
			delete(kv, k)
		}
	}
	if len(kv) > 0 {
		s.DispatchEvent(basis.EventSetMultiVariable, s.currentTarget, kv)
	}
}

func (s *VariableSupport) GetVar(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.vars[key]
	return val, ok
}

func (s *VariableSupport) CheckVar(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.vars[key]
	return ok
}

func (s *VariableSupport) RemoveVar(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vars, key)
}
