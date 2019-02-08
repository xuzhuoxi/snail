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
	Name    string `json:"name"`
	Network string `json:"network"`
	Addr    string `json:"addr"`
}

type ObjectConf struct {
	Name        string   `json:"name"`
	Module      string   `json:"module"`
	RpcList     []string `json:"rpc,omitempty"`
	ServiceList []string `json:"service,omitempty"`
	Log         string   `json:"log,omitempty"`
}

func (oc ObjectConf) LogDir() string {
	basePath := filepath.Dir(os.Args[0])
	return basePath + "/log/"
}

type Conf struct {
	Services []ServiceConf `json:"services,omitempty`
	Routes   []ObjectConf  `json:"routes,omitempty"`
	Admins   []ObjectConf  `json:"admins,omitempty"`
	Games    []ObjectConf  `json:"games,omitempty"`
	OnList   []string      `json:"onList"`

	mapService map[string]*ServiceConf
	mapObject  map[string]*ObjectConf
}

var DefaultConfig *Conf

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
	handleData(cfg)
	return cfg
}

func GetServiceConf(name string) (*ServiceConf, bool) {
	rs, ok := DefaultConfig.mapService[name]
	return rs, ok
}

func GetConfByName(name string) (ObjectConf, error) {
	arr := append(append(DefaultConfig.Routes, DefaultConfig.Admins...), DefaultConfig.Games...)
	for _, val := range arr {
		if val.Name == name {
			return val, nil
		}
	}
	return ObjectConf{}, &snail.Error{"Error"}
}

//private-----------------------

func handleData(c *Conf) {
	c.mapService = make(map[string]*ServiceConf)
	for _, val := range c.Services {
		c.mapService[val.Name] = &val
	}
	objectToMap := func(m map[string]*ObjectConf, objects []ObjectConf) {
		if len(objects) == 0 {
			return
		}
		for _, val := range objects {
			_, has := m[val.Name]
			if has {
				panic("ObjectName Repeat!")
			}
			m[val.Name] = &val
		}
	}
	c.mapObject = make(map[string]*ObjectConf)
	objectToMap(c.mapObject, c.Routes)
	objectToMap(c.mapObject, c.Games)
	objectToMap(c.mapObject, c.Admins)
}
