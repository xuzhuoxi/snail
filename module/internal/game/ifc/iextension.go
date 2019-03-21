//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package ifc

import "github.com/xuzhuoxi/infra-go/extendx/protox"

type IGameExtension interface {
	protox.IProtocolExtension
	protox.IRequestExtension
}
