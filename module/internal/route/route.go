package route

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/route/impl"
)

func newModuleRoute() imodule.IModule {
	return &impl.ModuleRoute{}
}

func init() {
	imodule.RegisterModule(imodule.ModRoute, newModuleRoute)
}
