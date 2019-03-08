//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

import (
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
