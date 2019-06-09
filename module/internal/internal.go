//
//Created by xuzhuoxi
//on 2019-02-08.
//@author xuzhuoxi
//
package internal

//这里不能import子模块如route,game,admin
//否则会循环引用
import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/snail/engine"
	"github.com/xuzhuoxi/snail/module/config"
	"github.com/xuzhuoxi/snail/module/imodule"
	"sync"
)

type state int

func (s state) String() string {
	if val, ok := stateDesc[s]; ok {
		return val
	} else {
		return "Unknown"
	}
}

const (
	StateDefault = 0 + iota
	StateInit
	StateRunning
	StateStopping
	StateDestroy
)

type internalMod struct {
	name  string
	mod   imodule.IModule
	state state
	index int
}

func (i *internalMod) running() bool {
	return i.state == StateRunning
}

func (i *internalMod) String() string {
	return fmt.Sprintf("{id(name):%s,\tmodule:%s,\tstate:%s,\trunning:%s}", i.name, i.mod.GetModuleName(), i.state.String(), fmt.Sprint(i.running()))
}

var (
	stateDesc = make(map[state]string)
	mods      []*internalMod
	modsMap   = make(map[string]*internalMod)
	mu        sync.Mutex
)

func init() {
	stateDesc[StateDefault] = "StateDefault"
	stateDesc[StateInit] = "StateInit"
	stateDesc[StateRunning] = "StateRunning"
	stateDesc[StateStopping] = "StateStopping"
	stateDesc[StateDestroy] = "StateDestroy"
}

func Start(name ...string) error {
	if len(name) == 0 {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	list := []*internalMod{}
	for _, n := range name {
		if Running(n) {
			engine.SnailLogger.Warnln("ModuleName " + n + "is running!")
			continue
		}
		internal, err := newInternal(n)
		if nil != err {
			return err
		}
		list = append(list, internal)
	}
	startModules(list...)
	return nil
}

func Stop(name ...string) {
	mu.Lock()
	defer mu.Unlock()
	var list []*internalMod
	for _, n := range name {
		if !Running(n) {
			engine.SnailLogger.Warnln("ModuleName " + n + "is not running!")
			continue
		}
		list = append(list, modsMap[n])
	}
	stopModules(list...)
}

func StopAll() {
	mu.Lock()
	defer mu.Unlock()
	var list []*internalMod
	for _, internal := range mods {
		if !internal.running() {
			continue
		}
		list = append(list, internal)
	}
	stopModules(list...)
}

func Running(name string) bool {
	internal, ok := modsMap[name]
	return ok && internal.running()
}

func ListInfo(state state) []string {
	var list []string
	for _, internal := range mods {
		if internal.state == state {
			list = append(list, internal.String())
		}
	}
	return list
}

//-------------------------------------

func newInternal(name string) (*internalMod, error) {
	c, has := config.GetObjectById(name)
	if !has {
		return nil, errors.New("No ModuleName Config :" + name)
	}
	m := imodule.ModuleName(c.ModuleName)
	if !m.Available() {
		return nil, errors.New("ModuleName Undefined:" + string(m))
	}
	rs := &internalMod{name: name, mod: m.NewModule()}
	rs.mod.SetConfig(c)
	return rs, nil
}

func cacheInternal(internal *internalMod) {
	if nil == internal {
		return
	}
	internal.index = len(mods)
	mods = append(mods, internal)
	modsMap[internal.name] = internal
}

func unCacheInternal(internal *internalMod) {
	if nil == internal {
		return
	}
	i, ok := modsMap[internal.name]
	if !ok {
		return
	}
	delete(modsMap, i.name)
	mods = append(mods[:i.index], mods[i.index+1:]...)
}

func initModule(m *internalMod) {
	if m.state != StateDefault && m.state != StateDestroy {
		return
	}
	engine.SnailLogger.Infoln(fmt.Sprintf("InitExtension..........[%s]", m.name))
	m.mod.Init()
	m.state = StateInit
}

func runModule(m *internalMod) {
	if m.state != StateInit {
		return
	}
	engine.SnailLogger.Infoln(fmt.Sprintf("Run...........[%s]", m.name))
	go m.mod.Run()
	m.state = StateRunning
}

func onDestroyModule(m *internalMod) {
	if m.state != StateRunning {
		return
	}
	m.mod.OnDestroy()
	m.state = StateStopping
	engine.SnailLogger.Infoln(fmt.Sprintf("OnDestroy.....[%s]", m.name))
}

func destroyModule(m *internalMod) {
	if m.state != StateStopping {
		return
	}
	m.mod.Destroy()
	m.state = StateDestroy
	engine.SnailLogger.Infoln(fmt.Sprintf("Destroy.......[%s]", m.name))
}

//---------------------------------------

func startModules(ms ...*internalMod) {
	if nil == ms || len(ms) == 0 {
		return
	}
	for _, mod := range ms {
		cacheInternal(mod)
	}
	for _, mod := range ms {
		initModule(mod)
	}
	for _, mod := range ms {
		runModule(mod)
	}
}

func stopModules(ms ...*internalMod) {
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
	for i := l - 1; i >= 0; i-- {
		unCacheInternal(ms[i])
	}
}

func foreach(mods []*internalMod, f func(i *internalMod) bool) []*internalMod {
	var rs []*internalMod
	for _, val := range mods {
		if f(val) {
			rs = append(rs, val)
		}
	}
	return rs
}
