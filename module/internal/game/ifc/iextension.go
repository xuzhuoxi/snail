//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package ifc

import (
	"github.com/xuzhuoxi/snail/engine/extension"
	"sync"
)

type IGameExtensionContainer = extension.ISnailExtensionContainer

type IGameExtension interface {
	extension.ISnailExtension
	IGameSingleCaseSetter
}

type GameExtensionConstructor func() IGameExtension

var (
	extConstructors []GameExtensionConstructor
	extMu           sync.RWMutex
)

func RegisterExtension(constructor GameExtensionConstructor) {
	extMu.Lock()
	defer extMu.Unlock()
	extConstructors = append(extConstructors, constructor)
}

func ForeachExtensionConstructor(eachFunc func(constructor GameExtensionConstructor)) {
	extMu.RLock()
	defer extMu.RUnlock()
	for _, c := range extConstructors {
		eachFunc(c)
	}
}
