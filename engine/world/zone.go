//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

type IZone interface {
	IEntity
	IChannelSupport
	Rooms() []string
	GetRoom(roomId string) IRoom
	CreateRoom(room *RoomConfig) (IRoom, error)

	SetZoneVariables(vars Variables)
}

type ZoneConfig struct {
}
