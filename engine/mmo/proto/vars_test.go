//
//Created by xuzhuoxi
//on 2019-03-18.
//@author xuzhuoxi
//
package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"github.com/xuzhuoxi/snail/engine/mmo/entity"
	"testing"
	"unsafe"
)

func TestPtr(t *testing.T) {
	var data interface{} = false
	var data2 *interface{} = &data
	fmt.Println("data:", data, *data2)
}

func TestBuff(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	var data interface{} = float32(23)
	binary.Write(buff, binary.BigEndian, &data)
	fmt.Println(buff.Bytes())
}

func TestSize(t *testing.T) {
	fmt.Println("int:", binary.Size(1111))
	fmt.Println("int16:", binary.Size(int16(1111)))
	fmt.Println("uint:", binary.Size(uint(1111)))
	fmt.Println("uint16:", binary.Size(uint16(1111)))
}

func TestType(t *testing.T) {
	catchType := func(e interface{}) {
		switch e := e.(type) {
		case bool:
			fmt.Println("bool")
		case *bool:
			fmt.Println("*bool")
		default:
			fmt.Println("default", e)
		}
	}
	var data interface{} = true
	var pdata *interface{} = &data
	catchType(data)
	catchType(pdata)
	fmt.Println("无敌分界线——————————————")
	var data2 = true
	catchType(data2)
	catchType(&data2)
	fmt.Println("无敌分界线——————————————")

	fmt.Println(unsafe.Sizeof(data), unsafe.Sizeof(&data), unsafe.Sizeof(data2), unsafe.Sizeof(&data2))
	fmt.Println(unsafe.Sizeof(struct{}{}), unsafe.Sizeof(make(map[string]struct{})), unsafe.Sizeof(make(map[string]string)))

	//var data3 = &data
	//var data4 = &data2
	//
	//fmt.Println(reflect.TypeOf(data3).Kind() == reflect.Ptr)
	//fmt.Println(reflect.TypeOf(data4).Kind() == reflect.Ptr)

	//这个interface{}真麻烦，具体类型转为interface{}时好像被内嵌了
}

func TestBytes(t *testing.T) {
	room := entity.NewRoomEntity("111", "顶你个肺")
	varSet := basis.NewVarSet()
	//varSet.Set("key0", false)
	//varSet.Set("key1", true)
	//varSet.Set("key2", 222)
	//varSet.Set("key3", 222.5)
	//varSet.Set("key4", "aaa，哈哈")
	//varSet.Set("key5", []bool{false, true})
	//varSet.Set("key6", []uint16{222, 333})
	varSet.Set("key7", []int{222, 333})
	//varSet.Set("key8", []string{"", "o只"})

	bs := VarSetToBytes(room, varSet)
	fmt.Println("结果：", bs)

	et, eid, vs := BytesToVarSet(bs)
	fmt.Println(et, eid, vs)
}
