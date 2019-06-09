package imodule

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/module/config"
)

type IModule interface {
	IBaseModule

	Init()
	Run()
	Save()
	OnDestroy()
	Destroy()
}

type IBaseModule interface {
	GetId() string
	GetModuleName() string
	GetConfig() config.ObjectConf
	SetConfig(config config.ObjectConf)
	GetLogger() logx.ILogger
}

//ModuleBase-------------------------------------------

type ModuleBase struct {
	cfg    config.ObjectConf
	Logger logx.ILogger
}

func (m *ModuleBase) GetId() string {
	return m.cfg.Id
}

func (m *ModuleBase) GetModuleName() string {
	return m.cfg.ModuleName
}

func (m *ModuleBase) GetConfig() config.ObjectConf {
	return m.cfg
}

func (m *ModuleBase) SetConfig(config config.ObjectConf) {
	m.cfg = config
	m.updateLog()
}

func (m *ModuleBase) GetLogger() logx.ILogger {
	return m.Logger
}

func (m *ModuleBase) updateLog() {
	logName := m.cfg.Log
	if "" == logName {
		m.Logger = logx.DefaultLogger()
		return
	}
	fileDir, fileBaseName, fileExtName := m.cfg.LogFileInfo()
	newLog := logx.NewLogger()
	newLog.SetPrefix("[" + m.cfg.Id + "] ")
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeDailyFile, Level: logx.LevelAll, FileDir: fileDir, FileName: fileBaseName, FileExtName: "." + fileExtName})
	m.Logger = newLog
}

//ModuleName--------------------------------------

type ModuleConstructor func() IModule

type ModuleName string

const (
	ModRoute ModuleName = "route"
	ModGame  ModuleName = "game"
	ModAdmin ModuleName = "admin"
)

var moduleMap = make(map[ModuleName]ModuleConstructor)

func (m ModuleName) NewModule() IModule {
	if !m.Available() {
		panic("No Such ModuleName:" + m)
	}
	return moduleMap[m]()
}

func (m ModuleName) Available() bool {
	c, ok := moduleMap[m]
	return ok || nil != c
}

func RegisterModule(m ModuleName, constructor ModuleConstructor) {
	if m.Available() {
		panic("Repeat ModuleName Constructor:" + m)
	}
	moduleMap[m] = constructor
}
