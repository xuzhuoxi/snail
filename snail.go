package snail

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module"
	utilCmd "github.com/xuzhuoxi/util-go/cmdx"
	"log"
)

const version = "1.0.0"

func init() {
	conf.Config = conf.ParseConfig("config.json")
}

func Run() {
	startModule()
	startCmd()
}

func startModule() {
	log.Println("snail.startModule..........")
	cfg := conf.Config
	modList := append(append(cfg.Games, cfg.Admins...), cfg.Routes...)
	for _, val := range modList {
		mod := module.CreateModule(val.Module)
		mod.SetConfig(val)
		module.Register(mod)
	}
	var nameOnList []string
	for _, name := range cfg.OnList {
		nameOnList = append(nameOnList, name)
	}
	module.Start(nameOnList)
	log.Println("snail.startModule..........end")
}

func startCmd() {
	log.Println("snail.startCmd..........")
	cmdLine := utilCmd.CreateCommandLineListener("请输入命令：", 0)
	cmdLine.MapCommand(module.CommandKeyList, module.CmdList)
	cmdLine.MapCommand(module.CommandKeyInfo, module.CmdInfo)
	cmdLine.MapCommand(module.CommandKeyStart, module.CmdStart)
	cmdLine.MapCommand(module.CommandKeyStop, module.CmdStop)

	cmdLine.MapCommand(module.CommandKeyLogin, module.CmdGameLogin)
	cmdLine.MapCommand(module.CommandKeyLogout, module.CmdGameLogout)

	cmdLine.StartListen() //这里会发生阻塞，保证程序不会结束
	log.Println("snail.startCmd..........end")
}
