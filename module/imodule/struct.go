package imodule

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/util-go/logx"
	"github.com/xuzhuoxi/util-go/osxu"
)

type ModuleBase struct {
	cfg conf.ObjectConf
	Log logx.ILogger
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
	m.updateLog()
}

func (m *ModuleBase) GetLogger() logx.ILogger {
	return m.Log
}

func (m *ModuleBase) updateLog() {
	logName := m.cfg.Log
	if "" == logName {
		m.Log = logx.DefaultLogger()
		return
	}
	fileName, extName := osxu.SplitFileName(logName)
	newLog := logx.NewLogger()
	newLog.SetPrefix("[" + m.cfg.Name + "] ")
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	newLog.SetConfig(logx.LogConfig{Type: logx.TypeDailyFile, Level: logx.LevelAll, FileDir: m.cfg.LogDir(), FileName: fileName, FileExtName: "." + extName})
	m.Log = newLog
}

//-----------------------------------------------------------

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
