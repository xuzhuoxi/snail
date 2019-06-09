package config

import (
	"github.com/xuzhuoxi/snail/engine"
	"log"
	"testing"
)

func TestParseConfig(t *testing.T) {
	conf := ParseModuleConfig(engine.DefaultFlagSet)
	log.Println(conf)
}
