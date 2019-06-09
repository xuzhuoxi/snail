package config

import (
	"log"
	"testing"
)

func TestParseMMOConfig(t *testing.T) {
	mmoCfg := ParseMMOConfig()
	mmoCfg.HandleData()
	log.Println(mmoCfg)
}
