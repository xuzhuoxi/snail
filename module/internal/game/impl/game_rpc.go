package impl

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/util-go/netx"
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
	remotes := m.GetConfig().Remotes
	for _, name := range remotes {
		checkAndConnRemote(m, name)
	}
}

func checkAndConnRemote(m *ModuleGame, toName string) {
	service, ok := m.GetConfig().GetServiceConf(toName)
	if !ok {
		m.Log.Fatalln(m.GetId(), ": Remotes Error At:", toName)
		return
	}
	client, ok2 := m.rpcRemoteMap[toName]
	if !ok2 || !client.IsConnected() {
		conn2Service(m, toName, service.Network, service.Addr)
	}
}

func conn2Service(m *ModuleGame, toName string, network string, addr string) {
	client := netx.NewRPCClient(netx.RpcNetworkTCP)
	err := client.Dial(addr)
	m.Log.Infoln(m.GetId(), " Connecting to", toName, "(", addr, ") with RPC(", network, ")!")
	if nil != err {
		return
	}
	m.rpcRemoteMap[toName] = client
	notifyConnected(m, toName)
	notifyState(m, toName)
}

func notifyRemotes(m *ModuleGame, f func(m *ModuleGame, to string)) {
	remotes := m.GetConfig().Remotes
	for _, remoteName := range remotes {
		client, ok := m.rpcRemoteMap[remoteName]
		if ok && client.IsConnected() {
			f(m, remoteName)
		}
	}
}

func notifyConnected(m *ModuleGame, toName string) {
	toClient := m.rpcRemoteMap[toName]
	config := m.GetConfig()

	module := imodule.ModGame
	link, _ := config.GetServiceConf(config.ServiceList[0])
	state := imodule.ServiceState{Name: m.GetId(), Weight: m.GetStatePriority()}
	m.gobBuffEncoder.EncodeToBuff(module, link, state)
	data := m.gobBuffEncoder.EncodedBytes()

	args := &imodule.RPCArgs{From: m.GetId(), Cmd: imodule.CmdRoute_OnConnected, Data: data}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}

func notifyDisConnected(m *ModuleGame, toName string) {
	toClient := m.rpcRemoteMap[toName]
	args := &imodule.RPCArgs{From: m.GetId(), Cmd: imodule.CmdRoute_OnDisconnected}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}

func notifyState(m *ModuleGame, toName string) {
	toClient := m.rpcRemoteMap[toName]

	state := m.ToSimpleState()
	m.gobBuffEncoder.EncodeToBuff(state)
	data := m.gobBuffEncoder.EncodedBytes()

	args := &imodule.RPCArgs{From: m.GetId(), Cmd: imodule.CmdRoute_UpdateState, Data: data}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}
