//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package ifc

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/engine/extension"
)

type IGameSingleCase interface {
	DataBlockHandler() bytex.IDataBlockHandler
	BuffEncoder() encodingx.IBuffEncoder
	BuffDecoder() encodingx.IBuffDecoder
	ExtensionContainer() extension.ISnailExtensionContainer
	AddressProxy() netx.IAddressProxy
	Logger() logx.ILogger

	Init()
	logx.ILoggerSetter
}
