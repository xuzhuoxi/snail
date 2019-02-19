//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

type IChannelSupport interface {
	ChannelId() string
	NotifyAll(speaker User, message []byte) int
}

type IChannel interface {
	IEntity
}
