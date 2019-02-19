package root

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

type ModuleGame struct {
	imodule.ModuleBase

	singleCase   intfc.IGameSingleCase
	status       *GameStatus
	extensionCfg *ExtensionConfig
}

func (m *ModuleGame) Init() {
	singleCase := newGameSingleCase(m.Logger)
	m.singleCase = singleCase
	m.status = NewGameStatus(m.GetConfig(), singleCase)
	m.extensionCfg = NewExtensionConfig(singleCase)
	m.extensionCfg.ConfigExtensions()
}

func (m *ModuleGame) Run() {
	m.status.Start()
	m.extensionCfg.InitExtensions()

	go m.status.CheckRPC()
}

func (m *ModuleGame) Save() {
	panic("implement me")
}

func (m *ModuleGame) OnDestroy() {
}

func (m *ModuleGame) Destroy() {
}
