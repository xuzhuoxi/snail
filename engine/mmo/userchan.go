//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package mmo

type IChannelSubscriberBlack interface {
	//通信频道黑名单，返回原始切片，如果要修改的，请先copy
	BlackChannels() []string
	//增加频道黑名单
	AddBlackChannel(chanId string) error
	//移除频道黑名单
	RemoveBlackChannel(chanId string) error
	//处于频道
	OnBlackChannel(chanId string) bool
}

type IChannelSubscriberWhite interface {
	//订阅的通信频道，返回原始切片，如果要修改的，请先copy
	WhiteChannels() []string
	//订阅频道
	AddWhiteChannel(chanId string) error
	//取消频道订阅
	RemoveWhiteChannel(chanId string) error
	//处于频道
	OnWhiteChannel(chanId string) bool
}

//频道参与者
type IChannelSubscriber interface {
	IChannelSubscriberWhite
	IChannelSubscriberBlack
	//处于激活频道
	OnChannel(chanId string) bool
}

//-----------------------------------------------

func NewIChannelSubscriber() IChannelSubscriber {
	return &ChannelSubscriber{ChannelSubscriberBlack: *NewChannelSubscriberBlack(), ChannelSubscriberWhite: *NewChannelSubscriberWhite()}
}

func NewIChannelSubscriberWhite() IChannelSubscriberWhite {
	return NewChannelSubscriberWhite()
}
func NewIChannelSubscriberBlack() IChannelSubscriberBlack {
	return NewChannelSubscriberBlack()
}

func NewChannelSubscriber() *ChannelSubscriber {
	return &ChannelSubscriber{ChannelSubscriberBlack: *NewChannelSubscriberBlack(), ChannelSubscriberWhite: *NewChannelSubscriberWhite()}
}

func NewChannelSubscriberWhite() *ChannelSubscriberWhite {
	return &ChannelSubscriberWhite{whiteGroup: NewEntityListGroup(EntityUser)}
}
func NewChannelSubscriberBlack() *ChannelSubscriberBlack {
	return &ChannelSubscriberBlack{blackGroup: NewEntityListGroup(EntityUser)}
}

type ChannelSubscriber struct {
	ChannelSubscriberBlack
	ChannelSubscriberWhite
}

func (c *ChannelSubscriber) OnChannel(chanId string) bool {
	return c.OnWhiteChannel(chanId) && !c.OnBlackChannel(chanId)
}

type ChannelSubscriberWhite struct {
	whiteGroup IEntityGroup
}

func (c *ChannelSubscriberWhite) WhiteChannels() []string {
	return c.whiteGroup.Entities()
}

func (c *ChannelSubscriberWhite) AddWhiteChannel(chanId string) error {
	return c.whiteGroup.Accept(chanId)
}

func (c *ChannelSubscriberWhite) RemoveWhiteChannel(chanId string) error {
	return c.whiteGroup.Drop(chanId)
}

func (c *ChannelSubscriberWhite) OnWhiteChannel(chanId string) bool {
	return c.whiteGroup.ContainEntity(chanId)
}

type ChannelSubscriberBlack struct {
	blackGroup IEntityGroup
}

func (c *ChannelSubscriberBlack) BlackChannels() []string {
	return c.blackGroup.Entities()
}

func (c *ChannelSubscriberBlack) AddBlackChannel(chanId string) error {
	return c.blackGroup.Accept(chanId)
}

func (c *ChannelSubscriberBlack) RemoveBlackChannel(chanId string) error {
	return c.blackGroup.Drop(chanId)
}

func (c *ChannelSubscriberBlack) OnBlackChannel(chanId string) bool {
	return c.blackGroup.ContainEntity(chanId)
}
