//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package entity

import "github.com/xuzhuoxi/snail/engine/mmo/basis"

type BroadcastType uint16

//const (
//	BroadcastNone BroadcastType = iota
//	BroadcastWorld
//	BroadcastZone
//	BroadcastRoom
//	BroadcastCorps
//	BroadcastTeam
//	BroadcastUser
//	BroadcastChannel
//)

type IBroadcastEntity interface {
	//消息广播
	Broadcast(speaker string, handler func(entity basis.IUserEntity))
	//消息指定目标广播
	BroadcastSome(speaker string, receiver []string, handler func(entity basis.IUserEntity))
}

func BroadcastRoom(speaker basis.IUserEntity, room basis.IRoomEntity, handler func(entity basis.IUserEntity)) {
	if speaker.OnBlack(room.UID()) {
		return
	}
	room.ForEachChildrenByType(basis.EntityUser, func(child basis.IEntity) {
		if speaker.OnBlack(child.UID()) {
			return
		}
		handler(child.(basis.IUserEntity))
	}, false)
}

func BroadcastTeam(speaker basis.IUserEntity, team basis.ITeamEntity, handler func(entity basis.IUserEntity)) {
	if speaker.OnBlack(team.UID()) {
		return
	}
	team.ForEachChildrenByType(basis.EntityUser, func(child basis.IEntity) {
		if speaker.OnBlack(child.UID()) {
			return
		}
		handler(child.(basis.IUserEntity))
	}, false)
}

func BroadcastContainer(speaker basis.IUserEntity, container basis.IEntity, handler func(entity basis.IUserEntity)) {
	speakerId := speaker.UID()
	if speaker.OnBlack(container.UID()) || speakerId == container.UID() {
		return
	}
	if entityContainer, ok := container.(basis.IEntityContainer); ok {
		var userEntity basis.IUserEntity
		entityContainer.ForEachChildren(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
			if speaker.OnBlack(child.UID()) {
				return false, true
			}
			if basis.EntityUser == child.EntityType() && speakerId != child.UID() {
				userEntity = child.(basis.IUserEntity)
				handler(userEntity)
			}
			return
		})
	}
}
