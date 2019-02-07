package module

import (
	"github.com/xuzhuoxi/snail/module/intfc"
	"github.com/xuzhuoxi/snail/module/mods/admin"
	"github.com/xuzhuoxi/snail/module/mods/game"
	"github.com/xuzhuoxi/snail/module/mods/route"
	"log"
	"strconv"
)

type internal struct {
	name    string
	mod     intfc.IModule
	running bool
}

func (i *internal) String() string {
	return i.name + "(running=" + strconv.FormatBool(i.running) + ",module=" + i.mod.GetConfig().Module + ")"
}

const (
	ModuleNameRoute = "route"
	ModuleNameGame  = "game"
	ModuleNameAdmin = "admin"
)

var (
	mods    []*internal
	modsMap = make(map[string]*internal)
)

func Register(m intfc.IModule) {
	mod := &internal{name: m.GetConfig().Name, mod: m}
	mods = append(mods, mod)
	modsMap[mod.name] = mod
}

func UnRegister(modName string) {
	_, ok := modsMap[modName]
	if ok {
		delete(modsMap, modName)
		for index, val := range mods {
			if val.name == modName {
				mods = append(mods[:index], mods[index+1:]...)
				return
			}
		}
	}
}

func Start(nameList []string) {
	log.Println("module.Start..........", nameList)
	var list []*internal
	for _, name := range nameList {
		i, ok := modsMap[name]
		if ok {
			list = append(list, i)
		}
	}
	startModules(list)
	log.Println("module.Start..........end")
}

func Stop() {
	log.Println("module.Stop..........")
	stopModules(mods)
	log.Println("module.Stop..........end")
}

func CreateModule(modName string) intfc.IModule {
	var rs intfc.IModule
	switch modName {
	case ModuleNameRoute:
		rs = route.NewModuleRoute()
	case ModuleNameAdmin:
		rs = admin.NewModuleAdmin()
	case ModuleNameGame:
		rs = game.NewModuleGame()
	}
	return rs
}

//private-----------------------------

func initModule(m *internal) {
	if m.running {
		return
	}
	log.Println("\t[" + m.name + "]:\tInit..........")
	m.mod.Init()
}

func runModule(m *internal) {
	if m.running {
		return
	}
	log.Println("\t[" + m.name + "]:\tRun..........")
	go m.mod.Run()
	m.running = true
}

func startModules(ms []*internal) {
	if nil == ms || len(ms) == 0 {
		return
	}
	for _, mod := range ms {
		initModule(mod)
	}
	for _, mod := range ms {
		runModule(mod)
	}
}

func onDestroyModule(m *internal) {
	if !m.running {
		return
	}
	log.Println("\t[" + m.name + "]:\tOnDestroy..........")
	m.mod.OnDestroy()
}

func destroyModule(m *internal) {
	if !m.running {
		return
	}
	log.Println("\t[" + m.name + "]:\tDestroy..........")
	m.mod.Destroy()
	m.running = false
}

func stopModules(ms []*internal) {
	if nil == ms || len(ms) == 0 {
		return
	}
	l := len(ms)
	for i := l - 1; i >= 0; i-- {
		onDestroyModule(ms[i])
	}
	for i := l - 1; i >= 0; i-- {
		destroyModule(ms[i])
	}
}

func foreach(mods []*internal, f func(i *internal) bool) []*internal {
	var rs []*internal
	for _, val := range mods {
		if f(val) {
			rs = append(rs, val)
		}
	}
	return rs
}
