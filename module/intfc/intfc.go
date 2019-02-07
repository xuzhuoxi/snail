package intfc

import "github.com/xuzhuoxi/snail/conf"

type IModule interface {
	IModuleBase

	Init()
	Run()
	Save()
	OnDestroy()
	Destroy()
}

type IModuleBase interface {
	GetModule() string
	GetName() string
	GetConfig() conf.ObjectConf
	SetConfig(config conf.ObjectConf)
}
