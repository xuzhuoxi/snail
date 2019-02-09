package conf

import (
	"encoding/json"
	"flag"
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

	conf *Conf
}

func (oc ObjectConf) LogDir() string {
	basePath := filepath.Dir(os.Args[0])
	return basePath + "/log/"
}

func (c ObjectConf) GetServiceConf(name string) (*ServiceConf, bool) {
	return c.conf.GetServiceConf(name)
}

type Conf struct {
	Services []ServiceConf `json:"services,omitempty`
	Routes   []ObjectConf  `json:"routes,omitempty"`
	Admins   []ObjectConf  `json:"admins,omitempty"`
	Games    []ObjectConf  `json:"games,omitempty"`
	OnList   []string      `json:"onList"`

	mapService map[string]ServiceConf
	mapObject  map[string]ObjectConf
}

func (c *Conf) GetServiceConf(name string) (*ServiceConf, bool) {
	rs, has := c.mapService[name]
	if has {
		return &rs, true
	}
	return nil, false
}

func (c *Conf) GetConfByName(name string) (*ObjectConf, bool) {
	val, has := c.mapObject[name]
	if has {
		return &val, true
	}
	return nil, false
}

func (c *Conf) handleData() {
	c.mapService = make(map[string]ServiceConf)
	for _, val := range c.Services {
		c.mapService[val.Name] = val
	}
	objectToMap := func(m map[string]ObjectConf, objects []ObjectConf) {
		if len(objects) == 0 {
			return
		}
		for _, val := range objects {
			_, has := m[val.Name]
			if has {
				panic("ObjectName Repeat!")
			}
			val.conf = c
			m[val.Name] = val
		}
	}
	c.mapObject = make(map[string]ObjectConf)
	objectToMap(c.mapObject, c.Routes)
	objectToMap(c.mapObject, c.Games)
	objectToMap(c.mapObject, c.Admins)
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
	cfg.handleData()
	return cfg
}

func GetServiceConf(name string) (*ServiceConf, bool) {
	return DefaultConfig.GetServiceConf(name)
}

func GetConfByName(name string) (*ObjectConf, bool) {
	return DefaultConfig.GetConfByName(name)
}
