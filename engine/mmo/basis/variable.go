//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
)

const (
	EventVariableChanged = "EventVariableChanged"
)

func NewVarSet() encodingx.IKeyValue {
	return encodingx.NewCodingMap()
}

//变量列表
type IVariableSupport interface {
	eventx.IEventDispatcher
	SetVar(key string, value interface{})
	SetVars(kv encodingx.IKeyValue)
	GetVar(key string) (interface{}, bool)
	Vars() encodingx.IKeyValue

	CheckVar(key string) bool
	RemoveVar(key string)
}
