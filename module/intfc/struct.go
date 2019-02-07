package intfc

import "github.com/xuzhuoxi/snail/conf"

type ModuleBase struct {
	cfg conf.ObjectConf
}

func (m *ModuleBase) GetModule() string {
	return m.cfg.Module
}

func (m *ModuleBase) GetName() string {
	return m.cfg.Name
}

func (m *ModuleBase) GetConfig() conf.ObjectConf {
	return m.cfg
}

func (m *ModuleBase) SetConfig(config conf.ObjectConf) {
	m.cfg = config
}

type GameServerState struct {
	//名称
	Name string

	//连接数
	LinkCount uint32
	//请求密度(次数/毫秒)
	ReqDensity uint32
	//响应密度(次数/毫秒)
	RespDensity uint32

	//最大响应时间
	MaxRT int64
}
