//
//Created by xuzhuoxi
//on 2019-06-08.
//@author xuzhuoxi
//
package engine

import (
	"github.com/xuzhuoxi/infra-go/cmdx"
	"os"
	"sync"
)

var (
	flagSet   *cmdx.FlagSetExtend
	flagSetMu sync.Mutex
)

func GetDefaultFlagSet() *cmdx.FlagSetExtend {
	flagSetMu.Lock()
	defer flagSetMu.Unlock()
	if nil == flagSet {
		flagSet = ParseFlag()
	}
	return flagSet
}

func ParseFlag() *cmdx.FlagSetExtend {
	flagSet := cmdx.NewDefaultFlagSetExtend()
	flagSet.String("c", "config_module.json", "Base Config! ")
	flagSet.String("mmo", "config_mmo.json", "MMO config! ")
	flagSet.Parse(os.Args[1:])
	return flagSet
}
