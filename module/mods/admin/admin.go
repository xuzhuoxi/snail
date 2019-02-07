package admin

import (
	"github.com/xuzhuoxi/snail/module/intfc"
	"github.com/xuzhuoxi/snail/module/mods/admin/impl"
)

func NewModuleAdmin() intfc.IModule {
	return &impl.ModuleAdmin{}
}
