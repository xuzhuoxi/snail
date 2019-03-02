//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

func NewGameExtensionSupport(Name string, SingleCase intfc.IGameSingleCase) GameExtensionSupport {
	support := protox.NewProtocolExtensionSupport(Name)
	return GameExtensionSupport{ProtocolExtensionSupport: support, SingleCase: SingleCase}
}

type GameExtensionSupport struct {
	protox.ProtocolExtensionSupport
	SingleCase intfc.IGameSingleCase
}

func (s *GameExtensionSupport) RequestDataType() protox.RequestDataType {
	return protox.StructValue
}

func (e *GameExtensionSupport) Logger() logx.ILogger {
	return e.SingleCase.Logger()
}
