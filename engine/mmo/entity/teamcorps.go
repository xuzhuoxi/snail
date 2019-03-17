//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package entity

import "github.com/xuzhuoxi/snail/engine/mmo/basis"

func NewITeamCorpsEntity(corpsId string, corpsName string) basis.ITeamCorpsEntity {
	return &TeamCorpsEntity{CorpsId: corpsId, CorpsName: corpsName}
}

func NewTeamCorpsEntity(corpsId string, corpsName string) *TeamCorpsEntity {
	return &TeamCorpsEntity{CorpsId: corpsId, CorpsName: corpsName}
}

type TeamCorpsEntity struct {
	CorpsId   string
	CorpsName string
	EntityChildSupport
	ListEntityContainer

	//TeamGroup *EntityListGroup
	VariableSupport
}

func (e *TeamCorpsEntity) UID() string {
	return e.CorpsId
}

func (e *TeamCorpsEntity) NickName() string {
	return e.CorpsName
}

func (e *TeamCorpsEntity) EntityType() basis.EntityType {
	return basis.EntityTeamCorps
}

func (e *TeamCorpsEntity) InitEntity() {
	e.ListEntityContainer = *NewListEntityContainer(0)
	//e.TeamGroup = NewEntityListGroup(EntityTeam)
	e.VariableSupport = *NewVariableSupport(e)
}

//func (e *TeamCorpsEntity) TeamList() []string {
//	return e.TeamGroup.Entities()
//}
//
//func (e *TeamCorpsEntity) ContainTeam(corpsId string) bool {
//	return e.TeamGroup.ContainEntity(corpsId)
//}
//
//func (e *TeamCorpsEntity) AddTeam(corpsId string) error {
//	return e.TeamGroup.Accept(corpsId)
//}
//
//func (e *TeamCorpsEntity) RemoveTeam(corpsId string) error {
//	return e.TeamGroup.Drop(corpsId)
//}
