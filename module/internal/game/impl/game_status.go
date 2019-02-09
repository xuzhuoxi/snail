package impl

import "time"

func (m *ModuleGame) GetPassTime() int64 {
	return m.state.GetPassNano() / int64(time.Second)
}

func (m *ModuleGame) GetStatePriority() float64 {
	return m.state.StatsWeight()
}
