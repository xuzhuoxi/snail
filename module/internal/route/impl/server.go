//
//Created by xuzhuoxi
//on 2019-02-09.
//@author xuzhuoxi
//
package impl

import (
	"fmt"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"sort"
	"sync"
	"time"
)

const Timeout = int64(time.Minute)

func newSockCollection() iSockCollection {
	return &sockCollection{}
}

//----------------

type sock struct {
	ModuleId   string
	ModuleName imodule.ModuleName

	conf.SockConf
	imodule.SockState

	lastTimestamp int64
}

func (s *sock) SockName() string {
	return s.SockConf.Name
}

func (s *sock) Timeout() bool {
	return (time.Now().UnixNano() - s.lastTimestamp) >= Timeout
}

//-----------------------------

type sockList []*sock

func (sl sockList) Len() int {
	return len(sl)
}

func (sl sockList) Less(i, j int) bool {
	return sl[i].SockWeight < sl[j].SockWeight
}

func (sl sockList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

func (sl sockList) List() []*sock {
	if nil == sl {
		return nil
	}
	return sl
}

//-----------------------------

type iSockCollection interface {
	AddSock(s *sock)
	RemoveSock(id string) *sock
	UpdateSockState(state imodule.SockState)

	CheckSockByName(id string) bool
	GetSockByName(id string) *sock
	GetSocksByModuleId(id string) []*sock
	GetSocksByModule(moduleName imodule.ModuleName) []*sock

	ClearTimeout() []string
	GetWeightSock() *sock
}

type sockCollection struct {
	socks sockList
	mu    sync.RWMutex
}

func (c *sockCollection) AddSock(server *sock) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addSock(server)
}

func (c *sockCollection) RemoveSock(id string) *sock {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.removeSockById(id)
}

func (c *sockCollection) UpdateSockState(state imodule.SockState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s := c.getSock(state.SockName); s != nil {
		s.SockState = state
	}
}

func (c *sockCollection) CheckSockByName(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s := c.getSock(name); nil != s {
		return true
	}
	return false
}

func (c *sockCollection) GetSockByName(name string) *sock {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getSock(name)
}

func (c *sockCollection) GetSocksByModuleId(id string) []*sock {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []*sock
	for _, server := range c.socks {
		if server.ModuleId == id {
			rs = append(rs, server)
		}
	}
	return rs
}

func (c *sockCollection) GetSocksByModule(moduleName imodule.ModuleName) []*sock {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []*sock
	for _, server := range c.socks {
		if server.ModuleName == moduleName {
			rs = append(rs, server)
		}
	}
	return rs
}

func (c *sockCollection) GetWeightSock() *sock {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.clearTimeout()
	if c.socks.Len() <= 0 {
		return nil
	}
	sort.Sort(c.socks)
	return c.socks[0]
}

func (c *sockCollection) ClearTimeout() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.clearTimeout()
}

//---------------------------

func (c *sockCollection) getSock(name string) *sock {
	for _, ss := range c.socks {
		if ss.SockName() == name {
			return ss
		}
	}
	return nil
}

func (c *sockCollection) addSock(sock *sock) {
	if s := c.getSock(sock.SockName()); nil != s {
		return
	}
	c.socks = append(c.socks.List(), sock)
}

func (c *sockCollection) removeSockByIndex(index int) *sock {
	if index < 0 || index >= c.socks.Len() {
		return nil
	}
	list := c.socks.List()
	rs := list[index]
	c.socks = append(list[:index], list[index+1:]...)
	return rs
}

func (c *sockCollection) removeSockById(id string) *sock {
	for index, server := range c.socks {
		if id == server.SockName() {
			return c.removeSockByIndex(index)
		}
	}
	return nil
}

//
//func (c *sockCollection) copyServer(sock *sock) *sock {
//	copy := *sock
//	return &copy
//}

func (c *sockCollection) clearTimeout() []string {
	var rs []string
	for index := c.socks.Len() - 1; index >= 0; index-- {
		if c.socks[index].Timeout() {
			rs = append(rs, c.removeSockByIndex(index).SockName())
		}
	}
	return rs
}

func (c *sockCollection) printSocks() {
	for index, s := range c.socks {
		fmt.Println("服务器：", index, s.SockName(), s.SockState)
	}
}
