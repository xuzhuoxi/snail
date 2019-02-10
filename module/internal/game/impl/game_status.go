package impl

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"time"
)

func (m *ModuleGame) GetPassTime() int64 {
	return m.state.GetPassNano() / int64(time.Second)
}

func (m *ModuleGame) GetStatePriority() float64 {
	return m.state.StatsWeight()
}

func (m *ModuleGame) ToSimpleState() imodule.ServiceState {
	return imodule.ServiceState{Name: m.state.Name, Weight: m.state.StatsWeight()}
}
