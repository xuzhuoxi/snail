package impl

import "time"

func (m *ModuleGame) GetPassTime() int64 {
	return time.Now().UnixNano() - m.starting
}

func (m *ModuleGame) GetStatePriority() float64 {
	return 0
}

//吞吐量
func (m *ModuleGame) GetTPS() float64 {
	return 0
}

//最大响应时间
func (m *ModuleGame) GetMaxRT() int64 {
	return 0
}
