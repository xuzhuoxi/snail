//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

type IWorldEntity interface {
	IEntity
	IInitEntity
	IChannelBehavior
	IVariableSupport

	//添加区域
	AddZone(zoneId string) error
	//移除区域
	RemoveZone(zoneId string) error
	//检查区域存在性
	ContainZone(zoneId string) bool
	//区域列表
	ZoneList() []string
}

//-----------------------------------------------

func CreateWorldEntity() IWorldEntity {
	return &WorldEntity{}
}

type WorldEntity struct {
	WorldId   string
	WorldName string
	ZoneGroup *EntityListGroup

	VariableSupport *VariableSupport
	ChannelEntity   *ChannelEntity
}

func (w *WorldEntity) UID() string {
	return w.WorldId
}

func (w *WorldEntity) NickName() string {
	return w.WorldName
}

func (w *WorldEntity) EntityType() EntityType {
	return EntityWorld
}

func (w *WorldEntity) InitEntity() {
	w.ZoneGroup = NewEntityListGroup(w.WorldId, w.WorldName, EntityZone)
	w.VariableSupport = NewVariableSupport()
	w.ChannelEntity = NewChannelEntity(w.WorldId, w.WorldName)
	w.ChannelEntity.InitEntity()
}

func (w *WorldEntity) ChannelId() string {
	return w.ChannelEntity.ChannelId()
}

func (w *WorldEntity) MyChannel() IChannelEntity {
	return w.ChannelEntity
}

func (w *WorldEntity) TouchChannel(subscriber string) {
	w.ChannelEntity.TouchChannel(subscriber)
}

func (w *WorldEntity) UnTouchChannel(subscriber string) {
	w.ChannelEntity.UnTouchChannel(subscriber)
}

func (w *WorldEntity) Broadcast(speaker string, handler func(receiver string)) int {
	return w.ChannelEntity.Broadcast(speaker, handler)
}

func (w *WorldEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	return w.ChannelEntity.BroadcastSome(speaker, receiver, handler)
}

func (w *WorldEntity) SetVar(key string, value interface{}) {
	w.VariableSupport.SetVar(key, value)
}

func (w *WorldEntity) GetVar(key string) interface{} {
	return w.VariableSupport.GetVar(key)
}

func (w *WorldEntity) CheckVar(key string) bool {
	return w.VariableSupport.CheckVar(key)
}

func (w *WorldEntity) RemoveVar(key string) {
	w.VariableSupport.RemoveVar(key)
}

func (w *WorldEntity) AddZone(zoneId string) error {
	return w.ZoneGroup.AppendEntity(zoneId)
}

func (w *WorldEntity) RemoveZone(zoneId string) error {
	return w.ZoneGroup.RemoveEntity(zoneId)
}

func (w *WorldEntity) ContainZone(zoneId string) bool {
	return w.ZoneGroup.CheckEntity(zoneId)
}

func (w *WorldEntity) ZoneList() []string {
	return w.ZoneGroup.Entities()
}
