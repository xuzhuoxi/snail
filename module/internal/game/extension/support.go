//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func NewGameExtensionSupport(Name string) GameExtensionSupport {
	support := protox.NewProtoExtensionSupport(Name)
	return GameExtensionSupport{ProtoExtensionSupport: support}
}

type GameExtensionSupport struct {
	protox.ProtoExtensionSupport
	SingleCase ifc.IGameSingleCase
}

func (e *GameExtensionSupport) SetSingleCase(singleCase ifc.IGameSingleCase) {
	e.SingleCase = singleCase
}

func (e *GameExtensionSupport) GetLogger() logx.ILogger {
	return e.SingleCase.GetLogger()
}
