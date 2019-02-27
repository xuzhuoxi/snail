package root

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx/gobx"
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

type ModuleGame struct {
	imodule.ModuleBase

	singleCase intfc.IGameSingleCase
	status     *GameStatus
	server     *GameServer
}

func (m *ModuleGame) Init() {
	config := m.GetConfig()
	m.singleCase = m.newSingleCase()
	m.status = NewGameStatus(config, m.singleCase)
	m.server = NewGameServer(config, m.singleCase)
	m.server.InitServer()
}

func (m *ModuleGame) Run() {
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
	rs.OnceSetBuffEncoder(gobx.NewGobBuffEncoder(rs.DataBlockHandler()))
	rs.OnceSetBuffDecoder(gobx.NewGobBuffDecoder(rs.DataBlockHandler()))
	rs.OnceSetExtensionContainer(extension.NewSnailExtensionContainer())
	rs.OnceSetLogger(m.Logger)
	return rs
}
