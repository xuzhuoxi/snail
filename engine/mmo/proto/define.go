//
//Created by xuzhuoxi
//on 2019-03-18.
//@author xuzhuoxi
//
package proto

import "github.com/xuzhuoxi/snail/engine/mmo/basis"

const (
	NameVar  = "m.v"
	NameChat = "m.c"
)

const IdNone = ""

const (
	IdVarWorld     = "vw"
	IdVarZone      = "vz"
	IdVarRoom      = "vr"
	IdVarUser      = "vu"
	IdVarTeamCrops = "vtc"
	IdVarTeam      = "vt"
	IdVarChannel   = "vc"
)

const (
	IdChatWorld     = "cw"
	IdChatZone      = "cz"
	IdChatRoom      = "cr"
	IdChatUser      = "cu"
	IdChatTeamCrops = "ctc"
	IdChatTeam      = "ct"
	IdChatChannel   = "cc"
)

var et2var = make(map[basis.EntityType]string)

func init() {
	et2var[basis.EntityWorld] = IdVarWorld
	et2var[basis.EntityZone] = IdVarZone
	et2var[basis.EntityRoom] = IdVarRoom
	et2var[basis.EntityUser] = IdVarUser
	et2var[basis.EntityTeamCorps] = IdVarTeamCrops
	et2var[basis.EntityTeam] = IdVarTeam
	et2var[basis.EntityChannel] = IdVarChannel
}

func RegisterIdVar(entityType basis.EntityType, id string) {
	et2var[entityType] = id
}

func GetVarId(entityType basis.EntityType) string {
	if rs, ok := et2var[entityType]; ok {
		return rs
	}
	return IdNone
}
