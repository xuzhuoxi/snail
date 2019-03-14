//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

import "github.com/xuzhuoxi/infra-go/eventx"

const (
	EventSetVariable      = "EventSetVariable"
	EventSetMultiVariable = "EventSetMultiVariable"
)

type VarSet map[string]interface{}

//变量列表
type IVariableSupport interface {
	eventx.IEventDispatcher
	SetVar(key string, value interface{})
	SetVars(kv VarSet)
	GetVar(key string) interface{}
	Vars() VarSet

	CheckVar(key string) bool
	RemoveVar(key string)
}
