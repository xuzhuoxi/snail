//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

import "github.com/pkg/errors"

type EntityType int

const (
	EntityNone EntityType = iota
	EntityChannel
	EntityUser
	EntityRoom
	EntityZone
	EntityWorld

	EntityMax
)

type IEntity interface {
	UID() string
}

type IEntityIndex interface {
	CheckZone(zoneId string) bool
	GetZone(zoneId string) IZone
	AddZone(zone IZone) error
	RemoveZone(zoneId string) (IZone, error)

	CheckRoom(roomId string) bool
	GetRoom(roomId string) IRoom
	AddRoom(room IRoom) error
	RemoveRoom(roomId string) (IRoom, error)

	CheckUser(userId string) bool
	GetUser(userId string) IUser
	AddUser(user IUser) error
	RemoveUser(userId string) (IUser, error)

	CheckChannel(chanId string) bool
	GetChannel(chanId string) IChannel
	AddChannel(channel IChannel) error
	RemoveChannel(chanId string) (IChannel, error)
}

func NewEntityIndex() IEntityIndex {
	zoneMap := make(map[string]IZone)
	roomMap := make(map[string]IRoom)
	userMap := make(map[string]IUser)
	chanMap := make(map[string]IChannel)
	return &EntityIndex{zoneMap: zoneMap, roomMap: roomMap, userMap: userMap, chanMap: chanMap}
}

type EntityIndex struct {
	zoneMap map[string]IZone
	roomMap map[string]IRoom
	userMap map[string]IUser
	chanMap map[string]IChannel
}

func (i *EntityIndex) CheckZone(zoneId string) bool {
	_, ok := i.zoneMap[zoneId]
	return ok
}

func (i *EntityIndex) GetZone(zoneId string) IZone {
	return i.zoneMap[zoneId]
}

func (i *EntityIndex) AddZone(zone IZone) error {
	if nil == zone {
		return errors.New("AddZone nil!")
	}
	zoneId := zone.UID()
	if i.CheckZone(zoneId) {
		return errors.New("Zone Repeat At :" + zoneId)
	}
	i.zoneMap[zoneId] = zone
	return nil
}

func (i *EntityIndex) RemoveZone(zoneId string) (IZone, error) {
	e, ok := i.zoneMap[zoneId]
	if ok {
		delete(i.zoneMap, zoneId)
		return e, nil
	}
	return nil, errors.New("RemoveZone Error: No Zone[" + zoneId + "]")
}

func (i *EntityIndex) CheckRoom(roomId string) bool {
	_, ok := i.roomMap[roomId]
	return ok
}

func (i *EntityIndex) GetRoom(roomId string) IRoom {
	return i.roomMap[roomId]
}

func (i *EntityIndex) AddRoom(room IRoom) error {
	if nil == room {
		return errors.New("AddRoom nil!")
	}
	roomId := room.UID()
	if i.CheckRoom(roomId) {
		return errors.New("Room Repeat At :" + roomId)
	}
	i.roomMap[roomId] = room
	return nil
}

func (i *EntityIndex) RemoveRoom(roomId string) (IRoom, error) {
	e, ok := i.roomMap[roomId]
	if ok {
		delete(i.roomMap, roomId)
		return e, nil
	}
	return nil, errors.New("RemoveRoom Error: No Room[" + roomId + "]")
}

func (i *EntityIndex) CheckUser(userId string) bool {
	_, ok := i.userMap[userId]
	return ok
}

func (i *EntityIndex) GetUser(userId string) IUser {
	return i.userMap[userId]
}

func (i *EntityIndex) AddUser(user IUser) error {
	if nil == user {
		return errors.New("AddUser nil!")
	}
	userId := user.UID()
	if i.CheckUser(userId) {
		return errors.New("User Repeat At :" + userId)
	}
	i.userMap[userId] = user
	return nil
}

func (i *EntityIndex) RemoveUser(userId string) (IUser, error) {
	e, ok := i.userMap[userId]
	if ok {
		delete(i.userMap, userId)
		return e, nil
	}
	return nil, errors.New("RemoveUser Error: No User[" + userId + "]")
}

func (i *EntityIndex) CheckChannel(chanId string) bool {
	_, ok := i.chanMap[chanId]
	return ok
}

func (i *EntityIndex) GetChannel(chanId string) IChannel {
	return i.chanMap[chanId]
}

func (i *EntityIndex) AddChannel(channel IChannel) error {
	if nil == channel {
		return errors.New("AddChannel nil!")
	}
	chanId := channel.UID()
	if i.CheckChannel(chanId) {
		return errors.New("Channel Repeat At :" + chanId)
	}
	i.chanMap[chanId] = channel
	return nil
}

func (i *EntityIndex) RemoveChannel(chanId string) (IChannel, error) {
	e, ok := i.chanMap[chanId]
	if ok {
		delete(i.chanMap, chanId)
		return e, nil
	}
	return nil, errors.New("RemoveChannel Error: No Channel[" + chanId + "]")
}
