//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

type IWorld interface {
	IEntity
	IChannelSupport
	Zones() []string
	GetZone(zoneId string) IZone
	GetRooms(zoneId string) []IRoom

	CreateZone(zone *ZoneConfig) (IZone, error)
	CreateRoomAt(zoneId string, room *RoomConfig) (IRoom, error)

	JoinWorld(user *User, zoneId, roomId string) error
	TransferWorld(user *User, zoneId, roomId string) error
	TransferRoom(user *User, roomId string) error

	SetWorldVariables(vars Variables)
}
