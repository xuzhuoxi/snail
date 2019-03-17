//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

type EntityType uint16

const (
	EntityWorld EntityType = 1 << iota
	EntityZone
	EntityRoom
	EntityUser
	EntityTeamCorps
	EntityTeam
	EntityChannel

	EntityNone EntityType = 0
	EntityAll  EntityType = EntityWorld | EntityZone | EntityRoom | EntityUser | EntityTeamCorps | EntityTeam | EntityChannel
)

func (t EntityType) Match(check EntityType) bool {
	return t&check > 0
}

func (t EntityType) Include(check EntityType) bool {
	return t&check == check
}

type IEntity interface {
	//唯一标识
	UID() string
	//昵称，显示使用
	NickName() string
	//实体类型
	EntityType() EntityType
}

type IInitEntity interface {
	//初始化实体
	InitEntity()
}

type IDestroyEntity interface {
	//释放实体
	DestroyEntity()
}

//世界实体
type IWorldEntity interface {
	IEntity
	IInitEntity

	IEntityContainer
	//IZoneGroup
	IVariableSupport
}

//区域实体
type IZoneEntity interface {
	IEntity
	IEntityChild
	IInitEntity

	IEntityContainer
	//IRoomGroup
	IVariableSupport
}

//兵团实体
type ITeamCorpsEntity interface {
	IEntity
	IEntityChild
	IInitEntity

	IEntityContainer
	//ITeamGroup
	IVariableSupport
}

//房间实体
type IRoomEntity interface {
	IEntity
	IInitEntity
	IEntityChild

	//IUserGroup
	IEntityContainer
	IVariableSupport
}

//队伍实体
type ITeamEntity interface {
	IEntity
	IInitEntity
	IEntityChild

	IEntityContainer
	IVariableSupport
	//ITeamControl
}

//用户实体
type IUserEntity interface {
	IEntity
	IInitEntity
	IDestroyEntity
	IUserSubscriber
	IVariableSupport
	//用户名
	UserName() string

	GetLocation() (idType EntityType, id string)
	SetLocation(idType EntityType, id string)

	GetTeamInfo() (teamId string, corpsId string)
	SetTeam(teamId string)
	SetCorps(corpsId string)

	GetPosition() XYZ
	SetPosition(pos XYZ)
}

//频道实体
type IChannelEntity interface {
	IEntity
	IInitEntity
	IChannelBehavior
}

func EntityEqual(entity1 IEntity, entity2 IEntity) bool {
	return nil != entity1 && nil != entity2 && entity1.UID() == entity2.UID() && entity1.EntityType() == entity2.EntityType() && entity1.NickName() == entity2.NickName()
}
