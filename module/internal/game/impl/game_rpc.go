package impl

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/util-go/netx"
	"log"
	"time"
)

//private connect-------------------------------

func CheckRPC(m *ModuleGame) {
Conn:
	//fmt.Println(m.GetConfig().Name, ": start checkConn")
	checkAndConnRemotes(m)
	time.Sleep(10 * time.Second)
	goto Conn
}

func checkAndConnRemotes(m *ModuleGame) {
	for _, name := range m.GetConfig().ServiceList {
		checkAndConnRemote(m, name)
	}
}

func checkAndConnRemote(m *ModuleGame, name string) {
	service, ok := conf.GetServiceConf(name)
	if !ok {
		log.Fatalln(m.GetName(), ": Remotes Error At:", name)
		return
	}
	client, ok2 := m.remoteMap[name]
	if !ok2 || !client.IsConnected() {
		conn2Service(m, name, service.Network, service.Addr)
	}
}

func conn2Service(m *ModuleGame, name string, network string, addr string) {
	client := netx.NewRPCClient(netx.RpcNetworkTCP)
	err := client.Dial(addr)
	log.Println(m.GetName(), " Connecting to", name, "(", addr, ") with RPC(", network, ")!")
	if nil != err {
		return
	}
	m.remoteMap[name] = client
	notifyConnected(m, name)
	notifyState(m, name)
}

func notifyAllRemote(m *ModuleGame, f func(m *ModuleGame, to string)) {
	for _, remoteName := range m.GetConfig().ServiceList {
		client, ok := m.remoteMap[remoteName]
		if ok && client.IsConnected() {
			f(m, remoteName)
		}
	}
}

func notifyConnected(m *ModuleGame, to string) {
	toClient := m.remoteMap[to]
	args := &imodule.RPCArgs{From: m.GetName(), Cmd: imodule.CmdRoute_OnConnected}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}

func notifyDisConnected(m *ModuleGame, to string) {
	toClient := m.remoteMap[to]
	args := &imodule.RPCArgs{From: m.GetName(), Cmd: imodule.CmdRoute_OnDisconnected}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}

func notifyState(m *ModuleGame, to string) {
	toClient := m.remoteMap[to]
	data := m.codecs.Encode(m.state)
	args := &imodule.RPCArgs{From: m.GetName(), Cmd: imodule.CmdRoute_UpdateState, Data: data}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}
