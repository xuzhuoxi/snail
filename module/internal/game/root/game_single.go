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

func newGameSingleCase(logger logx.ILogger) intfc.IGameSingleCase {
	return &GameSingleCase{
		dataBlockHandler:   bytex.NewDefaultDataBlockHandler(),
		logger:             logger,
		gobBuffEncoder:     encodingx.NewDefaultGobBuffEncoder(),
		gobBuffDecoder:     encodingx.NewDefaultGobBuffDecoder(),
		extensionContainer: extension.NewSnailExtensionContainer(),
	}
}

type GameSingleCase struct {
	dataBlockHandler   bytex.IDataBlockHandler
	logger             logx.ILogger
	gobBuffEncoder     encodingx.IGobBuffEncoder
	gobBuffDecoder     encodingx.IGobBuffDecoder
	extensionContainer extension.ISnailExtensionContainer
}

func (s *GameSingleCase) DataBlockHandler() bytex.IDataBlockHandler {
	return s.dataBlockHandler
}

func (s *GameSingleCase) GobBuffEncoder() encodingx.IGobBuffEncoder {
	return s.gobBuffEncoder
}

func (s *GameSingleCase) GobBuffDecoder() encodingx.IGobBuffDecoder {
	return s.gobBuffDecoder
}

func (s *GameSingleCase) ExtensionContainer() extension.ISnailExtensionContainer {
	return s.extensionContainer
}

func (s *GameSingleCase) Logger() logx.ILogger {
	return s.logger
}
