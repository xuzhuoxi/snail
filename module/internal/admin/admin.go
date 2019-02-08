package admin

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/admin/impl"
)

func newModuleAdmin() imodule.IModule {
	return &impl.ModuleAdmin{}
}

func init() {
	imodule.RegisterModule(imodule.ModAdmin, newModuleAdmin)
}
