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
	"github.com/xuzhuoxi/snail/module/internal/route/ifc"
	"sort"
	"sync"
	"time"
)

func newSockCollection() iSockCollection {
	return &sockCollection{}
}

//----------------

type sock struct {
	imodule.SockOwner

	conf.SockConf
	imodule.SockState

	lastTimestamp int64
}

func (s *sock) SockName() string {
	return s.SockConf.Name
}

func (s *sock) IsTimeout() bool {
	return (time.Now().UnixNano() - s.lastTimestamp) >= ifc.SockTimeout
}

//-----------------------------

type sockWeightList []*sock

func (sl sockWeightList) Len() int {
	return len(sl)
}

func (sl sockWeightList) Less(i, j int) bool {
	bi := sl[i].IsTimeout()
	bj := sl[j].IsTimeout()
	if bi == bj {
		return sl[i].SockWeight < sl[j].SockWeight
	} else {
		return bj
	}
}

func (sl sockWeightList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

//-----------------------------

type sockLinkList []*sock

func (sl sockLinkList) Len() int {
	return len(sl)
}

func (sl sockLinkList) Less(i, j int) bool {
	return sl[i].SockConnections < sl[j].SockConnections
}

func (sl sockLinkList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

//-----------------------------

type iSockCollection interface {
	AddSock(s *sock)
	RemoveSock(id string) *sock
	UpdateSockState(state imodule.SockState)
	SockList() []sock

	CheckSockByName(name string) bool
	GetSockByName(name string) sock

	GetSocks(params imodule.SockOwner) []sock

	GetWeightSock() (rs sock, ok bool)
	GetLinkSock() (rs sock, ok bool)
}

type sockCollection struct {
	socks []*sock
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

func (c *sockCollection) SockList() []sock {
	if len(c.socks) == 0 {
		return nil
	}
	var rs []sock
	for _, s := range c.socks {
		rs = append(rs, *s)
	}
	return rs
}

func (c *sockCollection) CheckSockByName(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s := c.getSock(name); nil != s {
		return true
	}
	return false
}

func (c *sockCollection) GetSockByName(name string) sock {
	c.mu.Lock()
	defer c.mu.Unlock()
	return *c.getSock(name)
}

func (c *sockCollection) GetSocks(params imodule.SockOwner) []sock {
	var rs []sock
	for _, server := range c.socks {
		if server.IsTimeout() {
			continue
		}
		if params.PlatformId != "" && params.PlatformId != server.PlatformId {
			continue
		}
		if params.ModuleId != "" && params.ModuleId != server.ModuleId {
			continue
		}
		if params.ModuleName != "" && params.ModuleName != server.ModuleName {
			continue
		}
		rs = append(rs, *server)
	}
	return rs
}

func (c *sockCollection) GetWeightSock() (rs sock, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.socks) <= 0 {
		return sock{}, false
	}
	sort.Sort(sockWeightList(c.socks))
	return *c.socks[0], true
}

func (c *sockCollection) GetLinkSock() (rs sock, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.socks) <= 0 {
		return sock{}, false
	}
	sort.Sort(sockLinkList(c.socks))
	return *c.socks[0], true
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
	c.socks = append(c.socks, sock)
}

func (c *sockCollection) removeSockByIndex(index int) *sock {
	if index < 0 || index >= len(c.socks) {
		return nil
	}
	list := c.socks
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

//func (c *sockCollection) copyServer(sock *sock) *sock {
//	copy := *sock
//	return &copy
//}

func (c *sockCollection) printSocks() {
	for index, s := range c.socks {
		fmt.Println("服务器：", index, s.SockName(), s.SockState)
	}
}
