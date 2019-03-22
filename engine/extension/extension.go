//
//Created by xuzhuoxi
//on 2019-02-27.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
)

type ISnailInitExtension interface {
	InitExtension() error
}

type ISnailSaveExtension interface {
	SaveExtension() error
}

type ISnailDestroyExtension interface {
	DestroyExtension() error
}

type ISnailExtension = protox.IProtocolExtension
