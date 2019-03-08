//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"sync"
)

//频道索引
type IChannelIndex interface {
	//检查Channel是否存在
	CheckChannel(chanId string) bool
	//获取Channel
	GetChannel(chanId string) IChannelEntity
	//从索引中增加一个Channel
	AddChannel(channel IChannelEntity) error
	//从索引中移除一个Channel
	RemoveChannel(chanId string) (IChannelEntity, error)
	//从索引中更新一个Channel
	UpdateChannel(channel IChannelEntity) error
}

func NewIChannelIndex() IChannelIndex {
	return &ChannelIndex{chanMap: make(map[string]IChannelEntity)}
}

func NewChannelIndex() ChannelIndex {
	return ChannelIndex{chanMap: make(map[string]IChannelEntity)}
}

type ChannelIndex struct {
	chanMap map[string]IChannelEntity
	mu      sync.RWMutex
}

func (i *ChannelIndex) CheckChannel(chanId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.chanMap[chanId]
	return ok
}

func (i *ChannelIndex) GetChannel(chanId string) IChannelEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.chanMap[chanId]
}

func (i *ChannelIndex) AddChannel(channel IChannelEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == channel {
		return errors.New("ChannelIndex.AddChannel Error: channel is nil")
	}
	chanId := channel.UID()
	if i.CheckChannel(chanId) {
		return errors.New("ChannelIndex.AddChannel Error: Channel(" + chanId + ") Duplicate")
	}
	i.chanMap[chanId] = channel
	return nil
}

func (i *ChannelIndex) RemoveChannel(chanId string) (IChannelEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if e, ok := i.chanMap[chanId]; ok {
		delete(i.chanMap, chanId)
		return e, nil
	}
	return nil, errors.New("ChannelIndex.RemoveChannel Error: No Channel(" + chanId + ")")
}

func (i *ChannelIndex) UpdateChannel(channel IChannelEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == channel {
		return errors.New("ChannelIndex.UpdateChannel Error: Channel is nil")
	}
	i.chanMap[channel.UID()] = channel
	return nil
}
