//
//Created by xuzhuoxi
//on 2019-05-12.
//@author xuzhuoxi
//
package user

import (
	"github.com/xuzhuoxi/snail/module/internal/game/ifc"
)

func init() {
	ifc.RegisterExtension(func() ifc.IGameExtension {
		return NewLoginExtension("Login")
	})
}
