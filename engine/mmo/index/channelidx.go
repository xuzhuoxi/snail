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

func NewIChannelIndex() basis.IChannelIndex {
	return &ChannelIndex{chanMap: make(map[string]basis.IChannelEntity)}
}

func NewChannelIndex() ChannelIndex {
	return ChannelIndex{chanMap: make(map[string]basis.IChannelEntity)}
}

type ChannelIndex struct {
	chanMap map[string]basis.IChannelEntity
	mu      sync.RWMutex
}

func (i *ChannelIndex) CheckChannel(chanId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.checkChannel(chanId)
}

func (i *ChannelIndex) checkChannel(chanId string) bool {
	_, ok := i.chanMap[chanId]
	return ok
}

func (i *ChannelIndex) GetChannel(chanId string) basis.IChannelEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.chanMap[chanId]
}

func (i *ChannelIndex) AddChannel(channel basis.IChannelEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == channel {
		return errors.New("ChannelIndex.AddChannel Error: channel is nil")
	}
	chanId := channel.UID()
	if i.checkChannel(chanId) {
		return errors.New("ChannelIndex.AddChannel Error: Channel(" + chanId + ") Duplicate")
	}
	i.chanMap[chanId] = channel
	return nil
}

func (i *ChannelIndex) RemoveChannel(chanId string) (basis.IChannelEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if e, ok := i.chanMap[chanId]; ok {
		delete(i.chanMap, chanId)
		return e, nil
	}
	return nil, errors.New("ChannelIndex.RemoveBlack Error: No Channel(" + chanId + ")")
}

func (i *ChannelIndex) UpdateChannel(channel basis.IChannelEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == channel {
		return errors.New("ChannelIndex.UpdateChannel Error: Channel is nil")
	}
	i.chanMap[channel.UID()] = channel
	return nil
}
