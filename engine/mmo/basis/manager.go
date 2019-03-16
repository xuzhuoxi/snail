//
//Created by xuzhuoxi
//on 2019-03-16.
//@author xuzhuoxi
//
package basis

import "github.com/xuzhuoxi/infra-go/logx"

type IManagerBase interface {
	logx.ILoggerSetter
	InitManager()
	DisposeManager()
}
