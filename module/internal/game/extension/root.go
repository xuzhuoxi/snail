//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

func NewGameExtensionSupport(ProtoId string, SingleCase intfc.IGameSingleCase) GameExtensionSupport {
	return GameExtensionSupport{ProtoId: ProtoId, SingleCase: SingleCase, SnailExtensionSupport: extension.SnailExtensionSupport{}}
}

type GameExtensionSupport struct {
	extension.SnailExtensionSupport

	ProtoId    string
	SingleCase intfc.IGameSingleCase
}

func (s *GameExtensionSupport) Key() string {
	return s.ProtoId
}

func (s *GameExtensionSupport) ProtocolId() string {
	return s.ProtoId
}
