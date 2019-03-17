//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewIRoomIndex() basis.IRoomIndex {
	return NewRoomIndex()
}

func NewRoomIndex() *RoomIndex {
	return &RoomIndex{EntityIndex: *NewEntityIndex("RoomIndex", basis.EntityRoom)}
}

type RoomIndex struct {
	EntityIndex
}

func (i *RoomIndex) CheckRoom(roomId string) bool {
	return i.EntityIndex.Check(roomId)
}

func (i *RoomIndex) GetRoom(roomId string) basis.IRoomEntity {
	entity := i.EntityIndex.Get(roomId)
	if nil != entity {
		return entity.(basis.IRoomEntity)
	}
	return nil
}

func (i *RoomIndex) AddRoom(room basis.IRoomEntity) error {
	return i.EntityIndex.Add(room)
}

func (i *RoomIndex) RemoveRoom(roomId string) (basis.IRoomEntity, error) {
	c, err := i.EntityIndex.Remove(roomId)
	if nil != c {
		return c.(basis.IRoomEntity), err
	}
	return nil, err
}

func (i *RoomIndex) UpdateRoom(room basis.IRoomEntity) error {
	return i.EntityIndex.Update(room)
}
