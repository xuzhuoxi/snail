//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/encodingx/gobx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func NewGameSingleCase() ifc.IGameSingleCase {
	return &GameSingleCase{}
}

type GameSingleCase struct {
	isInit             bool
	dataBlockHandler   bytex.IDataBlockHandler
	buffEncoder        encodingx.IBuffEncoder
	buffDecoder        encodingx.IBuffDecoder
	extensionContainer extension.ISnailExtensionContainer
	addressProxy       netx.IAddressProxy

	logger logx.ILogger
}

func (s *GameSingleCase) Init() {
	if s.isInit {
		return
	}
	s.isInit = true
	s.dataBlockHandler = bytex.NewDefaultDataBlockHandler()
	s.buffEncoder = gobx.NewGobBuffEncoder(s.dataBlockHandler)
	s.buffDecoder = gobx.NewGobBuffDecoder(s.dataBlockHandler)
	s.extensionContainer = extension.NewISnailExtensionContainer()
	s.addressProxy = netx.NewIAddressProxy()
}

func (s *GameSingleCase) SetLogger(logger logx.ILogger) {
	s.logger = logger
}

func (s *GameSingleCase) DataBlockHandler() bytex.IDataBlockHandler {
	return s.dataBlockHandler
}

func (s *GameSingleCase) BuffEncoder() encodingx.IBuffEncoder {
	return s.buffEncoder
}

func (s *GameSingleCase) BuffDecoder() encodingx.IBuffDecoder {
	return s.buffDecoder
}

func (s *GameSingleCase) ExtensionContainer() extension.ISnailExtensionContainer {
	return s.extensionContainer
}

func (s *GameSingleCase) AddressProxy() netx.IAddressProxy {
	return s.addressProxy
}

func (s *GameSingleCase) Logger() logx.ILogger {
	return s.logger
}
