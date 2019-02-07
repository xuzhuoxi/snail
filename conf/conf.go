package conf

import (
	"encoding/json"
	"flag"
	"github.com/xuzhuoxi/snail/snail"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ServiceConf struct {
	Name    string `json:"name,omitempty"`
	Network string `json:"network"`
	Addr    string `json:"addr"`
}

type ObjectConf struct {
	Name    string      `json:"name"`
	Service ServiceConf `json:"service"`
	Module  string      `json:"module"`
	Log     string      `json:"logx,omitempty"`
	RpcName string      `json:"rpc,omitempty"`
	Remotes []string    `json:"remotes,omitempty"`
}

type Conf struct {
	RpcList []ServiceConf `json:"rpcs"`
	Routes  []ObjectConf  `json:"routes,omitempty"`
	Admins  []ObjectConf  `json:"admins,omitempty"`
	Games   []ObjectConf  `json:"games,omitempty"`
	OnList  []string      `json:"onList"`
	mapRPC  map[string]ServiceConf
}

func (c *Conf) handleData() {
	c.mapRPC = make(map[string]ServiceConf)
	for _, val := range c.RpcList {
		c.mapRPC[val.Name] = val
	}
}

func (c ObjectConf) GetRpcInfo() *ServiceConf {
	if c.RpcName == "" {
		return nil
	}
	rs := Config.mapRPC[c.RpcName]
	return &rs
}

func (c *Conf) GetRpcInfo(name string) (*ServiceConf, bool) {
	rs, ok := c.mapRPC[name]
	return &rs, ok
}

var Config *Conf

func ParseConfig(configName string) *Conf {
	//读取运行参数配置文件
	var c = flag.String("c", configName, "GetConfig osxu for running")
	flag.Parse()
	//取当前运行路径
	basePath := filepath.Dir(os.Args[0])
	//读取配置文件
	cfgBody, err := ioutil.ReadFile(basePath + "/conf/" + *c)
	if nil != err {
		log.Fatal(err)
		return nil
	}
	cfg := &Conf{}
	json.Unmarshal(cfgBody, cfg)
	cfg.handleData()
	return cfg
}

func getConfByName(name string) (ObjectConf, error) {
	arr := append(append(Config.Routes, Config.Admins...), Config.Games...)
	for _, val := range arr {
		if val.Name == name {
			return val, nil
		}
	}
	return ObjectConf{}, &snail.Error{"Error"}
}
