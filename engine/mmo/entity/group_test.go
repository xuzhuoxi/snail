//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package entity

import (
	"fmt"
	"testing"
)

func TestMapLen(t *testing.T) {
	m := make(map[string]*struct{})
	fmt.Println(len(m))
	m["aaa"] = nil
	fmt.Println(len(m))
	m["bbb"] = nil
	fmt.Println(len(m))
}
