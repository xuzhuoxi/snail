package imodule

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/osxu"
	"github.com/xuzhuoxi/snail/conf"
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
	GetConfig() conf.ObjectConf
	SetConfig(config conf.ObjectConf)
	GetLogger() logx.ILogger
}

//ModuleBase-------------------------------------------

type ModuleBase struct {
	cfg conf.ObjectConf
	Log logx.ILogger
}

func (m *ModuleBase) GetId() string {
	return m.cfg.Id
}

func (m *ModuleBase) GetModuleName() string {
	return m.cfg.ModuleName
}

func (m *ModuleBase) GetConfig() conf.ObjectConf {
	return m.cfg
}

func (m *ModuleBase) SetConfig(config conf.ObjectConf) {
	m.cfg = config
	m.updateLog()
}

func (m *ModuleBase) GetLogger() logx.ILogger {
	return m.Log
}

func (m *ModuleBase) updateLog() {
	logName := m.cfg.Log
	if "" == logName {
		m.Log = logx.DefaultLogger()
		return
	}
	fileName, extName := osxu.SplitFileName(logName)
	newLog := logx.NewLogger()
	newLog.SetPrefix("[" + m.cfg.Id + "] ")
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeDailyFile, Level: logx.LevelAll, FileDir: m.cfg.LogDir(), FileName: fileName, FileExtName: "." + extName})
	m.Log = newLog
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
