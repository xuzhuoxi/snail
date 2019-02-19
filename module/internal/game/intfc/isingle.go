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
	GobBuffEncoder() encodingx.IGobBuffEncoder
	GobBuffDecoder() encodingx.IGobBuffDecoder
	ExtensionContainer() extension.ISnailExtensionContainer
	Logger() logx.ILogger
}
