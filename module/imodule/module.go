package imodule

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/util-go/logx"
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

//--------------------------------------

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
