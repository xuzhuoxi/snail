//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

import (
	"errors"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

type IZoneEntity interface {
	IEntity
	IEntityOwner
	IInitEntity
	IChannelBehavior
	IRoomGroup
	IVariableSupport
}

type IZoneIndex interface {
	//检查Zone是否存在
	CheckZone(zoneId string) bool
	//获取Zone
	GetZone(zoneId string) IZoneEntity
	//添加一个新Zone到索引中
	AddZone(zone IZoneEntity) error
	//从索引中移除一个Zone
	RemoveZone(zoneId string) (IZoneEntity, error)
	//更新一个新Zone到索引中
	UpdateZone(zone IZoneEntity) error
}

//-----------------------------------------------

type ZoneConfig struct {
}

func NewIZoneEntity(zoneId string, zoneName string) IZoneEntity {
	return &ZoneEntity{ZoneId: zoneId, ZoneName: zoneName}
}

func NewZoneEntity(zoneId string, zoneName string) *ZoneEntity {
	return &ZoneEntity{ZoneId: zoneId, ZoneName: zoneName}
}

type ZoneEntity struct {
	ZoneId   string
	ZoneName string

	roomList []string
	roomMu   sync.RWMutex

	EntityOwnerSupport
	ChannelEntity   *ChannelEntity
	VariableSupport *VariableSupport
}

func (e *ZoneEntity) UID() string {
	return e.ZoneId
}

func (e *ZoneEntity) NickName() string {
	return e.ZoneName
}

func (e *ZoneEntity) InitEntity() {
	e.ChannelEntity = NewChannelEntity(e.ZoneId, e.ZoneName)
	e.VariableSupport = NewVariableSupport()
	e.ChannelEntity.InitEntity()
}

func (e *ZoneEntity) GroupId() string {
	return e.ZoneId
}

func (e *ZoneEntity) GroupName() string {
	return e.ZoneName
}

func (e *ZoneEntity) Rooms() []string {
	return e.roomList
}

func (e *ZoneEntity) CopyRooms() []string {
	return slicex.CopyString(e.roomList)
}

func (e *ZoneEntity) CheckRoom(roomId string) bool {
	_, ok := slicex.IndexString(e.roomList, roomId)
	return ok
}

func (e *ZoneEntity) AddRoom(roomId string) error {
	e.roomMu.Lock()
	defer e.roomMu.Unlock()
	_, ok := slicex.IndexString(e.roomList, roomId)
	if ok {
		return errors.New("AddRoom Repeat At" + roomId)
	}
	e.roomList = append(e.roomList, roomId)
	return nil
}

func (e *ZoneEntity) RemoveRoom(roomId string) error {
	e.roomMu.Lock()
	defer e.roomMu.Unlock()
	index, ok := slicex.IndexString(e.roomList, roomId)
	if !ok {
		return errors.New("RemoveRoom Fail, No Room:" + roomId)
	}
	e.roomList = append(e.roomList[:index], e.roomList[index+1:]...)
	return nil
}

func (e *ZoneEntity) MyChannel() IChannelEntity {
	return e.ChannelEntity
}

func (e *ZoneEntity) TouchChannel(subscriber string) {
	e.ChannelEntity.TouchChannel(subscriber)
}

func (e *ZoneEntity) UnTouchChannel(subscriber string) {
	e.ChannelEntity.UnTouchChannel(subscriber)
}

func (e *ZoneEntity) Broadcast(speaker string, handler func(receiver string)) int {
	return e.ChannelEntity.Broadcast(speaker, handler)
}

func (e *ZoneEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	return e.ChannelEntity.BroadcastSome(speaker, receiver, handler)
}

func (e *ZoneEntity) SetVar(key string, value interface{}) {
	e.VariableSupport.SetVar(key, value)
}

func (e *ZoneEntity) GetVar(key string) interface{} {
	return e.VariableSupport.GetVar(key)
}

func (e *ZoneEntity) CheckVar(key string) bool {
	return e.VariableSupport.CheckVar(key)
}

func (e *ZoneEntity) RemoveVar(key string) {
	e.VariableSupport.RemoveVar(key)
}

//-----------------------------------------------

func NewIZoneIndex() IZoneIndex {
	return &ZoneIndex{zoneMap: make(map[string]IZoneEntity)}
}

func NewZoneIndex() ZoneIndex {
	return ZoneIndex{zoneMap: make(map[string]IZoneEntity)}
}

type ZoneIndex struct {
	zoneMap map[string]IZoneEntity
	mu      sync.RWMutex
}

func (i *ZoneIndex) CheckZone(zoneId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.zoneMap[zoneId]
	return ok
}

func (i *ZoneIndex) GetZone(zoneId string) IZoneEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.zoneMap[zoneId]
}

func (i *ZoneIndex) AddZone(zone IZoneEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == zone {
		return errors.New("ZoneIndex.AddZone Error: zone is nil")
	}
	zoneId := zone.UID()
	if i.CheckZone(zoneId) {
		return errors.New("ZoneIndex.AddZone Error: Zone(" + zoneId + ") Duplicate")
	}
	i.zoneMap[zoneId] = zone
	return nil
}

func (i *ZoneIndex) RemoveZone(zoneId string) (IZoneEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.zoneMap[zoneId]
	if ok {
		delete(i.zoneMap, zoneId)
		return e, nil
	}
	return nil, errors.New("ZoneIndex.RemoveZone Error: No Zone(" + zoneId + ")")
}

func (i *ZoneIndex) UpdateZone(zone IZoneEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == zone {
		return errors.New("ZoneIndex.UpdateZone Error: zone is nil")
	}
	i.zoneMap[zone.UID()] = zone
	return nil
}
