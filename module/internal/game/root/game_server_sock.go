//
//Created by xuzhuoxi
//on 2019-03-26.
//@author xuzhuoxi
//
package root

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/extendx"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/timex"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
	"time"
)

func NewGameSock(conf conf.ServiceConf, single ifc.IGameSingleCase) *GameSock {
	container := extension.NewISnailExtensionContainer()
	registerExtension(container, single)
	container.InitExtensions()

	server := netx.NewTCPServer()
	server.SetName(conf.Name)
	server.SetMax(100)
	server.SetLogger(single.GetLogger())

	Status := imodule.NewServiceStateDetail(conf.Name, imodule.DefaultStatsInterval)

	return &GameSock{Conf: conf, Server: server, Container: container, StateDetail: Status}
}

type GameSock struct {
	Conf        conf.ServiceConf
	Server      netx.ITCPServer
	Container   ifc.IGameExtensionContainer
	StateDetail *imodule.ServiceStateDetail
}

func (gs *GameSock) Running() bool {
	if nil != gs.Server {
		return gs.Server.Running()
	} else {
		return false
	}
}

func (gs *GameSock) GetPassSecond() int64 {
	return gs.StateDetail.GetPassNano() / int64(time.Second)
}

func (gs *GameSock) GetStateDetail() imodule.IServiceStateDetail {
	return gs.StateDetail
}

func (gs *GameSock) GetStateSimple() imodule.ServiceState {
	return imodule.ServiceState{Name: gs.StateDetail.Name, Weight: gs.StateDetail.StatsWeight()}
}

//-------------------

func (gs *GameSock) SockRun() {
	gs.StateDetail.Start()
	gs.Server.AddEventListener(netx.ServerEventConnOpened, gs.onConnOpened)
	gs.Server.AddEventListener(netx.ServerEventConnClosed, gs.onConnClosed)
	gs.Server.GetPackHandler().AppendPackHandler(gs.onPack)
	gs.Server.StartServer(netx.SockParams{Network: gs.Conf.Network, LocalAddress: gs.Conf.Addr}) //这里会阻塞
}

func (gs *GameSock) SockStop() {
	gs.Server.RemoveEventListener(netx.ServerEventConnClosed, gs.onConnClosed)
	gs.Server.RemoveEventListener(netx.ServerEventConnOpened, gs.onConnOpened)
	gs.Server.StopServer()
	gs.Server.GetPackHandler().ClearHandlers()
}

//------------------

func (gs *GameSock) onConnOpened(evd *eventx.EventData) {
	gs.StateDetail.AddLinkCount()
}

func (gs *GameSock) onConnClosed(evd *eventx.EventData) {
	address := evd.Data.(string)
	ifc.AddressProxy.RemoveByAddress(address)
	gs.StateDetail.RemoveLinkCount()
}

//消息处理入口，这里是并发方法
//msgData非共享的，但在parsePackMessage后这部分数据会发生变化
func (gs *GameSock) onPack(msgData []byte, senderAddress string, other interface{}) bool {
	gs.StateDetail.AddReqCount()
	name, pid, uid, data := gs.parsePackMessage(msgData)
	extension := gs.getProtocolExtension(name)
	if nil == extension {
		ifc.LoggerExtension.Warnln(fmt.Sprintf("Undefined Extension(%s)! Sender(%s)", name, uid))
		return false
	}
	if !extension.CheckProtocolId(pid) { //有效性检查
		ifc.LoggerExtension.Warnln(fmt.Sprintf("Undefined ProtoId(%s) Send to Extension(%s)! Sender(%s)", pid, name, uid))
		return false
	}
	func() { //记录时间状态
		tn := time.Now().UnixNano()
		defer func() {
			un := time.Now().UnixNano() - tn
			ifc.LoggerExtension.Infoln(name, pid, un, timex.FormatUnixMilli(un/1e6, "5.999999ms")) //记录响应时间
			gs.StateDetail.AddRespUnixNano(un)
		}()
		gs.handleExtension(extension, senderAddress, name, pid, uid, data)
	}()
	return true
}

