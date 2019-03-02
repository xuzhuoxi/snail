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
	"github.com/xuzhuoxi/infra-go/encodingx/jsonx"
	"github.com/xuzhuoxi/infra-go/extendx"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/conf"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
	"sync"
)

func NewGameServer(config conf.ObjectConf, singleCase intfc.IGameSingleCase) *GameServer {
	s := &GameServer{}
	s.config = config
	s.SingleCase = singleCase
	s.BuffToData = bytex.NewBuffToData(bytex.NewDefaultDataBlockHandler())
	return s
}

type GameServer struct {
	config     conf.ObjectConf
	SingleCase intfc.IGameSingleCase

	Server []netx.ISockServer

	BuffToData bytex.IBuffToData
	buffMu     sync.Mutex

	extensionCfg *ExtensionConfig

	index int
}

func (s *GameServer) InitServer() {
	s.extensionCfg = NewExtensionConfig(s.SingleCase)
	s.extensionCfg.ConfigExtensions()
}

func (s *GameServer) StartServer() {
	s.extensionCfg.InitExtensions()
	for _, service := range s.config.ServiceList {
		conf, ok := s.config.GetServiceConf(service)
		if !ok {
			panic("Service[" + service + "] Undefined!")
		}
		server := netx.NewTCPServer(100)
		server.SetLogger(s.SingleCase.Logger())
		server.SetPackHandler(newPackHandler(s.SingleCase, *s.extensionCfg).onPack)
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

//--------------------------------------------------

func newPackHandler(singleCase intfc.IGameSingleCase, extensionCfg ExtensionConfig) *packHandler {
	return &packHandler{
		buffToData:   bytex.NewBuffToData(bytex.NewDefaultDataBlockHandler()),
		extensionCfg: extensionCfg,
		singleCase:   singleCase,
		decoder:      jsonx.NewDefaultJsonCodingHandler(),
	}
}

type packHandler struct {
	buffToData   bytex.IBuffToData
	extensionCfg ExtensionConfig //克隆，减少资源竞争
	singleCase   intfc.IGameSingleCase
	decoder      encodingx.IDecodeHandler

	index int
}

func (h *packHandler) onPack(msgBytes []byte, info interface{}) {
	pid, uid, data := h.parsePackMessage(msgBytes)
	extension := h.getExtension(pid)
	if be, ok := extension.(protox.IBeforeRequestExtension); ok {
		be.BeforeRequest()
	}
	if re, ok := extension.(protox.IOnRequestExtension); ok {
		dataType := re.RequestDataType()
		switch {
		case dataType == protox.None || len(data) == 0:
			h.handleRequestNone(re, pid, uid)
		case dataType == protox.ByteArray:
			h.handleRequestByteArray(re, pid, uid, data)
		case dataType == protox.StructValue:
			h.handleRequestStructValue(re, pid, uid, data)
		}
	}
	if ae, ok := extension.(protox.IAfterRequestExtension); ok {
		ae.AfterRequest()
	}
	fmt.Println(h.index, pid, uid, data)
	h.index++
}

func (h *packHandler) handleRequestStructValue(extension protox.IOnRequestExtension, pid string, uid string, data [][]byte) {
	var list []interface{}
	for _, bs := range data {
		newData := extension.RequestData()
		h.decoder.HandleDecode(bs, &newData)
		list = append(list, newData)
	}
	if len(list) > 1 {
		if be, ok := extension.(protox.IBatchExtension); ok {
			if be.Batch() {
				extension.OnRequest(pid, uid, list[0], list[1:]...)
				return
			}
		}
		for _, val := range list {
			extension.OnRequest(pid, uid, val)
		}
	} else {
		extension.OnRequest(pid, uid, list[0])
	}
}

func (h *packHandler) handleRequestByteArray(extension protox.IOnRequestExtension, pid string, uid string, data [][]byte) {
	if len(data) > 1 {
		if be, ok := extension.(protox.IBatchExtension); ok {
			if be.Batch() {
				data2 := []interface{}{}
				for index := 1; index < len(data); index++ {
					data2 = append(data2, data[index])
				}
				extension.OnRequest(pid, uid, data[0], data2...)
				return
			}
		}
		for _, bs := range data {
			extension.OnRequest(pid, uid, bs)
		}
	} else {
		extension.OnRequest(pid, uid, data[0])
	}
}

func (h *packHandler) handleRequestNone(extension protox.IOnRequestExtension, pid string, uid string) {
	extension.OnRequest(pid, uid, nil)
}

//block0 : pid	utf8
//block1 : uid	utf8
//[n]其它信息
func (h *packHandler) parsePackMessage(msgBytes []byte) (pid string, uid string, data [][]byte) {
	h.buffToData.Reset()
	h.buffToData.WriteBytes(msgBytes)
	pid = string(h.buffToData.ReadData())
	uid = string(h.buffToData.ReadData())
	if h.buffToData.Len() > 0 {
		for h.buffToData.Len() > 0 {
			d := h.buffToData.ReadData()
			if nil == d {
				h.singleCase.Logger().Warnln("data is nil")
				break
			}
			data = append(data, d)
		}
	}
	return
}

func (h *packHandler) getExtension(pid string) extendx.IExtension {
	return h.singleCase.ExtensionContainer().GetExtension(pid)
}
