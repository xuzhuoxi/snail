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

type testA struct {
	A string
	B int
	C bool
}

func NewObjDemoExtension(Name string) *ObjDemoExtension {
	return &ObjDemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name)}
}

//Extension至少实现两个接口
//IProtocolExtension(必须)
//IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension(选一)
//IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension(可选)
type ObjDemoExtension struct {
	extension.GameExtensionSupport
}

func (e *ObjDemoExtension) Batch() bool {
	return true
}

func (e *ObjDemoExtension) InitProtocolId() {
	e.ProtoIdToValue["Obj_0"] = testA{}
}

func (e *ObjDemoExtension) BeforeRequest(protoId string) {
	e.GetLogger().Debugln("ObjDemoExtension.BeforeRequest!", protoId)
}

func (e *ObjDemoExtension) AfterRequest(protoId string) {
	e.GetLogger().Debugln("ObjDemoExtension.AfterRequest!", protoId)
}

func (e *ObjDemoExtension) GetRequestData(ProtoId string) (dataCopy interface{}) {
	dataCopy = e.ProtoIdToValue[ProtoId]
	return
}

func (e *ObjDemoExtension) OnRequest(resp extendx.IExtensionObjectResponse, protoId string, uid string, data interface{}, data2 ...interface{}) {
	e.GetLogger().Debugln("ObjDemoExtension.OnRequest", protoId, uid, data, data2)
}

func (e *ObjDemoExtension) InitExtension() error {
	e.GetLogger().Debugln("ObjDemoExtension.InitExtension", e.Name)
	return nil
}

func (e *ObjDemoExtension) SaveExtension() error {
	e.GetLogger().Debugln("ObjDemoExtension.SaveExtension", e.Name)
	return nil
}

func (e *ObjDemoExtension) DestroyExtension() error {
	e.GetLogger().Debugln("ObjDemoExtension.DestroyExtension", e.Name)
	return nil
}
