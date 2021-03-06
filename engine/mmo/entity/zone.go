//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package entity

import "github.com/xuzhuoxi/snail/engine/mmo/basis"

type ZoneConfig struct {
}

func NewIZoneEntity(zoneId string, zoneName string) basis.IZoneEntity {
	return &ZoneEntity{ZoneId: zoneId, ZoneName: zoneName}
}

func NewZoneEntity(zoneId string, zoneName string) *ZoneEntity {
	return &ZoneEntity{ZoneId: zoneId, ZoneName: zoneName}
}

type ZoneEntity struct {
	ZoneId   string
	ZoneName string
	EntityChildSupport
	ListEntityContainer

	//RoomGroup *EntityListGroup
	VariableSupport
}

func (e *ZoneEntity) UID() string {
	return e.ZoneId
}

func (e *ZoneEntity) NickName() string {
	return e.ZoneName
}

func (e *ZoneEntity) EntityType() basis.EntityType {
	return basis.EntityZone
}

func (e *ZoneEntity) InitEntity() {
	e.EntityChildSupport = *NewEntityChildSupport()
	e.ListEntityContainer = *NewListEntityContainer(0)
	//e.RoomGroup = NewEntityListGroup(EntityRoom)
	e.VariableSupport = *NewVariableSupport(e)
}

//func (e *ZoneEntity) RoomList() []string {
//	return e.RoomGroup.Entities()
//}
//
//func (e *ZoneEntity) ContainRoom(roomId string) bool {
//	return e.RoomGroup.ContainEntity(roomId)
//}
//
//func (e *ZoneEntity) AddRoom(roomId string) error {
//	return e.RoomGroup.Accept(roomId)
//}
//
//func (e *ZoneEntity) RemoveRoom(roomId string) error {
//	return e.RoomGroup.Drop(roomId)
//}
