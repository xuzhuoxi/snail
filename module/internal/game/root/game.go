package root

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

type ModuleGame struct {
	imodule.ModuleBase

	singleCase ifc.IGameSingleCase
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

func (m *ModuleGame) newSingleCase() ifc.IGameSingleCase {
	rs := NewGameSingleCase()
	rs.Init()
	rs.SetLogger(m.Logger)
	return rs
}
