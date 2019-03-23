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
	"github.com/xuzhuoxi/infra-go/netx"
)

var (
	AddressProxy     = netx.NewIAddressProxy()            //uid与address的交叉映射,整个game模块共享
	DataBlockHandler = bytex.NewDefaultDataBlockHandler() //数据封包处理

	//以下为对象池，全game共享
	PoolBuffToData        = bytex.NewPoolBuffToData()
	PoolBuffToBlock       = bytex.NewPoolBuffToBlock()
	PoolEncoder           = lang.NewObjectPoolSync()
	PoolDecoder           = lang.NewObjectPoolSync()
	PoolJsonCodingHandler = lang.NewObjectPoolSync()
)

func init() {
	PoolBuffToData.Register(func() bytex.IBuffToData {
		return bytex.NewBuffToData(DataBlockHandler)
	})
	PoolBuffToBlock.Register(func() bytex.IBuffToBlock {
		return bytex.NewBuffToBlock(DataBlockHandler)
	})
	PoolEncoder.Register(func() interface{} {
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
	PoolDecoder.Register(func() interface{} {
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

func HandleEncode(handler func(encodingx.IBuffEncoder)) {
	encoder := getBuffEncoder()
	defer PoolEncoder.Recycle(encoder)
	handler(encoder)
}

func HandleDecode(handler func(encodingx.IBuffDecoder)) {
	decoder := getBuffDecoder()
	defer PoolDecoder.Recycle(decoder)
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
	rs := PoolEncoder.GetInstance().(encodingx.IBuffEncoder)
	rs.Reset()
	return rs
}

func getBuffDecoder() encodingx.IBuffDecoder {
	rs := PoolDecoder.GetInstance().(encodingx.IBuffDecoder)
	rs.Reset()
	return rs
}
