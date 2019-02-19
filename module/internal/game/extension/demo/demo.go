//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package demo

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/module/internal/game/extension"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

func NewDemoExtension(ProtoId string, SingleCase intfc.IGameSingleCase) *DemoExtension {
	return &DemoExtension{GameExtensionSupport: extension.NewGameExtensionSupport(ProtoId, SingleCase)}
}

type DemoExtension struct {
	extension.GameExtensionSupport
}

func (e *DemoExtension) Batch() bool {
	return false
}

func (e *DemoExtension) HandleRequest(pId string, data interface{}, data2 ...interface{}) {
	e.Logger().Infoln("DemoExtension.HandleRequest", pId, data, data2)
}

func (e *DemoExtension) InitExtension() error {
	e.Logger().Infoln("DemoExtension.InitExtension", e.ProtoId)
	return nil
}

func (e *DemoExtension) SaveExtension() error {
	e.Logger().Infoln("DemoExtension.SaveExtension", e.ProtoId)
	return nil
}

func (e *DemoExtension) DestroyExtension() error {
	e.Logger().Infoln("DemoExtension.DestroyExtension", e.ProtoId)
	return nil
}

func (e *DemoExtension) Logger() logx.ILogger {
	return e.SingleCase.Logger()
}
