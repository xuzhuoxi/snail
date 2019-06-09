package impl

import (
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/module/config"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/route/ifc"
	"net/http"
	"sync"
	"time"
)

type ModuleRoute struct {
	imodule.ModuleBase //内嵌
	sockCollection     iSockCollection
	//Service
	httpServer netx.IHttpServer
	rpcServer  netx.IRPCServer

	mu sync.Mutex
}

func (m *ModuleRoute) Init() {
	m.sockCollection = newSockCollection()
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
		m.Logger.Infoln(":start rpc sock at:" + rpc.Addr)
		m.rpcServer.StartServer(rpc.Addr)
	}()
}

func (m *ModuleRoute) runForeignServices() {
	serviceName := m.GetConfig().SockList[0]
	service, ok := config.GetServiceConf(serviceName)
	if !ok {
		panic("Service Undefined :" + serviceName)
	}
	m.httpServer = netx.NewHttpServer()
	go func() {
		m.Logger.Infoln(":start http sock at:" + service.Addr)
		m.httpServer.MapFunc("/route", func(w http.ResponseWriter, r *http.Request) { m.onQueryRoute(w, r) })
		m.httpServer.StartServer(service.Addr)
	}()
}

func (m *ModuleRoute) onConnected(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	//m.Logger.Debugln("onConnected:", args.From, args.Data)
	var servers []sock
	ifc.HandleBuffDecodeFromPool(func(decoder encodingx.IBuffDecoder) {
		for _, bs := range args.Data {
			decoder.WriteBytes(bs)
		}
		var owner imodule.SockOwner
		decoder.DecodeDataFromBuff(&owner)
		for index := 1; index < len(args.Data); index++ {
			var cfg config.SockConf
			decoder.DecodeDataFromBuff(&cfg)
			server := &sock{SockOwner: owner, SockConf: cfg, SockState: imodule.SockState{SockName: cfg.Name}, lastTimestamp: time.Now().UnixNano()}
			servers = append(servers, *server)
			m.sockCollection.AddSock(server)
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
		if s := m.sockCollection.RemoveSock(serverName); nil != s {
			serverNames = append(serverNames, serverName)
		}
	}
	m.Logger.Infoln("ModuleRoute.onDisconnected:", serverNames)
	return nil
}

func (m *ModuleRoute) onUpdateState(args *imodule.RPCArgs, reply *imodule.RPCReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var state imodule.SockState
	ifc.HandleBuffDecodeFromPool(func(decoder encodingx.IBuffDecoder) {
		decoder.Reset()
		decoder.WriteBytes(args.Data[0])
		decoder.DecodeDataFromBuff(&state)
		m.sockCollection.UpdateSockState(state)
	})
	m.Logger.Infoln("ModuleRoute.onUpdateState:", args.From, state)
	return nil
}

func (m *ModuleRoute) onQueryRoute(w http.ResponseWriter, r *http.Request) {
	if sock, ok := m.sockCollection.GetWeightSock(); ok {
		bs, _ := jsoniter.Marshal(sock.SockConf)
		w.Write(bs)
		m.Logger.Infoln("ModuleRoute:onQueryRoute:", sock.SockConf, sock.SockConnections)
	} else {
		w.Write([]byte(""))
		m.Logger.Warnln("ModuleRoute:onQueryRoute:", "None Server")
	}
}
