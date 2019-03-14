//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

const (
	EventUserJoinRoom  = "EventUserJoinRoom"
	EventUserLeaveRoom = "EventUserLeaveRoom"
)

//黑名单
type IUserBlackList interface {
	//通信黑名单，返回原始切片，如果要修改的，请先copy
	Blacks() []string
	//增加黑名单
	AddBlack(targetId string) error
	//移除黑名单
	RemoveBlack(targetId string) error
	//处于
	OnBlack(targetId string) bool
}

//黑名单
type IUserWhiteList interface {
	//通信白名单，返回原始切片，如果要修改的，请先copy
	Whites() []string
	//增加白名单
	AddWhite(targetId string) error
	//移除白名单
	RemoveWhite(targetId string) error
	//处于
	OnWhite(targetId string) bool
}

//参与者
type IUserSubscriber interface {
	IUserWhiteList
	IUserBlackList
	//处于激活
	OnActive(targetId string) bool
}
