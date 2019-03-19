//
//Created by xuzhuoxi
//on 2019-03-18.
//@author xuzhuoxi
//
package proto

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/binaryx"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

//序列化
//格式:
// 	EntityType:uint16
//	string: 长度+[]byte
//	Kind: 	uint8
//	key:	长度+[]byte
//	value:	Kind + Value 或 Kind + Len(uint8) + Value
func VarSetToBytes(source basis.IEntity, varSet basis.VarSet) []byte {
	buff := bytex.NewBuffDataBlock(DataBlockHandler)
	buff.WriteBinary(source.EntityType()) //EntityType
	buff.WriteData([]byte(source.UID()))  //string
	fmt.Println("VarSetToBytes:", buff.Bytes())
	for key, val := range varSet {
		if !binaryx.CheckValue(val) { //非法值
			continue
		}
		buff.WriteData([]byte(key)) //key
		kind, ln := binaryx.GetValueKind(val)
		buff.WriteBinary(kind)         //Kind
		if binaryx.IsSliceKind(kind) { //Len
			buff.WriteBinary(uint8(ln))
		}
		switch kind {
		case binaryx.KindSliceString: //Value=[]string
			for _, str := range val.([]string) {
				buff.WriteData([]byte(str))
			}
		case binaryx.KindString: //Value=string
			buff.WriteData([]byte(val.(string)))
		default: //Value
			buff.WriteBinary(val)
		}
		fmt.Println("VarSetToBytes:", buff.Bytes())
	}
	return buff.ReadBytes()
}

//反序列化
func BytesToVarSet(bs []byte) (entityType basis.EntityType, entityId string, varSet basis.VarSet) {
	buff := bytex.NewBuffDataBlock(DataBlockHandler)
	buff.Write(bs)
	buff.ReadBinary(&entityType)       //EntityType
	entityId = string(buff.ReadData()) //string
	varSet = basis.NewVarSet()
	for buff.Len() > 0 {
		key := string(buff.ReadData()) //key
		var kind binaryx.ValueKind
		buff.ReadBinary(&kind) //Kind
		var ln uint8
		if binaryx.IsSliceKind(kind) { //Len
			buff.ReadBinary(&ln)
		}
		var val interface{}
		switch kind {
		case binaryx.KindSliceString: //Value=[]string
			var strSlice []string
			for ln < 0 {
				strSlice = append(strSlice, string(buff.ReadData()))
			}
			val = strSlice
		case binaryx.KindString: //Value=string
			val = string(buff.ReadData())
		default: //Value
			rs := binaryx.GetKindValue(kind, int(ln))
			if binaryx.IsSliceKind(kind) {
				buff.ReadBinary(rs)
			} else {
				buff.ReadBinary(&rs)
			}
			val = rs
		}
		varSet.Set(key, val)
		//fmt.Println("BytesToVarSet, buff.Len())
	}
	return
}
