//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package intfc

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/engine/extension"
)

type IGameSingleCase interface {
	DataBlockHandler() bytex.IDataBlockHandler
	BuffEncoder() encodingx.IBuffEncoder
	BuffDecoder() encodingx.IBuffDecoder
	ExtensionContainer() extension.ISnailExtensionContainer
	Logger() logx.ILogger

	OnceSetDataBlockHandler(handler bytex.IDataBlockHandler)
	OnceSetBuffEncoder(encoder encodingx.IBuffEncoder)
	OnceSetBuffDecoder(decoder encodingx.IBuffDecoder)
	OnceSetExtensionContainer(container extension.ISnailExtensionContainer)
	OnceSetLogger(logger logx.ILogger)
}
