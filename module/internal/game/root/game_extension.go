//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package root

import (
	"github.com/xuzhuoxi/infra-go/logx"
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
	container := singleCase.ExtensionContainer()
	container.AppendExtension(demo.NewDemoExtension("Demo1", singleCase))
	container.AppendExtension(demo.NewDemoExtension("Demo2", singleCase))
}

func (c *ExtensionConfig) InitExtensions() {
	c.singleCase.ExtensionContainer().InitExtensions()
}

func (c *ExtensionConfig) Logger() logx.ILogger {
	return c.singleCase.Logger()
}
