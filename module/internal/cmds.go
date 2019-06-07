package internal

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/cmdx"
	intfc2 "github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

const (
	CommandKeyList  = "list"
	CommandKeyInfo  = "info"
	CommandKeyStart = "start"
	CommandKeyStop  = "stop"

	CommandKeyLogin  = "login"
	CommandKeyLogout = "logout"
)

func printlnMod(title string, mod *internalMod) {
	if len(title) > 0 {
		fmt.Println(title)
	}
	fmt.Println(mod.String())
}

func printlnMods(title string, mods []*internalMod) {
	if len(title) > 0 {
		fmt.Println(title)
	}
	for _, mod := range mods {
		fmt.Println(mod.String())
	}
}

//list
//list -r=true(false) -m=ModuleName
func CmdList(flagSet *cmdx.FlagSetExtend, args []string) {
	running := flagSet.Bool("r", false, "-r=true(false)")
	moduleName := flagSet.String("m", "", "-m=ModuleName")
	flagSet.Parse(args)
	rb := flagSet.CheckKey("r")
	mb := flagSet.CheckKey("m")
	list := foreach(mods, func(imod *internalMod) bool {
		if rb && *running != imod.running() {
			return false
		}
		if mb && *moduleName != imod.mod.GetModuleName() {
			return false
		}
		return true
	})
	printlnMods("", list)
}

//info -n=Name
func CmdInfo(flagSet *cmdx.FlagSetExtend, args []string) {
	name := flagSet.String("n", "", "-n=Name")
	flagSet.Parse(args)
	nb := flagSet.CheckKey("n")
	if !nb {
		fmt.Println("Command \"" + flagSet.Name() + "\" args error!")
		return
	}
	mod, ok := modsMap[*name]
	if !ok {
		fmt.Println("ModuleName \"" + *name + "\" does not exist.")
		return
	}
	printlnMod("", mod)
}

//stop -m=ModuleName
//stop -n=Name
func CmdStop(flagSet *cmdx.FlagSetExtend, args []string) {
	cmdSwitchModule(flagSet, args, false)
}

//start -m=ModuleName
//start -n=Name
func CmdStart(flagSet *cmdx.FlagSetExtend, args []string) {
	cmdSwitchModule(flagSet, args, true)
}

//login -g=Name
func CmdGameLogin(flagSet *cmdx.FlagSetExtend, args []string) {
	cmdLogin(flagSet, args, true)
}

//logout -g=Name
func CmdGameLogout(flagSet *cmdx.FlagSetExtend, args []string) {
	cmdLogin(flagSet, args, false)
}

//private ------------------------------------------

func cmdSwitchModule(flagSet *cmdx.FlagSetExtend, args []string, on bool) {
	name := flagSet.String("n", "", "-n=Name")
	moduleName := flagSet.String("m", "", "-m=ModuleName")
	flagSet.Parse(args)
	nb := flagSet.CheckKey("n")
	if nb {
		mod, ok := modsMap[*name]
		if !ok {
			fmt.Println("Name \"" + *name + "\" does not exist.")
			return
		}
		if mod.running() == on {
			printlnMod("", mod)
			return
		}
		if on {
			initModule(mod)
			runModule(mod)
		} else {
			onDestroyModule(mod)
			destroyModule(mod)
		}
		return
	}
	mb := flagSet.CheckKey("m")
	if !mb {
		fmt.Println("Command \"" + flagSet.Name() + "\" args error!")
		return
	}
	list := foreach(mods, func(i *internalMod) bool {
		if nb && *name != i.name {
			return false
		}
		if mb && *moduleName != i.mod.GetModuleName() {
			return false
		}
		if on != i.running() {
			return false
		}
		return true
	})
	printlnMods("list:", list)
	if on {
		startModules(list...)
	} else {
		stopModules(list...)
	}
}

//login -g=Name
func cmdLogin(flagSet *cmdx.FlagSetExtend, args []string, login bool) {
	gameName := flagSet.String("g", "", "-n=Name")
	flagSet.Parse(args)
	gb := flagSet.CheckKey("g")
	if !gb {
		fmt.Println("Command \"" + flagSet.Name() + "\" args error!")
		return
	}
	mod, ok := modsMap[*gameName]
	if !ok {
		fmt.Println(fmt.Sprintf("Game module \"%s\" does not exist!", *gameName))
		return
	}
	if mod.mod.GetModuleName() != string(intfc2.ModGame) {
		fmt.Println(fmt.Sprintf("Game \"%s\" is not game module!", *gameName))
		return
	}
	if !mod.running() {
		fmt.Println(fmt.Sprintf("Game module \"%s\" is not running!", *gameName))
		return
	}
	gm, _ := mod.mod.(ifc.ILoginServer)
	if login {
		gm.Login()
	} else {
		gm.Logout()
	}
}
