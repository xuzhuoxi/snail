//
//Created by xuzhuoxi
//on 2019-02-17.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
)

type ISnailExtensionContainer interface {
	protox.IProtocolContainer
	InitExtensions()
	SaveExtensions()
	DestroyExtensions()
}

func NewISnailExtensionContainer() ISnailExtensionContainer {
	return &SnailExtensionContainer{ProtocolContainer: protox.NewProtocolExtensionContainer()}
}

func NewSnailExtensionContainer() SnailExtensionContainer {
	return SnailExtensionContainer{ProtocolContainer: protox.NewProtocolExtensionContainer()}
}

//----------------------------------------------------

type SnailExtensionContainer struct {
	protox.ProtocolContainer
}

func (c *SnailExtensionContainer) InitExtensions() {
	ln := c.Len()
	if ln == 0 {
		return
	}
	c.Range(func(_ int, extension extendx.IExtension) {
		if e, ok := extension.(ISnailInitExtension); ok {
			e.InitExtension()
		}
	})
}

func (c *SnailExtensionContainer) SaveExtensions() {
	ln := c.Len()
	if ln == 0 {
		return
	}
	c.Range(func(_ int, extension extendx.IExtension) {
		if e, ok := extension.(ISnailSaveExtension); ok {
			e.SaveExtension()
		}
	})
}

func (c *SnailExtensionContainer) DestroyExtensions() {
	ln := c.Len()
	if ln == 0 {
		return
	}
	c.RangeReverse(func(_ int, extension extendx.IExtension) {
		if e, ok := extension.(ISnailDestroyExtension); ok {
			e.DestroyExtension()
		}
	})
}
