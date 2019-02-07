package impl

import (
	"github.com/xuzhuoxi/snail/module/intfc"
	"github.com/xuzhuoxi/util-go/encodingx"
	"github.com/xuzhuoxi/util-go/netx"
	"log"
)

type ModuleRoute struct {
	intfc.ModuleBase //内嵌
	gameMap          map[string]intfc.GameServerState
	//Service
	httpServer netx.IHttpServer
	rpcServer  netx.IRPCServer

	codecs *encodingx.GobCodecs
}

func (m *ModuleRoute) Init() {
	m.gameMap = make(map[string]intfc.GameServerState)
	m.codecs = encodingx.NewCodecs()
	m.initRPCServer()
	m.initHttpServer()
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

func (m *ModuleRoute) initRPCServer() {
	m.rpcServer = netx.NewRPCServer()
	rpcHandler := new(intfc.RPCHandler)
	m.rpcServer.Register(rpcHandler)
	intfc.MapRPCHandler(rpcHandler, intfc.CmdRoute_OnConnected, m.onConnected)
	intfc.MapRPCHandler(rpcHandler, intfc.CmdRoute_OnDisconnected, m.onDisconnected)
	intfc.MapRPCHandler(rpcHandler, intfc.CmdRoute_UpdateState, m.onUpdateState)
	rpc := m.GetConfig().GetRpcInfo()
	go func() {
		log.Println(m.GetName(), ": start rpc server......")
		m.rpcServer.StartServer(rpc.Addr)
	}()
}

func (m *ModuleRoute) initHttpServer() {
	m.httpServer = netx.NewHttpServer()
	service := m.GetConfig().Service
	go func() {
		log.Println(m.GetName(), ": start http server......")
		m.httpServer.StartServer(service.Addr)
	}()
}

func (m *ModuleRoute) onConnected(args *intfc.RPCArgs, reply *intfc.RPCReply) error {
	name := args.From
	log.Println(m.GetName(), ": onConnected:", name)
	return nil
}

func (m *ModuleRoute) onDisconnected(args *intfc.RPCArgs, reply *intfc.RPCReply) error {
	name := args.From
	log.Println(m.GetName(), ": onDisconnected:", name)
	return nil
}

func (m *ModuleRoute) onUpdateState(args *intfc.RPCArgs, reply *intfc.RPCReply) error {
	var state intfc.GameServerState
	m.codecs.Decoder(args.Data, &state)
	m.gameMap[args.From] = state
	log.Println(m.GetName(), ": onUpdateState:", state)
	return nil
}
