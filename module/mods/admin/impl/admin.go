package impl

import (
	"github.com/xuzhuoxi/snail/module/intfc"
)

type ModuleAdmin struct {
	intfc.ModuleBase
}

func (m *ModuleAdmin) Init() {

}

func (m *ModuleAdmin) Run() {

}

func (m *ModuleAdmin) Save() {
	panic("implement me")
}

func (m *ModuleAdmin) OnDestroy() {

}

func (m *ModuleAdmin) Destroy() {

}
