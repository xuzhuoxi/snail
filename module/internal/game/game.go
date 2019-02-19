package game

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/root"
)

func newModuleGame() imodule.IModule {
	return &root.ModuleGame{}
}

func init() {
	imodule.RegisterModule(imodule.ModGame, newModuleGame)
}
