//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package demo

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/snail/module/internal/game/extension"
)

func NewDemoExtension(Name string) *DemoExtension {
	return &DemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name)}
}

func NewOrigin() interface{} {
	return &originObj{}
}

type originObj struct {
	Data  int
	Data2 string
}

type paramHandler struct {
}

func (h *paramHandler) SetCodingHandler(handler encodingx.ICodingHandler) {
	panic("implement me")
}

func (h *paramHandler) HandleRequestParam(data []byte) interface{} {
	panic("implement me")
}

func (h *paramHandler) HandleRequestParams(data [][]byte) []interface{} {
	panic("implement me")
}

func (h *paramHandler) HandleResponseParam(data interface{}) []byte {
	panic("implement me")
}

func (h *paramHandler) HandleResponseParams(data []interface{}) [][]byte {
	panic("implement me")
}

//Extension至少实现两个接口
//IProtocolExtension(必须)
//IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension(选一)
//IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension(可选)
type DemoExtension struct {
	extension.GameExtensionSupport
}

func (e *DemoExtension) InitExtension() error {
	e.GetLogger().Debugln("DemoExtension.InitExtension", e.Name)
	e.SetRequestHandler("N_0", e.onRequestNoneParam)
	e.SetRequestHandlerBinary("B_0", e.onRequestBinary)
	e.SetRequestHandlerString("J_0", e.onRequestJson)
	e.SetRequestHandlerObject("Obj_0", e.onRequestObj, NewOrigin, nil)
	return nil
}

func (e *DemoExtension) DestroyExtension() error {
	e.ClearRequestHandler("J_0")
	e.ClearRequestHandler("B_0")
	e.GetLogger().Debugln("DemoExtension.DestroyExtension", e.Name)
	return nil
}

func (e *DemoExtension) BeforeRequest(protoId string) {
	e.GetLogger().Debugln("DemoExtension.BeforeRequest!", protoId)
}

func (e *DemoExtension) onRequestNoneParam(resp protox.IExtensionResponse, req protox.IExtensionRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestNoneParam!", req, resp)
}

func (e *DemoExtension) onRequestBinary(resp protox.IExtensionBinaryResponse, req protox.IExtensionBinaryRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestBinary!", req, resp)
}

func (e *DemoExtension) onRequestJson(resp protox.IExtensionStringResponse, req protox.IExtensionStringRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestJson!", req, resp)
}

func (e *DemoExtension) onRequestObj(resp protox.IExtensionObjectResponse, req protox.IExtensionObjectRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestObj!", req, resp)
}

func (e *DemoExtension) AfterRequest(protoId string) {
	e.GetLogger().Debugln("DemoExtension.AfterRequest!", protoId)
}

func (e *DemoExtension) SaveExtension() error {
	e.GetLogger().Debugln("DemoExtension.SaveExtension", e.Name)
	return nil
}
