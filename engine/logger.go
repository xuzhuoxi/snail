//
//Created by xuzhuoxi
//on 2019-06-07.
//@author xuzhuoxi
//
package engine

import "github.com/xuzhuoxi/infra-go/logx"

var SnailLogger = logx.NewLogger()

func init() {
	SnailLogger.SetPrefix("[Snail] ")
	SnailLogger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
}
