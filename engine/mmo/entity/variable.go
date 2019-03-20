//
//Created by xuzhuoxi
//on 2019-03-03.
//@author xuzhuoxi
//
package entity

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIVariableSupport(currentTarget basis.IEntity) basis.IVariableSupport {
	return NewVariableSupport(currentTarget)
}

func NewVariableSupport(currentTarget basis.IEntity) *VariableSupport {
	return &VariableSupport{currentTarget: currentTarget, vars: basis.NewVarSet()}
}

//---------------------------------------------

type VariableSupport struct {
	currentTarget basis.IEntity
	eventx.EventDispatcher
	vars encodingx.IKeyValue
	mu   sync.RWMutex
}

func (s *VariableSupport) Vars() encodingx.IKeyValue {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vars
}

func (s *VariableSupport) SetVar(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if diff, ok := s.vars.Set(key, value); ok {
		s.DispatchEvent(basis.EventVariableChanged, s.currentTarget, diff)
	}
}

func (s *VariableSupport) SetVars(kv encodingx.IKeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	diff := s.vars.Merge(kv)
	if nil != diff {
		s.DispatchEvent(basis.EventVariableChanged, s.currentTarget, diff)
	}
}

func (s *VariableSupport) GetVar(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vars.Get(key)
}

func (s *VariableSupport) CheckVar(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vars.Check(key)
}

func (s *VariableSupport) RemoveVar(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vars.Delete(key)
}
