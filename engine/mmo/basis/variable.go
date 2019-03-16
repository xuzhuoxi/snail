//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/lang"
)

const (
	EventVariableChanged = "EventVariableChanged"
)

func NewVarSet() VarSet {
	return make(map[string]interface{})
}

type VarSet map[string]interface{}

func (v VarSet) Merge(vs VarSet) VarSet {
	var rm []string
	for key, val := range vs {
		if v2, ok := v[key]; ok && lang.Equal(v2, val) {
			rm = append(rm, key)
			continue
		}
		v[key] = val
	}
	if len(rm) > 0 { //有重复
		for _, key := range rm {
			delete(vs, key)
		}
	}
	if len(vs) == 0 {
		return nil
	}
	return vs
}
func (v VarSet) Set(key string, value interface{}) VarSet {
	if v2, ok := v[key]; ok && lang.Equal(v2, value) {
		return nil
	}
	v[key] = value
	rs := NewVarSet()
	rs[key] = value
	return rs
}

//变量列表
type IVariableSupport interface {
	eventx.IEventDispatcher
	SetVar(key string, value interface{})
	SetVars(kv VarSet)
	GetVar(key string) (interface{}, bool)
	Vars() VarSet

	CheckVar(key string) bool
	RemoveVar(key string)
}
