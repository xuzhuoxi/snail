//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package demo

import (
	"github.com/xuzhuoxi/snail/module/internal/game/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

type testA struct {
	A string
	B int
	C bool
}

func NewDemoExtension(ProtoId string, SingleCase intfc.IGameSingleCase) *DemoExtension {
	return &DemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(ProtoId, SingleCase)}
}

type DemoExtension struct {
	extension.GameExtensionSupport
}

func (e *DemoExtension) RequestData() interface{} {
	return testA{}
}

func (e *DemoExtension) BeforeRequest() {
	e.Logger().Debugln("DemoExtension.BeforeRequest!")
}

func (e *DemoExtension) AfterRequest() {
	e.Logger().Debugln("DemoExtension.AfterRequest!")
}

func (e *DemoExtension) OnRequest(pId string, data interface{}, data2 ...interface{}) {
	e.Logger().Debugln("DemoExtension.OnRequest", pId, data, data2)
}

func (e *DemoExtension) InitExtension() error {
	e.Logger().Debugln("DemoExtension.InitExtension", e.ProtoId)
	return nil
}

func (e *DemoExtension) SaveExtension() error {
	e.Logger().Debugln("DemoExtension.SaveExtension", e.ProtoId)
	return nil
}

func (e *DemoExtension) DestroyExtension() error {
	e.Logger().Debugln("DemoExtension.DestroyExtension", e.ProtoId)
	return nil
}
