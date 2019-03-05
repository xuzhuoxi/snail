//
//Created by xuzhuoxi
//on 2019-03-03.
//@author xuzhuoxi
//
package world

import "sync"

type VarSet map[string]interface{}

//变量列表
type IVariableSupport interface {
	SetVar(key string, value interface{})
	GetVar(key string) interface{}

	CheckVar(key string) bool
	RemoveVar(key string)
}

func NewIVariableSupport() IVariableSupport {
	return &VariableSupport{set: make(map[string]interface{})}
}

func NewVariableSupport() *VariableSupport {
	return &VariableSupport{set: make(map[string]interface{})}
}

//---------------------------------------------

type VariableSupport struct {
	set VarSet
	mu  sync.RWMutex
}

func (s *VariableSupport) SetVar(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[key] = value
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
