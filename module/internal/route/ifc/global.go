//
//Created by xuzhuoxi
//on 2019-03-23.
//@author xuzhuoxi
//
package ifc

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/encodingx/gobx"
	"github.com/xuzhuoxi/infra-go/lang"
	"time"
)

var (
	DataBlockHandler = bytex.NewDefaultDataBlockHandler() //数据封包处理
	PoolBuffDecoder  = lang.NewObjectPoolSync()
)

const (
	//更新超时时间
	SockTimeout = int64(time.Second * 90)
)

func init() {
	PoolBuffDecoder.Register(func() interface{} {
		return gobx.NewGobBuffDecoder(DataBlockHandler)
	}, func(instance interface{}) bool {
		if nil == instance {
			return false
		}
		if _, ok := instance.(encodingx.IBuffDecoder); ok {
			return true
		}
		return false
	})
}

func HandleBuffDecode(handler func(encodingx.IBuffDecoder)) {
	decoder := getBuffDecoder()
	defer PoolBuffDecoder.Recycle(decoder)
	handler(decoder)
}

//---------------------------

func getBuffDecoder() encodingx.IBuffDecoder {
	rs := PoolBuffDecoder.GetInstance().(encodingx.IBuffDecoder)
	rs.Reset()
	return rs
}
