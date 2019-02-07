package game

import (
	"github.com/xuzhuoxi/snail/module/intfc"
	"github.com/xuzhuoxi/snail/module/mods/game/impl"
)

func NewModuleGame() intfc.IModule {
	return &impl.ModuleGame{}
}
