package root

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
	"time"
)

func NewGameStatus(config conf.ObjectConf, singleCase ifc.IGameSingleCase) *GameStatus {
	gameId := config.Id
	return &GameStatus{
		gameId:       gameId,
		config:       config,
		singleCase:   singleCase,
		state:        imodule.NewServiceState(gameId, imodule.DefaultStatsInterval),
		rpcRemoteMap: make(map[string]netx.IRPCClient)}
}

type GameStatus struct {
	gameId       string
	config       conf.ObjectConf
	singleCase   ifc.IGameSingleCase
	state        *imodule.ServiceStateDetail
	rpcRemoteMap map[string]netx.IRPCClient
}

func (s *GameStatus) GetPassTime() int64 {
	return s.state.GetPassNano() / int64(time.Second)
}

func (s *GameStatus) GetStatePriority() float64 {
	return s.state.StatsWeight()
}

func (s *GameStatus) DetailState() *imodule.ServiceStateDetail {
	return s.state
}

func (s *GameStatus) ToSimpleState() imodule.ServiceState {
	return imodule.ServiceState{Name: s.state.Name, Weight: s.state.StatsWeight()}
}

func (s *GameStatus) logger() logx.ILogger {
	return s.singleCase.Logger()
}

func (s *GameStatus) encoder() encodingx.IBuffEncoder {
	return s.singleCase.BuffEncoder()
}

//---------------------------------------------

func (s *GameStatus) Start() {
	s.state.Start()
	go s.CheckRPC()
}

func (s *GameStatus) CheckRPC() {
Conn:
	//fmt.Println(m.GetConfig().Name, ": start checkConn")
	s.checkAndConnRemotes()
	time.Sleep(time.Duration(s.state.StatsInterval()))
	goto Conn
}

func (s *GameStatus) checkAndConnRemotes() {
	remotes := s.config.Remotes
	for _, name := range remotes {
		s.checkAndConnRemote(name)
	}
}

func (s *GameStatus) checkAndConnRemote(toName string) {
	service, ok := s.config.GetServiceConf(toName)
	if !ok {
		s.logger().Fatalln(s.gameId, ": Remotes Error At:", toName)
		return
	}
	client, ok2 := s.rpcRemoteMap[toName]
	if !ok2 || !client.IsConnected() {
		s.conn2Service(toName, service.Network, service.Addr)
	}
}

func (s *GameStatus) conn2Service(toName string, network string, addr string) {
	client := netx.NewRPCClient(netx.RpcNetworkTCP)
	err := client.Dial(addr)
	if nil != err {
		return
	}
	s.logger().Infoln(s.gameId, "Connected to", toName, "(", addr, ") with RPC(", network, ")!")
	s.rpcRemoteMap[toName] = client
	s.notifyConnected(toName)
	s.notifyState(toName)
}

func (s *GameStatus) notifyRemotes(f func(to string)) {
	remotes := s.config.Remotes
	for _, remoteName := range remotes {
		client, ok := s.rpcRemoteMap[remoteName]
		if ok && client.IsConnected() {
			f(remoteName)
		}
	}
}

func (s *GameStatus) notifyConnected(toName string) {
	toClient := s.rpcRemoteMap[toName]
	config := s.config

	module := imodule.ModGame
	link, _ := config.GetServiceConf(config.ServiceList[0])
	state := imodule.ServiceState{Name: s.gameId, Weight: s.GetStatePriority()}

	//s.encoder().EncodeDataToBuff(module)
	//dataM := s.encoder().ReadBytes()
	//s.encoder().EncodeDataToBuff(link)
	//dataL := s.encoder().ReadBytes()
	//s.encoder().EncodeDataToBuff(state)
	//dataS := s.encoder().ReadBytes()
	//
	//s.encoder().EncodeDataToBuff(module)
	//dataM2 := s.encoder().ReadBytes()
	//s.encoder().EncodeDataToBuff(link)
	//dataL2 := s.encoder().ReadBytes()
	//s.encoder().EncodeDataToBuff(state)
	//dataS2 := s.encoder().ReadBytes()
	//
	//s.logger().Debugln(11, dataM)
	//s.logger().Debugln(12, dataM2)
	//s.logger().Debugln(21, dataL)
	//s.logger().Debugln(22, dataL2)
	//s.logger().Debugln(31, dataS)
	//s.logger().Debugln(32, dataS2)
	//args := &imodule.RPCArgs{From: s.gameId, Cmd: imodule.CmdRoute_OnConnected, Data: append(append(dataM, dataL...), dataS...)}

	s.encoder().EncodeDataToBuff(module, link, state)
	data := s.encoder().ReadBytes()

	args := &imodule.RPCArgs{From: s.gameId, Cmd: imodule.CmdRoute_OnConnected, Data: data}
	//s.logger().Debugln("GameStatus.Debug.notifyConnected:", *args)
	//s.logger().Debugln("GameStatus.Debug.notifyConnected1:", data)
	//s.logger().Debugln("GameStatus.Debug.notifyConnected2:", data2)

	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}

func (s *GameStatus) notifyDisConnected(toName string) {
	toClient := s.rpcRemoteMap[toName]
	args := &imodule.RPCArgs{From: s.gameId, Cmd: imodule.CmdRoute_OnDisconnected}
	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}

func (s *GameStatus) notifyState(toName string) {
	toClient := s.rpcRemoteMap[toName]

	state := s.ToSimpleState()
	s.encoder().EncodeDataToBuff(state)
	data := s.encoder().ReadBytes()

	args := &imodule.RPCArgs{From: s.gameId, Cmd: imodule.CmdRoute_UpdateState, Data: data}
	//s.logger().Debugln("GameStatus.Debug.notifyState:", args.Data, state)

	reply := &imodule.RPCReply{}
	toClient.Call(imodule.ServiceMethod_OnRPCCall, args, reply)
}
