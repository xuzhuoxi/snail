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
	ModuleId   string
	ModuleName imodule.ModuleName

	ServiceConf  conf.ServiceConf
	ServiceState imodule.ServiceState

	lastTimestamp int64
}

func (s *server) ServerId() string {
	return s.ServiceConf.Name
}

func (s *server) Timeout() bool {
	return (time.Now().UnixNano() - s.lastTimestamp) >= Timeout
}

//-----------------------------

type serverList []*server

func (c serverList) Len() int {
	return len(c)
}

func (c serverList) Less(i, j int) bool {
	return c[i].ServiceState.Weight < c[j].ServiceState.Weight
}

func (c serverList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

//-----------------------------

type iCollection interface {
	AddServer(server)
	RemoveServer(id string) *server
	UpdateServerState(state imodule.ServiceState)

	CheckServerById(id string) bool
	GetServerById(id string) *server
	GetServersByModuleId(id string) []*server
	GetServersByModule(moduleName imodule.ModuleName) []*server

	ClearTimeout() []string
}

type collection struct {
	servers serverList
	maps    map[string]*server
	mu      sync.Mutex
}

func (c *collection) AddServer(server server) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addServer(&server)
}

func (c *collection) RemoveServer(id string) *server {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.removeServerById(id)
}

func (c *collection) UpdateServerState(state imodule.ServiceState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	name := state.Name
	if c.hasServer(name) {
		c.maps[name].ServiceState = state
		c.maps[name].lastTimestamp = time.Now().UnixNano()
	}
}

func (c *collection) CheckServerById(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hasServer(id)
}

func (c *collection) GetServerById(id string) *server {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.maps[id]
	if ok {
		return copyServer(val)
	}
	return nil
}

func (c *collection) GetServersByModuleId(id string) []*server {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []*server
	for _, server := range c.servers {
		if server.ModuleId == id {
			rs = append(rs, copyServer(server))
		}
	}
	return rs
}

func (c *collection) GetServersByModule(moduleName imodule.ModuleName) []*server {
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

func (c *collection) ClearTimeout() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []string
	for index := len(c.servers) - 1; index >= 0; index-- {
		if c.servers[index].Timeout() {
			rs = append(rs, c.removeServerByIndex(index).ServerId())
		}
	}
	return rs
}

//---------------------------

func (c *collection) hasServer(id string) bool {
	_, ok := c.maps[id]
	return ok
}

func (c *collection) addServer(server *server) {
	if c.hasServer(server.ServerId()) {
		return
	}
	c.servers = append(c.servers, server)
	c.maps[server.ServerId()] = server
}

func (c *collection) removeServerById(id string) *server {
	for index := 0; index < c.servers.Len(); index-- {
		if id == c.servers[index].ServerId() {
			rs := c.servers[index]
			c.servers = append(c.servers[:index], c.servers[index+1:]...)
			delete(c.maps, id)
			return rs
		}
	}
	return nil
}

func (c *collection) removeServerByIndex(index int) *server {
	if index < 0 || index >= c.servers.Len() {
		return nil
	}
	rs := c.servers[index]
	c.servers = append(c.servers[:index], c.servers[index+1:]...)
	delete(c.maps, rs.ServerId())
	return rs
}

func (c *collection) sortServers(servers serverList) []*server {
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
