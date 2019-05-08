package conf

import (
	"flag"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/osxu"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func BaseLogPath() string {
	return osxu.RunningBaseDir() + "log/"
}

type SockConf struct {
	Name    string `json:"name"`
	Network string `json:"network"`
	Addr    string `json:"addr"`
}

type ObjectConf struct {
	Id         string   `json:"id"`
	ModuleName string   `json:"module"`
	RpcList    []string `json:"rpc,omitempty"`
	SockList   []string `json:"socks,omitempty"`
	Remotes    []string `json:"remotes,omitempty"`
	Log        string   `json:"log,omitempty"`

	conf *Conf
}

func (oc ObjectConf) LogFileInfo() (fileDir string, fileBaseName string, fileExtName string) {
	fullPath := BaseLogPath() + oc.Log
	var fileName string
	fileDir, fileName = osxu.SplitFilePath(fullPath)
	fileBaseName, fileExtName = osxu.SplitFileName(fileName)
	return
}

func (c ObjectConf) GetServiceConf(name string) (SockConf, bool) {
	return c.conf.GetServiceConf(name)
}

type Conf struct {
	Socks  []SockConf   `json:"socks,omitempty"`
	Routes []ObjectConf `json:"routes,omitempty"`
	Admins []ObjectConf `json:"admins,omitempty"`
	Games  []ObjectConf `json:"games,omitempty"`
	OnList []string     `json:"onList"`

	mapService map[string]SockConf
	mapObject  map[string]ObjectConf
}

func (c *Conf) GetServiceConf(name string) (SockConf, bool) {
	rs, has := c.mapService[name]
	if has {
		return rs, true
	}
	return SockConf{}, false
}

func (c *Conf) GetObjectById(id string) (ObjectConf, bool) {
	val, has := c.mapObject[id]
	if has {
		return val, true
	}
	return ObjectConf{}, false
}

func (c *Conf) handleData() {
	c.mapService = make(map[string]SockConf)
	for _, val := range c.Socks {
		c.mapService[val.Name] = val
	}
	objectToMap := func(m map[string]ObjectConf, objects []ObjectConf) {
		if len(objects) == 0 {
			return
		}
		for _, val := range objects {
			_, has := m[val.Id]
			if has {
				panic("Id Repeat!")
			}
			val.conf = c
			m[val.Id] = val
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
	jsoniter.Unmarshal(cfgBody, cfg)
	cfg.handleData()
	return cfg
}

func GetServiceConf(name string) (SockConf, bool) {
	return DefaultConfig.GetServiceConf(name)
}

func GetObjectById(name string) (ObjectConf, bool) {
	return DefaultConfig.GetObjectById(name)
}
