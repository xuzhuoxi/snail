//
//Created by xuzhuoxi
//on 2019-03-08.
//@author xuzhuoxi
//
package mmo

type BroadcastType uint16

const (
	BroadcastNone BroadcastType = iota
	BroadcastWorld
	BroadcastZone
	BroadcastRoom
	BroadcastCurrent
)

type IBroadcastManager interface {
	//消息广播
	Broadcast(speaker string, handler func(entity IUserEntity))
	//消息指定目标广播
	BroadcastSome(speaker string, receiver []string, handler func(entity IUserEntity))
}
