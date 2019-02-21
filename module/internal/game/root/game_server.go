//
//Created by xuzhuoxi
//on 2019-02-21.
//@author xuzhuoxi
//
package root

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

func NewGameServer(config conf.ObjectConf, singleCase intfc.IGameSingleCase) *GameServer {
	s := &GameServer{}
	s.config = config
	s.PackHandler = s.onPack
	s.SingleCase = singleCase
	return s
}

type GameServer struct {
	config      conf.ObjectConf
	SingleCase  intfc.IGameSingleCase
	PackHandler netx.PackHandler

	Server []netx.ISockServer
}

func (s *GameServer) StartServer() {
	for _, service := range s.config.ServiceList {
		conf, ok := s.config.GetServiceConf(service)
		if !ok {
			panic("Service[" + service + "] Undefined!")
		}
		server := netx.NewTCPServer(100)
		server.SetLogger(s.SingleCase.Logger())
		s.Server = append(s.Server, server)
		go server.StartServer(netx.SockParams{Network: conf.Network, LocalAddress: conf.Addr})
	}
}

func (s *GameServer) StopServer() {
	for index := len(s.Server) - 1; index >= 0; index-- {
		s.Server[index].StopServer()
	}
	s.Server = nil
}

func (s *GameServer) onPack(msgBytes []byte, info interface{}) {
	fmt.Println(11111)
}
