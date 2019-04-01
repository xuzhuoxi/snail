//
//Created by xuzhuoxi
//on 2019-03-23.
//@author xuzhuoxi
//
package ifc

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/netx"
	"time"
)

var (
	AddressProxy     = netx.NewIAddressProxy()            //uid与address的交叉映射,整个game模块共享
	DataBlockHandler = bytex.NewDefaultDataBlockHandler() //数据封包处理
)

const (
	//通知Route间隔
	GameNotifyRouteInterval = time.Second * 30
	//统计时间区间
	DefaultStatsInterval = int64(5 * time.Minute)
)
