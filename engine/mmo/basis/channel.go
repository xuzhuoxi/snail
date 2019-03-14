//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

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
