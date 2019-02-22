package root

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

type ModuleGame struct {
	imodule.ModuleBase

	singleCase   intfc.IGameSingleCase
	status       *GameStatus
	server       *GameServer
	extensionCfg *ExtensionConfig
}

func (m *ModuleGame) Init() {
	config := m.GetConfig()
	m.singleCase = m.newSingleCase()
	m.status = NewGameStatus(config, m.singleCase)
	m.extensionCfg = NewExtensionConfig(m.singleCase)
	m.server = NewGameServer(config, m.singleCase)

	m.extensionCfg.ConfigExtensions()
}

func (m *ModuleGame) Run() {
	m.extensionCfg.InitExtensions()
	m.server.StartServer()
	m.status.Start()
}

func (m *ModuleGame) Save() {
}

func (m *ModuleGame) OnDestroy() {
}

func (m *ModuleGame) Destroy() {
}

func (m *ModuleGame) newSingleCase() intfc.IGameSingleCase {
	rs := NewGameSingleCase()
	rs.OnceSetDataBlockHandler(bytex.NewDefaultDataBlockHandler())
	rs.OnceSetBuffEncoder(encodingx.NewGobBuffEncoder(rs.DataBlockHandler()))
	rs.OnceSetBuffDecoder(encodingx.NewGobBuffDecoder(rs.DataBlockHandler()))
	rs.OnceSetExtensionContainer(extension.NewSnailExtensionContainer())
	rs.OnceSetLogger(m.Logger)
	return rs
}
