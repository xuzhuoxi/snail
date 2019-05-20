//
//Created by xuzhuoxi
//on 2019-02-17.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
)

type ISnailExtensionContainer interface {
	protox.IProtocolContainer
}

func NewISnailExtensionContainer() ISnailExtensionContainer {
	return NewSnailExtensionContainer()
}

func NewSnailExtensionContainer() *SnailExtensionContainer {
	return &SnailExtensionContainer{ProtocolContainer: protox.NewProtocolExtensionContainer()}
}

//----------------------------------------------------

type SnailExtensionContainer struct {
	protox.ProtocolContainer
}
