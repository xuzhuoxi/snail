//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package demo

import (
	"github.com/xuzhuoxi/snail/module/internal/game/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

type testA struct {
	A string
	B int
	C bool
}

func NewDemoExtension(Name string, SingleCase ifc.IGameSingleCase) *DemoExtension {
	return &DemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(Name, SingleCase)}
}

type DemoExtension struct {
	extension.GameExtensionSupport
}

func (e *DemoExtension) InitProtocolId() {
	e.ProtoIdToValue["D_P2"] = testA{}
}

func (e *DemoExtension) GetRequestData(ProtoId string) (DataCopy interface{}) {
	DataCopy = e.ProtoIdToValue[ProtoId]
	return
}

func (e *DemoExtension) BeforeRequest(ProtoId string) {
	e.Logger().Debugln("DemoExtension.BeforeRequest!", ProtoId)
}

func (e *DemoExtension) AfterRequest(ProtoId string) {
	e.Logger().Debugln("DemoExtension.AfterRequest!", ProtoId)
}

func (e *DemoExtension) OnRequest(ProtoId string, Uid string, Data interface{}, Data2 ...interface{}) {
	e.Logger().Debugln("DemoExtension.OnRequest", ProtoId, Uid, Data, Data2)
}

func (e *DemoExtension) InitExtension() error {
	e.Logger().Debugln("DemoExtension.InitExtension", e.Name)
	return nil
}

func (e *DemoExtension) SaveExtension() error {
	e.Logger().Debugln("DemoExtension.SaveExtension", e.Name)
	return nil
}

func (e *DemoExtension) DestroyExtension() error {
	e.Logger().Debugln("DemoExtension.DestroyExtension", e.Name)
	return nil
}
