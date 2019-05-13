//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package demo

import (
	"github.com/xuzhuoxi/infra-go/extendx"
	"github.com/xuzhuoxi/snail/module/internal/game/extension"
)

func NewNoneDemoExtension(Name string) *NoneDemoExtension {
	return &NoneDemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name)}
}

//Extension规范：
//IGameExtension(必须)
//IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension(选一)
//IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension(可选)
type NoneDemoExtension struct {
	extension.GameExtensionSupport
}

func (e *NoneDemoExtension) InitProtocolId() {
	e.ProtoIdToValue["N_0"] = struct{}{}
}

func (e *NoneDemoExtension) BeforeRequest(protoId string) {
	e.GetLogger().Debugln("NoneDemoExtension.BeforeRequest!", protoId)
}

func (e *NoneDemoExtension) OnRequest(resp extendx.IExtensionResponse, protoId string, uid string) {
	e.GetLogger().Debugln("NoneDemoExtension.OnRequest", protoId, uid)
}

func (e *NoneDemoExtension) AfterRequest(protoId string) {
	e.GetLogger().Debugln("NoneDemoExtension.AfterRequest!", protoId)
}

func (e *NoneDemoExtension) InitExtension() error {
	e.GetLogger().Debugln("NoneDemoExtension.InitExtension", e.Name)
	return nil
}

func (e *NoneDemoExtension) SaveExtension() error {
	e.GetLogger().Debugln("NoneDemoExtension.SaveExtension", e.Name)
	return nil
}

func (e *NoneDemoExtension) DestroyExtension() error {
	e.GetLogger().Debugln("NoneDemoExtension.DestroyExtension", e.Name)
	return nil
}
