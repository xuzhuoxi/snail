//
//Created by xuzhuoxi
//on 2019-06-10.
//@author xuzhuoxi
//
package manager

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/osxu"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"github.com/xuzhuoxi/snail/engine/mmo/config"
	"testing"
)

var path = filex.Combine(osxu.GetRunningDir(), "conf/config_mmo.json")

func TestEntityManager_ConstructWorld(t *testing.T) {
	cfg := config.ParseMMOConfigByPath(path)
	fmt.Println(cfg)
	eMgr := NewIEntityManager()
	eMgr.ConstructWorld(cfg)
	eMgr.World().ForEachChild(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
		logx.Traceln(child.UID(), child.(basis.IEntityChild).GetParent())
		return
	})
}
