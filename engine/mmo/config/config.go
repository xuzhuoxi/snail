//
//Created by xuzhuoxi
//on 2019-06-09.
//@author xuzhuoxi
//
package config

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/cmdx"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/osxu"
	"io/ioutil"
)

type Entity struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	MaxUser int    `json:"max"`
}

type Entities struct {
	Worlds []Entity `json:"worlds"`
	Zones  []Entity `json:"zones"`
	Rooms  []Entity `json:"rooms"`
}

type Zone struct {
	ZoneId string   `json:"zone"`
	Rooms  []string `json:"rooms"`
}

type MMO struct {
	World string `json:"world"`
	Zones []Zone `json:"zones"`

	WorldEntity  *Entity
	ZoneEntities []*Entity
	RoomEntities []*Entity
	entityMap    map[string]*Entity
}

func (m *MMO) GetEntity(entityId string) (entity *Entity, ok bool) {
	entity, ok = m.entityMap[entityId]
	return
}

type MMOConfig struct {
	Entities Entities `json:"entities"`
	MMO      MMO      `json:"mmo"`
}

func (c *MMOConfig) HandleData() {
	eMap, err := c.cacheEntityMap()
	if nil != err {
		panic(err.Error())
	}
	c.MMO.entityMap = eMap
	world, zones, rooms, isErr, errorId := c.makeGroup()
	if isErr {
		panic("Entity Undefined: " + errorId)
	}
	c.MMO.WorldEntity = world
	c.MMO.ZoneEntities = zones
	c.MMO.RoomEntities = rooms
}

func (c *MMOConfig) cacheEntityMap() (eMap map[string]*Entity, err error) {
	eMap = make(map[string]*Entity)
	cache2map := func(eMap map[string]*Entity, entity *Entity) error {
		if _, ok := eMap[entity.Id]; ok {
			return errors.New("Entity duplicate definition at id:" + entity.Id)
		}
		eMap[entity.Id] = entity
		return nil
	}
	cacheList := func(eMap map[string]*Entity, list []Entity) error {
		for index, _ := range list {
			err := cache2map(eMap, &list[index])
			if nil != err {
				return err
			}
		}
		return nil
	}
	list := append(append(c.Entities.Worlds, c.Entities.Zones...), c.Entities.Rooms...)
	err = cacheList(eMap, list)
	return
}

func (c *MMOConfig) makeGroup() (world *Entity, zones []*Entity, rooms []*Entity, err bool, errorId string) {
	var ok bool
	eMap := c.MMO.entityMap
	if world, ok = eMap[c.MMO.World]; !ok {
		err, errorId = true, c.MMO.World
		return
	}
	zs := c.MMO.Zones
	for _, z := range zs {
		if _, ok = eMap[z.ZoneId]; !ok {
			err, errorId = true, z.ZoneId
			return
		}
		zones = append(zones, eMap[z.ZoneId])
		for _, rId := range z.Rooms {
			if _, ok = eMap[rId]; !ok {
				err, errorId = true, rId
				return
			}
			rooms = append(rooms, eMap[rId])
		}
	}
	err = false
	return
}

//------------------------------------------

var DefaultMMOConfig *MMOConfig

func ParseMMOConfigByFlag(flagSet *cmdx.FlagSetExtend) *MMOConfig {
	if !flagSet.CheckKey("mmo") {
		return nil
	}
	cfgName, _ := flagSet.GetString("mmo")
	path := filex.Combine(osxu.GetRunningDir(), "conf", cfgName)
	return ParseMMOConfigByPath(path)
}

func ParseMMOConfigByPath(path string) *MMOConfig {
	//读取配置文件
	fmt.Println("Path:", path)
	cfgBody, err := ioutil.ReadFile(path)
	if nil != err {
		panic("mmo does not exist! ")
	}
	return ParseMMOConfigByContent(cfgBody)
}

func ParseMMOConfigByContent(content []byte) *MMOConfig {
	mmoCfg := &MMOConfig{}
	jsoniter.Unmarshal(content, mmoCfg)
	mmoCfg.HandleData()
	return mmoCfg
}
