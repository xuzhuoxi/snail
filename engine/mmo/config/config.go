//
//Created by xuzhuoxi
//on 2019-06-09.
//@author xuzhuoxi
//
package config

import (
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/cmdx"
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

type MMOConfig struct {
	Entities Entities `json:"entities"`
	MMO      MMO      `json:"mmo"`
}

func (c *MMOConfig) HandleData() {
	eMap := c.cacheEntityMap()
	c.MMO.entityMap = eMap
	world, zones, rooms, err, errorId := c.makeGroup()
	if err {
		panic("Entity Undefined: " + errorId)
	}
	c.MMO.WorldEntity = world
	c.MMO.ZoneEntities = zones
	c.MMO.RoomEntities = rooms
}

func (c *MMOConfig) cacheEntityMap() map[string]*Entity {
	eMap := make(map[string]*Entity)
	for index, _ := range c.Entities.Worlds {
		eMap[c.Entities.Worlds[index].Id] = &c.Entities.Worlds[index]
	}
	for index, _ := range c.Entities.Zones {
		eMap[c.Entities.Zones[index].Id] = &c.Entities.Zones[index]
	}
	for index, _ := range c.Entities.Rooms {
		eMap[c.Entities.Rooms[index].Id] = &c.Entities.Rooms[index]
	}
	return eMap
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

func ParseMMOConfig(flagSet *cmdx.FlagSetExtend) *MMOConfig {
	if !flagSet.CheckKey("mmo") {
		return nil
	}
	cfgName, _ := flagSet.GetString("mmo")
	//读取配置文件
	cfgBody, err := ioutil.ReadFile(osxu.RunningBaseDir() + "/conf/" + cfgName)
	if nil != err {
		panic("mmo does not exist! ")
	}
	mmoCfg := &MMOConfig{}
	jsoniter.Unmarshal(cfgBody, mmoCfg)
	mmoCfg.HandleData()
	return mmoCfg
}
