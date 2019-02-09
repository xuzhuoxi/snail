package impl

import (
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/util-go/encodingx"
	"github.com/xuzhuoxi/util-go/logx"
	"github.com/xuzhuoxi/util-go/netx"
	"net/http"
	"time"
)

type ModuleRoute struct {
	imodule.ModuleBase //内嵌
	gameMap            map[string]imodule.ServiceState
	//Service
	httpServer netx.IHttpServer
	rpcServer  netx.IRPCServer

	codecs *encodingx.GobCodecs
}

func (m *ModuleRoute) Init() {
	m.gameMap = make(map[string]imodule.ServiceState)
	m.codecs = encodingx.NewCodecs()

	time.Now().Unix()
}

func (m *ModuleRoute) Run() {
	m.runRPCServices()
	m.runForeignServices()
}

func (m *ModuleRoute) Save() {
	panic("implement me")
}

func (m *ModuleRoute) OnDestroy() {
	m.httpServer.StopServer()
	m.rpcServer.StopServer()
}

func (m *ModuleRoute) Destroy() {
}

//-----------------------------------------

func (m *ModuleRoute) runRPCServices() {
	objConf := m.GetConfig()
	rpcName := objConf.RpcList[0]
	if "" == rpcName {
		return
	}
	rpc, ok := objConf.GetServiceConf(rpcName)
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
		m.Log.Infoln(m.GetName(), ":start rpc server at:"+rpc.Addr)
		m.rpcServer.StartServer(rpc.Addr)
	}()
}

func (m *ModuleRoute) runForeignServices() {
	serviceName := m.GetConfig().ServiceList[0]
	service, ok := conf.GetServiceConf(serviceName)
	if !ok {
		panic("Service Undefined :" + serviceName)
	}
	m.httpServer = netx.NewHttpServer()
	go func() {
		m.Log.Infoln(m.GetName(), ":start http server at:"+service.Addr)
		m.httpServer.MapFunc("/route", m.onQueryRoute)
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
	var state imodule.ServiceState
	m.codecs.Decoder(args.Data, &state)
	m.gameMap[args.From] = state
	m.Log.Infoln(m.GetName(), ": onUpdateState:", state)
	return nil
}

//格式:
func (m *ModuleRoute) onQueryRoute(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("Route: " + tm))
}
