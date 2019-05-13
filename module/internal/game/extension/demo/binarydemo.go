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

func NewBinaryDemoExtension(Name string) *BinaryDemoExtension {
	return &BinaryDemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name)}
}

//Extension至少实现两个接口
//IProtocolExtension(必须)
//IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension(选一)
//IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension(可选)
type BinaryDemoExtension struct {
	extension.GameExtensionSupport
}

func (e *BinaryDemoExtension) InitProtocolId() {
	e.ProtoIdToValue["B_0"] = struct{}{}
}

func (e *BinaryDemoExtension) BeforeRequest(protoId string) {
	e.GetLogger().Debugln("BinaryDemoExtension.BeforeRequest!", protoId)
}

func (e *BinaryDemoExtension) OnRequest(resp extendx.IExtensionBinaryResponse, protoId string, uid string, data []byte, data2 ...[]byte) {
	e.GetLogger().Debugln("BinaryDemoExtension.BeforeRequest!", protoId, uid, data, data2)
}

func (e *BinaryDemoExtension) AfterRequest(protoId string) {
	e.GetLogger().Debugln("BinaryDemoExtension.AfterRequest!", protoId)
}

func (e *BinaryDemoExtension) InitExtension() error {
	e.GetLogger().Debugln("BinaryDemoExtension.InitExtension", e.Name)
	return nil
}

func (e *BinaryDemoExtension) SaveExtension() error {
	e.GetLogger().Debugln("BinaryDemoExtension.SaveExtension", e.Name)
	return nil
}

func (e *BinaryDemoExtension) DestroyExtension() error {
	e.GetLogger().Debugln("BinaryDemoExtension.DestroyExtension", e.Name)
	return nil
}
