//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

func NewGameSingleCase() intfc.IGameSingleCase {
	return &GameSingleCase{}
}

type GameSingleCase struct {
	dataBlockHandler bytex.IDataBlockHandler
	buffEncoder      encodingx.IBuffEncoder
	buffDecoder      encodingx.IBuffDecoder
	container        extension.ISnailExtensionContainer
	logger           logx.ILogger
}

func (s *GameSingleCase) OnceSetBuffEncoder(encoder encodingx.IBuffEncoder) {
	if nil == s.buffEncoder {
		s.buffEncoder = encoder
	}
}

func (s *GameSingleCase) OnceSetBuffDecoder(decoder encodingx.IBuffDecoder) {
	if nil == s.buffDecoder {
		s.buffDecoder = decoder
	}
}

func (s *GameSingleCase) OnceSetDataBlockHandler(handler bytex.IDataBlockHandler) {
	if nil == s.dataBlockHandler {
		s.dataBlockHandler = handler
	}
}

func (s *GameSingleCase) OnceSetExtensionContainer(container extension.ISnailExtensionContainer) {
	if nil == s.container {
		s.container = container
	}
}

func (s *GameSingleCase) OnceSetLogger(logger logx.ILogger) {
	if nil == s.logger {
		s.logger = logger
	}
}

//--------------------------------

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
	return s.container
}

func (s *GameSingleCase) Logger() logx.ILogger {
	return s.logger
}
