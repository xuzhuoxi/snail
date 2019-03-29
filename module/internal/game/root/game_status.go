package root

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
	"sync"
	"time"
)

func NewGameStatus(config conf.ObjectConf, singleCase ifc.IGameSingleCase, server *GameServer) *GameStatus {
	gameId := config.Id
	return &GameStatus{
		gameId:        gameId,
		config:        config,
		singleCase:    singleCase,
		server:        server,
		linkingServer: make(map[string]struct{}),
		rpcRemoteMap:  make(map[string]netx.IRPCClient)}
}

type GameStatus struct {
	gameId     string
	config     conf.ObjectConf
	singleCase ifc.IGameSingleCase

	server        *GameServer
	linkingServer map[string]struct{}
	rpcRemoteMap  map[string]netx.IRPCClient
	remoteMu      sync.RWMutex
}

func (s *GameStatus) logger() logx.ILogger {
	return s.singleCase.GetLogger()
}

//---------------------------------------------

func (s *GameStatus) StartNotify() {
	s.server.AddEventListener(netx.ServerEventStart, s.onGameSockStart)
	s.server.AddEventListener(netx.ServerEventStop, s.onGameSockStop)
	//go s.CheckRPC()
}

func (s *GameStatus) StopNotify() {
	s.server.RemoveEventListener(netx.ServerEventStop, s.onGameSockStop)
	s.server.RemoveEventListener(netx.ServerEventStart, s.onGameSockStart)
	//go s.CheckRPC()
}

func (s *GameStatus) onGameSockStart(evd *eventx.EventData) {
	//fmt.Println("GameStatus.onGameSockStart")
	gameSock := evd.Data.(*GameSock)

	s.remoteMu.Lock()
	s.linkingServer[gameSock.Server.GetName()] = struct{}{}
	s.remoteMu.Unlock()

	s.notifyConnected(gameSock)
}

func (s *GameStatus) onGameSockStop(evd *eventx.EventData) {
	//fmt.Println("GameStatus.onGameSockStop")
	gameSockName := evd.Data.(string)

	s.remoteMu.Lock()
	delete(s.linkingServer, gameSockName)
	s.remoteMu.Unlock()

	s.notifyDisConnected(gameSockName)
}

//-------------------------------

func (s *GameStatus) notifyConnected(gs *GameSock) {
	ifc.HandleBuffEncode(func(encoder encodingx.IBuffEncoder) {
		var data [][]byte //0:ModGame,[]

		owner := imodule.SockOwner{PlatformId: "", ModuleId: s.gameId, ModuleName: imodule.ModGame}
		encoder.EncodeDataToBuff(owner)
		data = append(data, encoder.ReadBytes()) //[0]

		sockConf := gs.Conf //conf.SockConf
		encoder.EncodeDataToBuff(sockConf)
		data = append(data, encoder.ReadBytes()) //[n]

		s.doNotifyRoutes(imodule.CmdRoute_OnConnected, data)
	})

	go func(gs *GameSock) {
	ReCheck:
		time.Sleep(ifc.GameNotifyRouteInterval)
		s.remoteMu.Lock()
		_, ok := s.linkingServer[gs.Server.GetName()]
		s.remoteMu.Unlock()
		if ok {
			s.notifyState(gs)
			goto ReCheck
		}
	}(gs)
}

func (s *GameStatus) notifyDisConnected(gameSockName string) {
	s.doNotifyRoutes(imodule.CmdRoute_OnDisconnected, [][]byte{[]byte(gameSockName)})
}

func (s *GameStatus) notifyState(gs *GameSock) {
	ifc.HandleBuffEncode(func(encoder encodingx.IBuffEncoder) {
		var data [][]byte
		state := gs.GetSockState()
		encoder.EncodeDataToBuff(state)
		data = append(data, encoder.ReadBytes())
		s.doNotifyRoutes(imodule.CmdRoute_UpdateState, data)
	})
}

//-------------------------------

func (s *GameStatus) doNotifyRoutes(Cmd string, data [][]byte) {
	//fmt.Println("GameStatus.doNotifyRoutes", Cmd, data)
	remotes := s.config.Remotes
	for _, remoteName := range remotes {
		client, ok, _ := s.getRemoteClient(remoteName)
		if !ok || !client.IsConnected() {
			continue
		}
		s.doNotifyRoute(remoteName, Cmd, data)
	}
}

func (s *GameStatus) doNotifyRoute(remoteName string, Cmd string, data [][]byte) {
	client := s.rpcRemoteMap[remoteName]
	args := &imodule.RPCArgs{From: s.gameId, Cmd: Cmd, Data: data}
	reply := &imodule.RPCReply{}
	err := client.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
	if nil != err {
		s.remoteMu.Lock()
		defer s.remoteMu.Unlock()
		s.cacheRemoteClient(remoteName, nil)
	}
}

//-----------------------------

func (s *GameStatus) getRemoteClient(remoteName string) (client netx.IRPCClient, ok bool, isNew bool) {
	s.remoteMu.RLock()
	defer s.remoteMu.RUnlock()
	client, ok = s.rpcRemoteMap[remoteName]
	if ok {
		return
	}
	service, ok2 := s.config.GetServiceConf(remoteName)
	if !ok2 {
		s.logger().Fatalln(s.gameId, ": Remotes Error At:", remoteName)
		return nil, false, false
	}
	client = netx.NewRPCClient(netx.RpcNetworkTCP)
	err := client.Dial(service.Addr)
	if nil != err {
		return nil, false, false
	}
	s.logger().Infoln(fmt.Sprintf("Connected to %s(%s) with RPC(%s)", remoteName, service.Addr, service.Network))
	s.cacheRemoteClient(remoteName, client)
	return client, true, true
}

func (s *GameStatus) cacheRemoteClient(remoteName string, client netx.IRPCClient) {
	if nil == client {
		delete(s.rpcRemoteMap, remoteName)
	} else {
		s.rpcRemoteMap[remoteName] = client
	}
}
