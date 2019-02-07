package route

import (
	"github.com/xuzhuoxi/snail/module/intfc"
	"github.com/xuzhuoxi/snail/module/mods/route/impl"
)

func NewModuleRoute() intfc.IModule {
	return &impl.ModuleRoute{}
}
