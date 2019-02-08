package impl

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/util-go/encodingx"
	"github.com/xuzhuoxi/util-go/logx"
	"github.com/xuzhuoxi/util-go/netx"
)

type ModuleRoute struct {
	imodule.ModuleBase //内嵌
	gameMap            map[string]imodule.GameServerState
	//Service
	httpServer netx.IHttpServer
	rpcServer  netx.IRPCServer

	codecs *encodingx.GobCodecs
}

func (m *ModuleRoute) Init() {
	m.gameMap = make(map[string]imodule.GameServerState)
	m.codecs = encodingx.NewCodecs()
	m.initRPCServices()
	m.initForeignServices()
}

func (m *ModuleRoute) Run() {
}

func (m *ModuleRoute) Save() {
	panic("implement me")
}

func (m *ModuleRoute) OnDestroy() {
}

func (m *ModuleRoute) Destroy() {
}

//-----------------------------------------

func (m *ModuleRoute) initRPCServices() {
	rpcName := m.GetConfig().RpcList[0]
	if "" == rpcName {
		return
	}
	rpc, ok := conf.GetServiceConf(rpcName)
	if !ok {
		panic("RPC Undefined :" + rpcName)
	}
	m.rpcServer = netx.NewRPCServer()
	rpcHandler := new(imodule.RPCHandler)

	m.rpcServer.(logx.ILoggerSupport).SetLogger(m.Log)
	rpcHandler.Log = m.Log

	m.rpcServer.Register(rpcHandler)
	imodule.MapRPCHandler(rpcHandler, imodule.CmdRoute_OnConnected, m.onConnected)
	imodule.MapRPCHandler(rpcHandler, imodule.CmdRoute_OnDisconnected, m.onDisconnected)
	imodule.MapRPCHandler(rpcHandler, imodule.CmdRoute_UpdateState, m.onUpdateState)
	go func() {
		m.Log.Infoln(m.GetName(), ": start rpc server......")
		m.rpcServer.StartServer(rpc.Addr)
	}()
}

func (m *ModuleRoute) initForeignServices() {
	serviceName := m.GetConfig().ServiceList[0]
	service, ok := conf.GetServiceConf(serviceName)
	if !ok {
		panic("Service Undefined :" + serviceName)
	}
	m.httpServer = netx.NewHttpServer()
	go func() {
		m.Log.Infoln(m.GetName(), ": start http server......")
		m.httpServer.StartServer(service.Addr)
	}()
}

func (m *ModuleRoute) onConnected(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	name := args.From
	m.Log.Infoln(m.GetName(), ": onConnected:", name)
	return nil
}

func (m *ModuleRoute) onDisconnected(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	name := args.From
	m.Log.Infoln(m.GetName(), ": onDisconnected:", name)
	return nil
}

func (m *ModuleRoute) onUpdateState(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	var state imodule.GameServerState
	m.codecs.Decoder(args.Data, &state)
	m.gameMap[args.From] = state
	m.Log.Infoln(m.GetName(), ": onUpdateState:", state)
	return nil
}
