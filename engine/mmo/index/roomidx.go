//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"errors"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIRoomIndex() basis.IRoomIndex {
	return &RoomIndex{roomMap: make(map[string]basis.IRoomEntity)}
}

func NewRoomIndex() RoomIndex {
	return RoomIndex{roomMap: make(map[string]basis.IRoomEntity)}
}

type RoomIndex struct {
	roomMap map[string]basis.IRoomEntity
	mu      sync.RWMutex
}

func (i *RoomIndex) CheckRoom(roomId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.checkRoom(roomId)
}

func (i *RoomIndex) checkRoom(roomId string) bool {
	_, ok := i.roomMap[roomId]
	return ok
}

func (i *RoomIndex) GetRoom(roomId string) basis.IRoomEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.roomMap[roomId]
}

func (i *RoomIndex) AddRoom(room basis.IRoomEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == room {
		return errors.New("RoomIndex.AddRoom Error: room is nil")
	}
	roomId := room.UID()
	if i.checkRoom(roomId) {
		return errors.New("RoomIndex.AddRoom Error: Room(" + roomId + ") Duplicate")
	}
	i.roomMap[roomId] = room
	return nil
}

func (i *RoomIndex) RemoveRoom(roomId string) (basis.IRoomEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.roomMap[roomId]
	if ok {
		delete(i.roomMap, roomId)
		return e, nil
	}
	return nil, errors.New("RoomIndex.RemoveRoom Error: No Room(" + roomId + ")")
}

func (i *RoomIndex) UpdateRoom(room basis.IRoomEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == room {
		return errors.New("RoomIndex.UpdateRoom Error: room is nil")
	}
	i.roomMap[room.UID()] = room
	return nil
}
