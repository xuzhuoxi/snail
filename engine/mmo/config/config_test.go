package config

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"testing"
)

var path = osxu.GetRunningDir() + "/conf/config_mmo.json"

func TestParseMMOConfig(t *testing.T) {
	cfg := ParseMMOConfigByPath(path)
	fmt.Println(*cfg)
}
