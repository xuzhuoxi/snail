//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package ifc

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

type IGameSingleCase interface {
	logx.ILoggerGetter
	logx.ILoggerSetter

	Init()
}

type IGameSingleCaseSetter interface {
	SetSingleCase(singleCase IGameSingleCase)
}
