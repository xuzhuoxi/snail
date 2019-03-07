//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

type ChannelType uint16

const (
	//无效
	None ChannelType = iota
	//状态
	StatusChannel
	//聊天
	ChatChannel
	//事件
	EventChannel
)

//频道实体
type IChannelEntity interface {
	IEntity
	IInitEntity
	IChannelBehavior
}

//频道行为
type IChannelBehavior interface {
	MyChannel() IChannelEntity
	//订阅频道
	TouchChannel(subscriber string)
	//取消频道订阅
	UnTouchChannel(subscriber string)
	//消息广播
	Broadcast(speaker string, handler func(receiver string)) int
	//消息指定目标广播
	BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int
}

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

//-----------------------------------------------

func NewIChannelEntity(chanId string, chanName string) IChannelEntity {
	return &ChannelEntity{ChanId: chanId, ChanName: chanName}
}

func NewChannelEntity(chanId string, chanName string) *ChannelEntity {
	return &ChannelEntity{ChanId: chanId, ChanName: chanName}
}

type ChannelEntity struct {
	ChanId     string
	ChanName   string
	Subscriber []string
	Mu         sync.RWMutex
}

func (c *ChannelEntity) UID() string {
	return c.ChanId
}

func (c *ChannelEntity) NickName() string {
	return c.ChanName
}

func (c *ChannelEntity) EntityType() EntityType {
	return EntityChannel
}

func (c *ChannelEntity) InitEntity() {
}

func (c *ChannelEntity) MyChannel() IChannelEntity {
	return c
}

func (c *ChannelEntity) TouchChannel(subscriber string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if c.hasSubscriber(subscriber) {
		return
	}
	c.Subscriber = append(c.Subscriber, subscriber)
}

func (c *ChannelEntity) UnTouchChannel(subscriber string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	index, ok := slicex.IndexString(c.Subscriber, subscriber)
	if !ok {
		return
	}
	c.Subscriber = append(c.Subscriber[:index], c.Subscriber[index+1:]...)
}

func (c *ChannelEntity) Broadcast(speaker string, handler func(receiver string)) int {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	rs := len(c.Subscriber)
	for _, r := range c.Subscriber {
		if r == speaker {
			continue
		}
		handler(r)
	}
	return rs - 1
}

func (c *ChannelEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	count := 0
	for _, v := range c.Subscriber {
		if _, ok := slicex.IndexString(receiver, v); ok && speaker != v {
			handler(v)
			count++
		}
	}
	return count
}

func (c *ChannelEntity) hasSubscriber(subscriber string) bool {
	_, ok := slicex.IndexString(c.Subscriber, subscriber)
	return ok
}

//-----------------------------------------------

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
