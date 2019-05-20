//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/snail/engine/extension"
	_ "github.com/xuzhuoxi/snail/module/internal/game/extension/demo"
	_ "github.com/xuzhuoxi/snail/module/internal/game/extension/user"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func injectExtensions(container extension.ISnailExtensionContainer, single ifc.IGameSingleCase) {
	if nil == container || nil == single {
		return
	}
	ifc.ForeachExtensionConstructor(func(constructor ifc.GameExtensionConstructor) {
		extension := constructor()
		extension.SetSingleCase(single)
		container.AppendExtension(extension)
	})
}
