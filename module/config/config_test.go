package config

import (
	"github.com/xuzhuoxi/infra-go/osxu"
	"log"
	"testing"
)

var path = osxu.RunningBaseDir() + "conf/config_module.json"

func TestParseModuleConfig(t *testing.T) {
	conf := ParseModuleConfigByPath(path)
	log.Println(conf)
}
