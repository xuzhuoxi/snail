//
//Created by xuzhuoxi
//on 2019-02-21.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
	"github.com/xuzhuoxi/snail/module/config"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func NewGameServer(config config.ObjectConf, singleCase ifc.IGameSingleCase) *GameServer {
	s := &GameServer{}
	s.config = config
	s.SingleCase = singleCase
	return s
}

type GameServer struct {
	eventx.EventDispatcher
	config     config.ObjectConf
	SingleCase ifc.IGameSingleCase
	GameSocks  []*GameSock
}

func (s *GameServer) InitServer() {
	for _, service := range s.config.SockList {
		conf, ok := s.config.GetServiceConf(service)
		if !ok {
			panic("Service[" + service + "] Undefined!")
		}
		s.GameSocks = append(s.GameSocks, NewGameSock(conf, s.SingleCase))
	}

}

func (s *GameServer) StartServer() {
	for _, gs := range s.GameSocks {
		gs.Server.AddEventListener(netx.ServerEventStart, s.onSockServerStart)
		gs.Server.AddEventListener(netx.ServerEventStop, s.onSockServerStop)
		go gs.SockRun()
	}
	ifc.AddressProxy.AddEventListener(netx.EventAddressRemoved, s.onAddressRemap)
}

func (s *GameServer) StopServer() {
	ifc.AddressProxy.RemoveEventListener(netx.EventAddressRemoved, s.onAddressRemap)
	for index := len(s.GameSocks) - 1; index >= 0; index-- {
		s.GameSocks[index].Server.RemoveEventListener(netx.ServerEventStop, s.onSockServerStop)
		s.GameSocks[index].Server.RemoveEventListener(netx.ServerEventStart, s.onSockServerStart)
		s.GameSocks[index].SockStop()
	}
	s.GameSocks = nil
}

func (s *GameServer) onSockServerStart(evd *eventx.EventData) {
	//fmt.Println(s.config.Id, "GameServer.onSockServerStart")
	server := evd.CurrentTarget.(tcpx.ITCPServer)
	gs, ok := s.getGameSock(server.GetName())
	if ok {
		s.DispatchEvent(netx.ServerEventStart, s, gs)
	}
}

func (s *GameServer) onSockServerStop(evd *eventx.EventData) {
	//fmt.Println("GameServer.onSockServerStop")
	server := evd.CurrentTarget.(tcpx.ITCPServer)
	s.DispatchEvent(netx.ServerEventStop, s, server.GetName())
}

func (s *GameServer) onAddressRemap(evd *eventx.EventData) {
	address := evd.Data.(string)
	if "" == address {
		return
	}
	for _, gameSock := range s.GameSocks {
		err, ok := gameSock.Server.CloseConnection(address)
		if ok {
			if nil != err {
				s.SingleCase.GetLogger().Warnln(err)
			}
			return
		}
	}
}

func (s *GameServer) getGameSock(name string) (*GameSock, bool) {
	for _, gs := range s.GameSocks {
		if name == gs.Server.GetName() {
			return gs, true
		}
	}
	return nil, false
}
