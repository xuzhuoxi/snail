//
//Created by xuzhuoxi
//on 2019-02-09.
//@author xuzhuoxi
//
package impl

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"sort"
	"sync"
	"time"
)

const Timeout = int64(time.Minute)

type server struct {
	Id         string
	ModuleName imodule.ModuleName
	Link       conf.ServiceConf
	State      imodule.ServiceState

	lastTimestamp int64
}

func (s *server) Timeout() bool {
	return (time.Now().UnixNano() - s.lastTimestamp) >= Timeout
}

//-----------------------------

type serverlist []*server

func (c serverlist) Len() int {
	return len(c)
}

func (c serverlist) Less(i, j int) bool {
	return c[i].State.Weight < c[j].State.Weight
}

func (c serverlist) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

//-----------------------------

type iCollection interface {
	InitServer(server)
	HasServer(id string) bool
	UpdateServerState(state imodule.ServiceState)
	ClearTimeout() []string
	GetServer(id string) *server
	GetServers(moduleName imodule.ModuleName) []*server
}

type collection struct {
	servers serverlist
	maps    map[string]*server
	mu      sync.Mutex
}

func (c *collection) InitServer(server server) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addServer(server)
}

func (c *collection) HasServer(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hasServer(id)
}

func (c *collection) UpdateServerState(state imodule.ServiceState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	name := state.Name
	if c.hasServer(name) {
		c.maps[name].State = state
		c.maps[name].lastTimestamp = time.Now().UnixNano()
	}
}

func (c *collection) ClearTimeout() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []string
	for index := len(c.servers) - 1; index >= 0; index-- {
		if c.servers[index].Timeout() {
			rs = append(rs, c.removeServer(index).Id)
		}
	}
	return rs
}

func (c *collection) GetServer(id string) *server {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.maps[id]
	if ok {
		return copyServer(val)
	}
	return nil
}

func (c *collection) GetServers(moduleName imodule.ModuleName) []*server {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []*server
	for _, server := range c.servers {
		if server.ModuleName == moduleName {
			rs = append(rs, copyServer(server))
		}
	}
	return rs
}

func (c *collection) hasServer(id string) bool {
	_, ok := c.maps[id]
	return ok
}

func (c *collection) addServer(server server) {
	if c.hasServer(server.Id) {
		return
	}
	r := &server
	c.servers = append(c.servers, r)
	c.maps[server.Id] = r
}

func (c *collection) removeServer(index int) *server {
	if index < 0 || index >= len(c.servers) {
		return nil
	}
	rs := c.servers[index]
	c.servers = append(c.servers[:index], c.servers[index+1:]...)
	delete(c.maps, rs.Id)
	return rs
}

func (c *collection) sortServers(servers serverlist) []*server {
	sort.Sort(servers)
	return servers
}

func copyServer(server *server) *server {
	copy := *server
	return &copy
}

func newCollection() iCollection {
	return &collection{servers: nil, maps: make(map[string]*server)}
}
