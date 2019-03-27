package impl

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/route/ifc"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ModuleRoute struct {
	imodule.ModuleBase //内嵌
	gameCollection     iCollection
	//Service
	httpServer netx.IHttpServer
	rpcServer  netx.IRPCServer

	mu sync.Mutex
}

func (m *ModuleRoute) Init() {
	m.gameCollection = newCollection()
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
	m.mu.Lock()
	defer m.mu.Unlock()
	//m.Logger.Debugln("onConnected:", args.From, args.Data)
	var servers []server
	ifc.HandleBuffDecode(func(decoder encodingx.IBuffDecoder) {
		for _, bs := range args.Data {
			decoder.WriteBytes(bs)
		}
		var module imodule.ModuleName
		decoder.DecodeDataFromBuff(&module)
		for index := 1; index < len(args.Data); index++ {
			var conf conf.ServiceConf
			var weight float64
			decoder.DecodeDataFromBuff(&conf, &weight)
			server := server{ModuleId: args.From, ModuleName: module, ServiceConf: conf, ServiceState: imodule.ServiceState{Name: conf.Name, Weight: weight}, lastTimestamp: time.Now().UnixNano()}
			servers = append(servers, server)
			m.gameCollection.AddServer(server)
		}
	})
	m.Logger.Infoln("ModuleRoute.onConnected:", args.From, servers)
	return nil
}

func (m *ModuleRoute) onDisconnected(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	//name := args.From
	data := args.Data
	var serverNames []string
	for _, d := range data {
		serverName := string(d)
		if s := m.gameCollection.RemoveServer(serverName); nil != s {
			serverNames = append(serverNames, serverName)
		}
	}
	m.Logger.Infoln("ModuleRoute.onDisconnected:", serverNames)
	return nil
}

func (m *ModuleRoute) onUpdateState(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var states []imodule.ServiceState
	ifc.HandleBuffDecode(func(decoder encodingx.IBuffDecoder) {
		for index := 0; index < len(args.Data); index++ {
			decoder.Reset()
			decoder.WriteBytes(args.Data[index])
			var state imodule.ServiceState
			decoder.DecodeDataFromBuff(&state)
			states = append(states, state)
			m.gameCollection.UpdateServerState(state)
		}
	})
	m.Logger.Infoln("ModuleRoute.onUpdateState:", args.From, states)
	return nil
}

func (m *ModuleRoute) onQueryRoute(w http.ResponseWriter, r *http.Request) {
	m.Logger.Infoln(":onQueryRoute:", len(m.gameCollection.GetServersByModule(imodule.ModGame)))
	w.Write([]byte(strconv.Itoa(len(m.gameCollection.GetServersByModule(imodule.ModGame)))))
}
