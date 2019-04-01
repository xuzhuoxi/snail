//
//Created by xuzhuoxi
//on 2019-04-02.
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
)

var (
	// block(带长度的[]byte) -> []byte
	PoolBuffToData = bytex.NewPoolBuffToData()
	// []byte -> block(带长度的[]byte)
	PoolBuffToBlock = bytex.NewPoolBuffToBlock()
	// Gob序列化与反序列化
	// 暂时用于通知Route的对象序列化与反序列
	PoolBuffEncoder = lang.NewObjectPoolSync()
	PoolBuffDecoder = lang.NewObjectPoolSync()
	// Json序列化与反序列化
	// 暂时用于Extension的对象序列化与反序列化
	PoolJsonCodingHandler = lang.NewObjectPoolSync()
	// 用于记录日志：Extension响应时间
	LoggerExtension = logx.NewLogger()
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

func HandleBuffEncodeFromPool(handler func(encodingx.IBuffEncoder)) {
	encoder := getBuffEncoderFromPool()
	defer PoolBuffEncoder.Recycle(encoder)
	handler(encoder)
}

func HandleBuffDecodeFromPool(handler func(encodingx.IBuffDecoder)) {
	decoder := getBuffDecoderFromPool()
	defer PoolBuffDecoder.Recycle(decoder)
	handler(decoder)
}

func HandleBuffToDataFromPool(handler func(bytex.IBuffToData)) {
	buffToData := PoolBuffToData.GetInstance()
	defer PoolBuffToData.Recycle(buffToData)
	handler(buffToData)
}

func HandleBuffToBlockFromPool(handler func(bytex.IBuffToBlock)) {
	buffToBlock := PoolBuffToBlock.GetInstance()
	defer PoolBuffToBlock.Recycle(buffToBlock)
	handler(buffToBlock)
}

func HandleJsonCodingFromPool(handler func(codingHandler encodingx.ICodingHandler)) {
	codingHandler := PoolJsonCodingHandler.GetInstance().(encodingx.ICodingHandler)
	defer PoolJsonCodingHandler.Recycle(codingHandler)
	handler(codingHandler)
}

//---------------------------

func getBuffEncoderFromPool() encodingx.IBuffEncoder {
	rs := PoolBuffEncoder.GetInstance().(encodingx.IBuffEncoder)
	rs.Reset()
	return rs
}

func getBuffDecoderFromPool() encodingx.IBuffDecoder {
	rs := PoolBuffDecoder.GetInstance().(encodingx.IBuffDecoder)
	rs.Reset()
	return rs
}
