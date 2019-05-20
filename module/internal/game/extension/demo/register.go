//
//Created by xuzhuoxi
//on 2019-05-12.
//@author xuzhuoxi
//
package demo

import (
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func init() {
	ifc.RegisterExtension(func() ifc.IGameExtension {
		return NewDemoExtension("Demo")
	})
}
