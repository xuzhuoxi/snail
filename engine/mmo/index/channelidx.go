//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewIChannelIndex() basis.IChannelIndex {
	return NewChannelIndex()
}

func NewChannelIndex() *ChannelIndex {
	return &ChannelIndex{EntityIndex: *NewEntityIndex("ChannelIndex", basis.EntityChannel)}
}

type ChannelIndex struct {
	EntityIndex
}

func (i *ChannelIndex) CheckChannel(chanId string) bool {
	return i.EntityIndex.Check(chanId)
}

func (i *ChannelIndex) GetChannel(chanId string) basis.IChannelEntity {
	entity := i.EntityIndex.Get(chanId)
	if nil != entity {
		return entity.(basis.IChannelEntity)
	}
	return nil
}

func (i *ChannelIndex) AddChannel(channel basis.IChannelEntity) error {
	return i.EntityIndex.Add(channel)
}

func (i *ChannelIndex) RemoveChannel(chanId string) (basis.IChannelEntity, error) {
	c, err := i.EntityIndex.Remove(chanId)
	if nil != c {
		return c.(basis.IChannelEntity), err
	}
	return nil, err
}

func (i *ChannelIndex) UpdateChannel(channel basis.IChannelEntity) error {
	return i.EntityIndex.Update(channel)
}
