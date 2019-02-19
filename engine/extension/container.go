//
//Created by xuzhuoxi
//on 2019-02-17.
//@author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx"
	"github.com/xuzhuoxi/infra-go/protocolx"
)

type ISnailExtensionContainer interface {
	protocolx.IProtocolContainer
	InitExtensions()
	SaveExtensions()
	DestroyExtensions()
}

func NewSnailExtensionContainer() ISnailExtensionContainer {
	return &SnailExtensionContainer{IProtocolContainer: protocolx.NewProtocolExtensionContainer()}
}

//----------------------------------------------------

type SnailExtensionContainer struct {
	protocolx.IProtocolContainer
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
