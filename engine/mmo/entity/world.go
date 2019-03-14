//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package entity

import "github.com/xuzhuoxi/snail/engine/mmo/basis"

func CreateWorldEntity(worldId string, worldName string) basis.IWorldEntity {
	return &WorldEntity{WorldId: worldId, WorldName: worldName}
}

type WorldEntity struct {
	WorldId   string
	WorldName string
	ListEntityContainer

	//ZoneGroup *EntityListGroup
	VariableSupport
}

func (w *WorldEntity) UID() string {
	return w.WorldId
}

func (w *WorldEntity) NickName() string {
	return w.WorldName
}

func (w *WorldEntity) EntityType() basis.EntityType {
	return basis.EntityWorld
}

func (w *WorldEntity) InitEntity() {
	w.ListEntityContainer = *NewListEntityContainer(0)
	//w.ZoneGroup = NewEntityListGroup(EntityZone)
	w.VariableSupport = *NewVariableSupport()
}

//func (w *WorldEntity) ZoneList() []string {
//	return w.ZoneGroup.Entities()
//}
//
//func (w *WorldEntity) ContainZone(zoneId string) bool {
//	return w.ZoneGroup.ContainEntity(zoneId)
//}
//
//func (w *WorldEntity) AddZone(zoneId string) error {
//	return w.ZoneGroup.Accept(zoneId)
//}
//
//func (w *WorldEntity) RemoveZone(zoneId string) error {
//	return w.ZoneGroup.Drop(zoneId)
//}
