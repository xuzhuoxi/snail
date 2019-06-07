package module

import (
	"github.com/xuzhuoxi/infra-go/cmdx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/engine/proxy"
	"github.com/xuzhuoxi/snail/module/internal"
	_ "github.com/xuzhuoxi/snail/module/internal/admin"
	_ "github.com/xuzhuoxi/snail/module/internal/game"
	_ "github.com/xuzhuoxi/snail/module/internal/route"
)

//--------------------------------------------

func StartDefault() {
	nameList := conf.DefaultConfig.OnList
	proxy.SnailLogger.Infoln("module.Start..........", nameList)
	internal.Start(nameList...)
	proxy.SnailLogger.Infoln("module.Start..........end")
}

func Start(nameList []string) {
	proxy.SnailLogger.Infoln("module.Start..........", nameList)
	internal.Start(nameList...)
	proxy.SnailLogger.Infoln("module.Start..........end")
}

func StopRunning() {
	proxy.SnailLogger.Infoln("module.Stop..........")
	internal.StopAll()
	proxy.SnailLogger.Infoln("module.Stop..........end")
}

func Stop(nameList []string) {
	proxy.SnailLogger.Infoln("module.Stop..........", nameList)
	internal.Stop(nameList...)
	proxy.SnailLogger.Infoln("module.Stop..........end")
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
