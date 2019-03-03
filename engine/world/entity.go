//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

type EntityType int

const (
	EntityNone EntityType = iota
	EntityChannel
	EntityUser
	EntityRoom
	EntityZone
	EntityWorld

	EntityMax
)

type IEntity interface {
	//唯一标识
	UID() string
	//昵称，显示使用
	NickName() string
}

type IInitEntity interface {
	//初始化实体
	InitEntity()
}
