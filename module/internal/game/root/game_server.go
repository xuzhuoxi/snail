//
//Created by xuzhuoxi
//on 2019-02-21.
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
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/engine/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
	"sync"
)

func NewGameServer(config conf.ObjectConf, singleCase ifc.IGameSingleCase) *GameServer {
	s := &GameServer{}
	s.config = config
	s.SingleCase = singleCase
	return s
}

type GameServer struct {
	config     conf.ObjectConf
	SingleCase ifc.IGameSingleCase

	Servers    []netx.ITCPServer
	Containers []extension.ISnailExtensionContainer
}

func (s *GameServer) InitServer() {
}

func (s *GameServer) StartServer() {
	for _, service := range s.config.ServiceList {
		conf, ok := s.config.GetServiceConf(service)
		if !ok {
			panic("Service[" + service + "] Undefined!")
		}
		s.startService(conf)
	}
	ifc.AddressProxy.AddEventListener(netx.EventAddressRemoved, s.onAddressRemap)
}

func (s *GameServer) StopServer() {
	ifc.AddressProxy.RemoveEventListener(netx.EventAddressRemoved, s.onAddressRemap)
	for index := len(s.Servers) - 1; index >= 0; index-- {
		s.Servers[index].RemoveEventListener(netx.ServerEventConnClosed, s.onConnClosed)
		s.Servers[index].StopServer()
	}
	s.Servers = nil
}

//--------------------------------------------------

//tcp
//json
func (s *GameServer) startService(conf conf.ServiceConf) {
	container := extension.NewISnailExtensionContainer()
	registerExtension(container, s.SingleCase)
	container.InitExtensions()
	s.Containers = append(s.Containers, container)

	server := netx.NewTCPServer()
	s.Servers = append(s.Servers, server)

	server.SetLinkMax(100)
	server.SetLogger(s.SingleCase.GetLogger())
	server.GetPackHandler().AppendPackHandler(newPackHandler(s.SingleCase, server, container).onPack)
	server.AddEventListener(netx.ServerEventConnClosed, s.onConnClosed)

	go server.StartServer(netx.SockParams{Network: conf.Network, LocalAddress: conf.Addr})
}

func (s *GameServer) onConnClosed(evd *eventx.EventData) {
	address := evd.Data.(string)
	ifc.AddressProxy.RemoveByAddress(address)
}

func (s *GameServer) onAddressRemap(evd *eventx.EventData) {
	address := evd.Data.(string)
	if "" == address {
		return
	}
	for _, server := range s.Servers {
		err, ok := server.CloseConnection(address)
		if ok {
			if nil != err {
				s.SingleCase.GetLogger().Warnln(err)
			}
			return
		}
	}
}

//---------------------------

func newPackHandler(singleCase ifc.IGameSingleCase, server netx.ISockServer, container extension.ISnailExtensionContainer) *packHandler {
	return &packHandler{
		singleCase: singleCase,
		server:     server,
		container:  container,
	}
}

type packHandler struct {
	singleCase ifc.IGameSingleCase
	container  extension.ISnailExtensionContainer
	server     netx.ISockServer
	syncLock   sync.Mutex
}

func (h *packHandler) GetLogger() logx.ILogger {
	return h.singleCase.GetLogger()
}

//消息处理入口，这里是并发方法
//msgData非共享的，但就是会发生并发问题，现在没搞明白
func (h *packHandler) onPack(msgData []byte, senderAddress string, other interface{}) bool {
	h.syncLock.Lock()
	defer h.syncLock.Unlock()
	name, pid, uid, data := h.parsePackMessage(msgData)
	extension := h.getProtocolExtension(name)
	if nil == extension {
		h.GetLogger().Warnln(fmt.Sprintf("Undefined Extension(%s)! Sender(%s)", name, uid))
		return false
	}
	if !extension.CheckProtocolId(pid) { //有效性检查
		h.GetLogger().Warnln(fmt.Sprintf("Undefined ProtoId(%s) Send to Extension(%s)! Sender(%s)", pid, name, uid))
		return false
	}
	if be, ok := extension.(protox.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(pid)
	}
	//请求处理
	response := &extendx.SockServerResponse{SockServer: h.server, Address: senderAddress, AddressProxy: ifc.AddressProxy}
	switch ne := extension.(type) {
	case protox.IOnNoneRequestExtension:
		ne.OnRequest(response, pid, uid)
	case protox.IOnBinaryRequestExtension:
		h.handleRequestBinary(response, ne, pid, uid, data)
	case protox.IOnObjectRequestExtension:
		h.handleRequestObject(response, ne, pid, uid, data)
	}
	if ae, ok := extension.(protox.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(pid)
	}
	return true
}

func (h *packHandler) handleRequestObject(response extendx.IExtensionResponse, extension protox.IOnObjectRequestExtension, pid string, uid string, data [][]byte) {
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

func (h *packHandler) handleRequestBinary(response extendx.IExtensionResponse, extension protox.IOnBinaryRequestExtension, pid string, uid string, data [][]byte) {
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
func (h *packHandler) parsePackMessage(msgBytes []byte) (name string, pid string, uid string, data [][]byte) {
	ifc.HandleBuffToData(func(buffToData bytex.IBuffToData) {
		buffToData.WriteBytes(msgBytes)
		name = string(buffToData.ReadData())
		pid = string(buffToData.ReadData())
		uid = string(buffToData.ReadData())
		if buffToData.Len() > 0 {
			for buffToData.Len() > 0 {
				d := buffToData.ReadData()
				//h.singleCase.GetLogger().Traceln("parsePackMessage", uid, d)
				if nil == d {
					//h.singleCase.GetLogger().Warnln("data is nil")
					break
				}
				data = append(data, d)
			}
		}
	})
	return name, pid, uid, data
}

func (h *packHandler) getProtocolExtension(pid string) ifc.IGameExtension {
	e := h.container.GetExtension(pid)
	if pe, ok := e.(ifc.IGameExtension); ok {
		return pe
	}
	return nil
}
