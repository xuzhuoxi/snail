package impl

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"net/http"
	"strconv"
	"time"
)

type ModuleRoute struct {
	imodule.ModuleBase //内嵌
	gameCollection     iCollection
	//Service
	httpServer netx.IHttpServer
	rpcServer  netx.IRPCServer

	gobBuffEncoder encodingx.IBuffEncoder
	gobBuffDecoder encodingx.IBuffDecoder
}

func (m *ModuleRoute) Init() {
	m.gameCollection = newCollection()
	m.gobBuffEncoder = encodingx.NewDefaultGobBuffEncoder()
	m.gobBuffDecoder = encodingx.NewDefaultGobBuffDecoder()
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
	m.rpcServer.(logx.ILoggerSupport).SetLogger(m.Logger)

	rpcHandler := imodule.NewRPCHandler(m.Logger)
	imodule.MapRPCFunction(rpcHandler, imodule.CmdRoute_OnConnected, m.onConnected)
	imodule.MapRPCFunction(rpcHandler, imodule.CmdRoute_OnDisconnected, m.onDisconnected)
	imodule.MapRPCFunction(rpcHandler, imodule.CmdRoute_UpdateState, m.onUpdateState)

	m.rpcServer.Register(rpcHandler)

	go func() {
		m.Logger.Infoln(":start rpc server at:" + rpc.Addr)
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
		m.Logger.Infoln(":start http server at:" + service.Addr)
		m.httpServer.MapFunc("/route", func(w http.ResponseWriter, r *http.Request) { m.onQueryRoute(w, r) })
		m.httpServer.StartServer(service.Addr)
	}()
}

func (m *ModuleRoute) onConnected(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	name := args.From
	//fmt.Println(222, args.Data)
	m.gobBuffDecoder.WriteBytes(args.Data)
	var module imodule.ModuleName
	var link conf.ServiceConf
	var state imodule.ServiceState
	m.gobBuffDecoder.DecodeFromBuff(&module, &link, &state)
	//fmt.Println(222, module, link, state)
	server := server{Id: state.Name, ModuleName: module, Link: link, State: state, lastTimestamp: time.Now().UnixNano()}
	m.gameCollection.InitServer(server)
	m.Logger.Infoln(":onConnected:", name, server)
	return nil
}

func (m *ModuleRoute) onDisconnected(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	name := args.From
	m.Logger.Infoln(":onDisconnected:", name)
	return nil
}

func (m *ModuleRoute) onUpdateState(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	m.gobBuffDecoder.WriteBytes(args.Data)
	var state imodule.ServiceState
	m.gobBuffDecoder.DecodeFromBuff(&state)
	m.gameCollection.UpdateServerState(state)
	m.Logger.Infoln(":onUpdateState:", state)
	return nil
}

func (m *ModuleRoute) onQueryRoute(w http.ResponseWriter, r *http.Request) {
	m.Logger.Infoln(":onQueryRoute:", len(m.gameCollection.GetServers(imodule.ModGame)))
	w.Write([]byte(strconv.Itoa(len(m.gameCollection.GetServers(imodule.ModGame)))))
}
