package module

import (
	"github.com/xuzhuoxi/infra-go/cmdx"
	"github.com/xuzhuoxi/snail/engine"
	"github.com/xuzhuoxi/snail/module/config"
	"github.com/xuzhuoxi/snail/module/internal"
	_ "github.com/xuzhuoxi/snail/module/internal/admin"
	_ "github.com/xuzhuoxi/snail/module/internal/game"
	_ "github.com/xuzhuoxi/snail/module/internal/route"
)

//--------------------------------------------

func init() {
	config.DefaultModuleConfig = config.ParseModuleConfig(engine.GetDefaultFlagSet())
}

func StartDefault() {
	nameList := config.DefaultModuleConfig.OnList
	engine.SnailLogger.Infoln("module.Start..........", nameList)
	internal.Start(nameList...)
	engine.SnailLogger.Infoln("module.Start..........end")
}

func Start(nameList []string) {
	engine.SnailLogger.Infoln("module.Start..........", nameList)
	internal.Start(nameList...)
	engine.SnailLogger.Infoln("module.Start..........end")
}

func StopRunning() {
	engine.SnailLogger.Infoln("module.Stop..........")
	internal.StopAll()
	engine.SnailLogger.Infoln("module.Stop..........end")
}

func Stop(nameList []string) {
	engine.SnailLogger.Infoln("module.Stop..........", nameList)
	internal.Stop(nameList...)
	engine.SnailLogger.Infoln("module.Stop..........end")
}

func StartCmdListener() {
	cmdLine := cmdx.CreateCommandLineListener("请输入命令：", 0)
	cmdLine.MapCommand(internal.CommandKeyList, internal.CmdList)
	cmdLine.MapCommand(internal.CommandKeyInfo, internal.CmdInfo)
	cmdLine.MapCommand(internal.CommandKeyStart, internal.CmdStart)
	cmdLine.MapCommand(internal.CommandKeyStop, internal.CmdStop)

	cmdLine.MapCommand(internal.CommandKeyLogin, internal.CmdGameLogin)
	cmdLine.MapCommand(internal.CommandKeyLogout, internal.CmdGameLogout)

	cmdLine.StartListen() //这里会发生阻塞，保证程序不会结束
}
