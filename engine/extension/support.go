//
//Created by xuzhuoxi
//on 2019-02-17.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/protocolx"
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

type ISnailExtension interface {
	protocolx.IProtocolExtension
}

//----------------------------------------------------

type SnailExtensionSupport struct {
}
