//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/extension/demo"
	"github.com/xuzhuoxi/snail/module/internal/game/extension/user"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func registerExtension(container extension.ISnailExtensionContainer, single ifc.IGameSingleCase) {
	if nil == container || nil == single {
		return
	}
	funcAppend := func(container extension.ISnailExtensionContainer, extension extension.ISnailExtension) {
		extension.InitProtocolId()
		container.AppendExtension(extension)
	}
	funcAppend(container, demo.NewNoneDemoExtension("NoneDemo", single))
	funcAppend(container, demo.NewBinaryDemoExtension("BinaryDemo", single))
	funcAppend(container, demo.NewObjDemoExtension("ObjDemo", single))
	funcAppend(container, user.NewLoginExtension("Login", single))
}
