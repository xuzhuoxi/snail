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
	"github.com/xuzhuoxi/infra-go/encodingx/jsonx"
	"github.com/xuzhuoxi/infra-go/lang"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/module/imodule"
	"time"
)

var (
	AddressProxy     = netx.NewIAddressProxy()            //uid与address的交叉映射,整个game模块共享
	DataBlockHandler = bytex.NewDefaultDataBlockHandler() //数据封包处理

	//以下为对象池，全game共享
	PoolBuffToData        = bytex.NewPoolBuffToData()
	PoolBuffToBlock       = bytex.NewPoolBuffToBlock()
	PoolBuffEncoder       = lang.NewObjectPoolSync()
	PoolBuffDecoder       = lang.NewObjectPoolSync()
	PoolJsonCodingHandler = lang.NewObjectPoolSync()

	LoggerExtension = logx.NewLogger()
)

const (
	GameNotifyRouteInterval = time.Duration(imodule.DefaultStatsInterval)
)

func init() {
	PoolBuffToData.Register(func() bytex.IBuffToData {
		return bytex.NewBuffToData(DataBlockHandler)
	})
	PoolBuffToBlock.Register(func() bytex.IBuffToBlock {
		return bytex.NewBuffToBlock(DataBlockHandler)
	})
	PoolBuffEncoder.Register(func() interface{} {
		return gobx.NewGobBuffEncoder(DataBlockHandler)
	}, func(instance interface{}) bool {
		if nil == instance {
			return false
		}
		if _, ok := instance.(encodingx.IBuffEncoder); ok {
			return true
		}
		return false
	})
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
	PoolJsonCodingHandler.Register(func() interface{} {
		return jsonx.NewJsonCodingHandlerSync()
	}, func(instance interface{}) bool {
		if nil == instance {
			return false
		}
		if _, ok := instance.(encodingx.ICodingHandler); ok {
			return true
		}
		return false
	})
}

func HandleBuffEncode(handler func(encodingx.IBuffEncoder)) {
	encoder := getBuffEncoder()
	defer PoolBuffEncoder.Recycle(encoder)
	handler(encoder)
}

func HandleBuffDecode(handler func(encodingx.IBuffDecoder)) {
	decoder := getBuffDecoder()
	defer PoolBuffDecoder.Recycle(decoder)
	handler(decoder)
}

func HandleBuffToData(handler func(bytex.IBuffToData)) {
	buffToData := PoolBuffToData.GetInstance()
	defer PoolBuffToData.Recycle(buffToData)
	handler(buffToData)
}

func HandleBuffToBlock(handler func(bytex.IBuffToBlock)) {
	buffToBlock := PoolBuffToBlock.GetInstance()
	defer PoolBuffToBlock.Recycle(buffToBlock)
	handler(buffToBlock)
}

func HandleJsonCoding(handler func(codingHandler encodingx.ICodingHandler)) {
	codingHandler := PoolJsonCodingHandler.GetInstance().(encodingx.ICodingHandler)
	defer PoolJsonCodingHandler.Recycle(codingHandler)
	handler(codingHandler)
}

//---------------------------

func getBuffEncoder() encodingx.IBuffEncoder {
	rs := PoolBuffEncoder.GetInstance().(encodingx.IBuffEncoder)
	rs.Reset()
	return rs
}

func getBuffDecoder() encodingx.IBuffDecoder {
	rs := PoolBuffDecoder.GetInstance().(encodingx.IBuffDecoder)
	rs.Reset()
	return rs
}
