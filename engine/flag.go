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
		flagSet = ParseFlag("config.json", "config_mmo.json")
	}
	return flagSet
}

func ParseFlag(config string, mmo string) *cmdx.FlagSetExtend {
	flagSet := cmdx.NewDefaultFlagSetExtend()
	flagSet.String("c", config, "Base Config! ")
	flagSet.String("mmo", mmo, "MMO config! ")
	flagSet.Parse(os.Args[1:])
	return flagSet
}
