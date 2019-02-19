//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

type IRoom interface {
	IEntity
	IChannelSupport
	RoomName() string

	SetRoomVariables(vars Variables)
}

type IRoomGroup interface {
	Rooms() []string
}

type RoomConfig struct {
	Id        string
	Name      string
	Private   bool
	MaxMember int
}
