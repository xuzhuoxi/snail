package internal

import (
	"flag"
	"fmt"
	intfc2 "github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
	"strings"
)

const (
	CommandKeyList  = "list"
	CommandKeyInfo  = "info"
	CommandKeyStart = "start"
	CommandKeyStop  = "stop"

	CommandKeyLogin  = "login"
	CommandKeyLogout = "logout"
)

//list
//list -r=true(false) -m=ModuleName
func CmdList(cmdArgs []string) {
	flagMap, ok := parseCmdArgs(cmdArgs, 2, "r", "m")
	var list []*internalMod
	if ok {
		isRunning := flagMap["r"]
		moduleName := flagMap["m"]
		var eachFunc = func(i *internalMod) bool {
			m := moduleName == "" || moduleName == i.mod.GetConfig().ModuleName
			r := isRunning == "" || (isRunning == "false" || isRunning == "0") != i.running()
			return m && r
		}
		list = foreach(mods, eachFunc)
	} else {
		list = mods
	}
	fmt.Println(list)
}

//info -n=Name
func CmdInfo(cmdArgs []string) {
	flagMap, ok := parseCmdArgs(cmdArgs, 2, "g")
	if ok {
		name := flagMap["n"]
		if name == "" {
			fmt.Println("Command \"" + cmdArgs[0] + "\" args error!")
			return
		}
		mod, ok := modsMap[name]
		if !ok {
			fmt.Println("ModuleName \"" + name + "\" does not exist.")
			return
		}
		fmt.Println(mod)
	}
}

//stop -m=ModuleName
//stop -n=Name
func CmdStop(cmdArgs []string) {
	cmdSwitchModule(cmdArgs, false)
}

//start -m=ModuleName
//start -n=Name
func CmdStart(cmdArgs []string) {
	cmdSwitchModule(cmdArgs, true)
}

//login -g=Name
func CmdGameLogin(cmdArgs []string) {
	cmdLogin(cmdArgs, true)
}

//logout -g=Name
func CmdGameLogout(cmdArgs []string) {
	cmdLogin(cmdArgs, false)

}

//private ------------------------------------------

func parseCmdArgs(cmdArgs []string, minLen int, name ...string) (map[string]string, bool) {
	if len(cmdArgs) < minLen {
		fmt.Println("Command \"" + cmdArgs[0] + "\" args error!")
		return nil, false
	}
	if len(name) <= 0 {
		return nil, true
	}
	cmdFlag := flag.NewFlagSet(cmdArgs[0], flag.ContinueOnError)
	var values []*string
	rs := make(map[string]string)
	for _, key := range name {
		value := cmdFlag.String(key, "", "No Usage")
		values = append(values, value)
	}
	cmdFlag.Parse(cmdArgs[1:])
	for index, key := range name {
		rs[key] = *values[index]
	}
	return rs, true
}

func cmdSwitchModule(cmdArgs []string, on bool) {
	flagMap, ok := parseCmdArgs(cmdArgs, 2, "n", "m")
	if ok {
		name := flagMap["n"]
		module := flagMap["m"]
		if "" == name {
			var eachFunc = func(i *internalMod) bool {
				return on != i.running() && i.mod.GetModuleName() == module
			}
			list := foreach(mods, eachFunc)
			fmt.Println("list:", list)
			if on {
				startModules(list...)
			} else {
				stopModules(list...)
			}
		} else {
			mod, ok := modsMap[name]
			if !ok {
				fmt.Println("ModuleName \"" + name + "\" does not exist.")
				return
			}
			if mod.running() == on || (mod.mod.GetModuleName() != module && module != "") {
				fmt.Println(mod)
				return
			}
			if on {
				initModule(mod)
				runModule(mod)
			} else {
				onDestroyModule(mod)
				destroyModule(mod)
			}
		}
	}
}

//login -g=Name
func cmdLogin(cmdArgs []string, login bool) {
	flagMap, ok := parseCmdArgs(cmdArgs, 2, "g")
	if ok {
		gameName := flagMap["g"]
		mod, okm := modsMap[gameName]
		if !okm || mod.mod.GetModuleName() != string(intfc2.ModGame) {
			fmt.Println(strings.Replace("Game module \"${name}\" does not exist!", "${name}", gameName, -1))
			return
		}
		if !mod.running() {
			fmt.Println(strings.Replace("Game module \"${name}\" is not running!", "${name}", gameName, -1))
			return
		}
		gm, _ := mod.mod.(ifc.ILoginServer)
		if login {
			gm.Login()
		} else {
			gm.Logout()
		}
	}
}
