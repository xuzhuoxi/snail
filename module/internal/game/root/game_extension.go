//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/timex"
	"github.com/xuzhuoxi/snail/module/imodule"
	_ "github.com/xuzhuoxi/snail/module/internal/game/extension/demo"
	_ "github.com/xuzhuoxi/snail/module/internal/game/extension/user"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
	"time"
)

func injectExtensions(container ifc.IGameExtensionContainer, single ifc.IGameSingleCase) {
	if nil == container || nil == single {
		return
	}
	ifc.ForeachExtensionConstructor(func(constructor ifc.GameExtensionConstructor) {
		extension := constructor()
		extension.SetSingleCase(single)
		container.AppendExtension(extension)
	})
}

func NewSnailGameExtensionManager(SockStateDetail *imodule.SockStateDetail) protox.IExtensionManager {
	rs := &SnailGameExtensionManager{ExtensionManager: *protox.NewExtensionManager(), SockStateDetail: SockStateDetail}
	return rs
}

type SnailGameExtensionManager struct {
	protox.ExtensionManager

	SockStateDetail *imodule.SockStateDetail
}

func (m *SnailGameExtensionManager) StartManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Container.InitExtensions()
	m.SockServer.GetPackHandlerContainer().AppendPackHandler(m.onSnailGamePack)
}

func (m *SnailGameExtensionManager) StopManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.SockServer.GetPackHandlerContainer().ClearHandler(m.onSnailGamePack)
	m.Container.DestroyExtensions()
}

func (m *SnailGameExtensionManager) onSnailGamePack(msgData []byte, senderAddress string, other interface{}) bool {
	//m.Logger.Infoln("ExtensionManager.onPack", senderAddress, msgData)
	m.SockStateDetail.AddReqCount()
	name, pid, uid, data := m.ParseMessage(msgData)
	extension, ok := m.Verify(name, pid, uid)
	if !ok {
		return false
	}
	//参数处理
	response, request := m.GenParams(extension, senderAddress, name, pid, uid, data)
	defer func() {
		protox.DefaultRequestPool.Recycle(request)
		protox.DefaultResponsePool.Recycle(response)
	}()
	//响应处理
	if be, ok := extension.(protox.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(request)
	}
	if re, ok := extension.(protox.IRequestExtension); ok {
		func() { //记录时间状态
			tn := time.Now().UnixNano()
			defer func() {
				un := time.Now().UnixNano() - tn
				ifc.LoggerExtension.Infoln(name, pid, un, timex.FormatUnixMilli(un/1e6, "5.999999ms")) //记录响应时间
				m.SockStateDetail.AddRespUnixNano(un)
			}()
			re.OnRequest(response, request)
		}()
	}
	if ae, ok := extension.(protox.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(response, request)
	}
	return true
}
