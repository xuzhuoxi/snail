package snail

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module"
)

const version = "1.0.0"

func init() {
	conf.DefaultConfig = conf.ParseConfig("config.json")
}

func Run(openGo bool) {
	startModule()
	startCmd(openGo)
}

func startModule() {
	logx.Infoln("snail.startModule..........")
	module.StartDefault()
	logx.Infoln("snail.startModule..........end")
}

func startCmd(openGo bool) {
	logx.Infoln("snail.startCmd..........")
	if openGo {
		go module.StartCmdListener()
	} else {
		module.StartCmdListener()
		logx.Infoln("snail.startCmd..........end")
	}
}
