//
//Created by xuzhuoxi
//on 2019-02-09.
//@author xuzhuoxi
//
package impl

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"time"
	"sync"
	"sort"
)

const Timeout = int64(time.Minute)

type server struct {
	Name   string
	Module imodule.Module
	Link   conf.ServiceConf
	State  imodule.ServiceState

	lastTimestamp int64
}

func (s *server) Timeout() bool {
	return (time.Now().UnixNano() - s.lastTimestamp) >= Timeout
}

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

type collection struct {
	servers serverlist
	maps    map[string]*server
	mu      sync.Mutex
}

func (c *collection) HasServer(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hasServer(name)
}

func (c *collection) UpdateServerState(state imodule.ServiceState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	name := state.Name
	if c.hasServer(name) {
		c.mu.Lock()
		defer c.mu.Unlock()
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
			rs = append(rs, c.removeServer(index).Name)
		}
	}
	return rs
}

func (c *collection) GetServer(name string) *server {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.maps[name]
	if ok {
		return copyServer(val)
	}
	return nil
}

func (c *collection) GetServers(module imodule.Module) []*server {
	c.mu.Lock()
	defer c.mu.Unlock()
	var rs []*server
	for _, server := range c.servers {
		if server.Module == module {
			rs = append(rs, copyServer(server))
		}
	}
	return rs
}

func (c *collection) hasServer(name string) bool {
	_, ok := c.maps[name]
	return ok
}

func (c *collection) addServer(server server) {
	r := &server
	c.servers = append(c.servers, r)
}

func (c *collection) removeServer(index int) *server {
	if index < 0 || index >= len(c.servers) {
		return nil
	}
	rs := c.servers[index]
	c.servers = append(c.servers[:index], c.servers[index+1:]...)
	delete(c.maps, rs.Name)
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

func newCollection() *collection {
	return &collection{servers: nil, maps: make(map[string]*server)}
}
