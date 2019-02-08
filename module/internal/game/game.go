package game

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/impl"
)

func newModuleGame() imodule.IModule {
	return &impl.ModuleGame{}
}

func init() {
	imodule.RegisterModule(imodule.ModGame, newModuleGame)
}
