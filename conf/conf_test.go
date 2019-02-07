package conf

import (
	"log"
	"testing"
)

func TestParseConfig(t *testing.T) {
	conf := ParseConfig("config.json")
	log.Println(conf)
}