func (gs *GameSock) handleExtension(extension ifc.IGameExtension, senderAddress string, name string, pid string, uid string, data [][]byte) {
	if be, ok := extension.(protox.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(pid)
	}
	//请求处理
	response := &extendx.SockServerResponse{SockServer: gs.Server, Address: senderAddress, AddressProxy: ifc.AddressProxy}
	switch ne := extension.(type) {
	case protox.IOnNoneRequestExtension:
		ne.OnRequest(response, pid, uid)
	case protox.IOnBinaryRequestExtension:
		gs.handleRequestBinary(response, ne, pid, uid, data)
	case protox.IOnObjectRequestExtension:
		gs.handleRequestObject(response, ne, pid, uid, data)
	}
	if ae, ok := extension.(protox.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(pid)
	}
}

func (gs *GameSock) handleRequestObject(response extendx.IExtensionResponse, extension protox.IOnObjectRequestExtension, pid string, uid string, data [][]byte) {
	dataLn := len(data)
	if 0 == dataLn {
		extension.OnRequest(response, pid, uid, nil)
		return
	}
	var list []interface{}
	for _, bs := range data {
		newData := extension.GetRequestData(pid)
		ifc.HandleJsonCoding(func(codingHandler encodingx.ICodingHandler) {
			codingHandler.HandleDecode(bs, &newData)
		})
		list = append(list, newData)
	}
	if dataLn > 1 {
		if be, ok := extension.(protox.IBatchExtension); ok && be.Batch() {
			extension.OnRequest(response, pid, uid, list[0], list[1:]...)
		} else {
			for _, val := range list {
				extension.OnRequest(response, pid, uid, val)
			}
		}
	} else {
		extension.OnRequest(response, pid, uid, list[0])
	}
}

func (gs *GameSock) handleRequestBinary(response extendx.IExtensionResponse, extension protox.IOnBinaryRequestExtension, pid string, uid string, data [][]byte) {
	dataLn := len(data)
	if 0 == dataLn {
		extension.OnRequest(response, pid, uid, nil)
		return
	}
	if len(data) > 1 {
		if be, ok := extension.(protox.IBatchExtension); ok && be.Batch() {
			extension.OnRequest(response, pid, uid, data[0], data[1:]...)
		} else {
			for _, bs := range data {
				extension.OnRequest(response, pid, uid, bs)
			}
		}
	} else {
		extension.OnRequest(response, pid, uid, data[0])
	}
}

//block0 : eName utf8
//block1 : pid	utf8
//block2 : uid	utf8
//[n]其它信息
//这里为并发区域，但没有共享资源，可是就是出并发问题
func (gs *GameSock) parsePackMessage(msgBytes []byte) (name string, pid string, uid string, data [][]byte) {
	index := 0
	ifc.HandleBuffToData(func(buffToData bytex.IBuffToData) {
		buffToData.WriteBytes(msgBytes)
		name = string(buffToData.ReadData())
		pid = string(buffToData.ReadData())
		uid = string(buffToData.ReadData())
		if buffToData.Len() > 0 {
			for buffToData.Len() > 0 {
				n, d := buffToData.ReadDataTo(msgBytes[index:]) //由于msgBytes前部分数据已经处理完成，可以利用这部分空间
				//h.singleCase.GetLogger().Traceln("parsePackMessage", uid, d)
				if nil == d {
					//h.singleCase.GetLogger().Warnln("data is nil")
					break
				}
				data = append(data, d)
				index += n
			}
		}
	})
	return name, pid, uid, data
}

func (gs *GameSock) getProtocolExtension(pid string) ifc.IGameExtension {
	e := gs.Container.GetExtension(pid)
	if pe, ok := e.(ifc.IGameExtension); ok {
		return pe
	}
	return nil
}
