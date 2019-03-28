package root

import (
	"github.com/xuzhuoxi/infra-go/logx"
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
	m.initLoggerExtension()
	config := m.GetConfig()
	m.singleCase = m.newSingleCase()
	m.server = NewGameServer(config, m.singleCase)
	m.server.InitServer()

	m.status = NewGameStatus(config, m.singleCase, m.server)
}

func (m *ModuleGame) initLoggerExtension() {
	dir, baseName, extName := m.GetConfig().LogFileInfo()
	ifc.LoggerExtension.SetPrefix("[" + "RespTime" + "] ")
	//ifc.LoggerExtension.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	ifc.LoggerExtension.SetConfig(logx.LogConfig{Type: logx.TypeDailyFile, Level: logx.LevelInfo, FileDir: dir + "/extension/", FileName: baseName, FileExtName: "." + extName})
}

func (m *ModuleGame) Run() {
	m.status.StartNotify()
	m.server.StartServer()
}

func (m *ModuleGame) Save() {
}

func (m *ModuleGame) OnDestroy() {
}

func (m *ModuleGame) Destroy() {
	m.server.StopServer()
	m.status.StopNotify()
}

func (m *ModuleGame) newSingleCase() ifc.IGameSingleCase {
	rs := NewGameSingleCase()
	rs.Init()
	rs.SetLogger(m.Logger)
	return rs
}
