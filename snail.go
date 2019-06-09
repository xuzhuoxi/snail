package snail

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/module"
)

const version = "1.0.0"

func StartModule() {
	startModule()
}

func StartMMO() {
	startMMO()
}

func StartCmd(openGo bool) {
	startCmd(openGo)
}

func startModule() {
	logx.Infoln("snail.startModule..........")
	module.StartDefault()
	logx.Infoln("snail.startModule..........end")
}

func startMMO() {
	logx.Infoln("snail.startMMO..........")
	module.StartDefault()
	logx.Infoln("snail.startMMO..........end")
}

func startCmd(openGo bool) {
	logx.Infoln("snail.startCmd..........")
	if openGo {
		go module.StartCmdListener()
	} else {
		module.StartCmdListener()
	}
}
