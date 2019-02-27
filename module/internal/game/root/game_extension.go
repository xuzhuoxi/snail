//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/snail/module/internal/game/extension/demo"
	"github.com/xuzhuoxi/snail/module/internal/game/intfc"
)

func NewExtensionConfig(singleCase intfc.IGameSingleCase) *ExtensionConfig {
	return &ExtensionConfig{singleCase: singleCase}
}

type ExtensionConfig struct {
	singleCase intfc.IGameSingleCase
}

func (c *ExtensionConfig) ConfigExtensions() {
	singleCase := c.singleCase
	c.appendConfig("Demo1", demo.NewDemoExtension("Demo1", singleCase))
	c.appendConfig("Demo2", demo.NewDemoExtension("Demo2", singleCase))
}

func (c *ExtensionConfig) InitExtensions() {
	c.singleCase.ExtensionContainer().InitExtensions()
}

func (c *ExtensionConfig) appendConfig(pid string, extension intfc.IGameExtension) {
	c.singleCase.ExtensionContainer().AppendExtension(extension)
}
