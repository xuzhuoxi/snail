package imodule

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/util-go/logx"
	"github.com/xuzhuoxi/util-go/osxu"
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
	GetModule() string
	GetName() string
	GetConfig() conf.ObjectConf
	SetConfig(config conf.ObjectConf)
	GetLogger() logx.ILogger
}

//ModuleBase-------------------------------------------

type ModuleBase struct {
	cfg conf.ObjectConf
	Log logx.ILogger
}

func (m *ModuleBase) GetModule() string {
	return m.cfg.Module
}

func (m *ModuleBase) GetName() string {
	return m.cfg.Name
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
	newLog.SetPrefix("[" + m.cfg.Name + "] ")
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeDailyFile, Level: logx.LevelAll, FileDir: m.cfg.LogDir(), FileName: fileName, FileExtName: "." + extName})
	m.Log = newLog
}

//Module--------------------------------------

type ModuleConstructor func() IModule

type Module string

const (
	ModRoute Module = "route"
	ModGame  Module = "game"
	ModAdmin Module = "admin"
)

var moduleMap = make(map[Module]ModuleConstructor)

func (m Module) New() IModule {
	if !m.Available() {
		panic("No Such Module:" + m)
	}
	return moduleMap[m]()
}

func (m Module) Available() bool {
	c, ok := moduleMap[m]
	return ok || nil != c
}

func RegisterModule(m Module, constructor ModuleConstructor) {
	if m.Available() {
		panic("Repeat Module Constructor:" + m)
	}
	moduleMap[m] = constructor
}
