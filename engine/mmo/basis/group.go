//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

type IZoneGroup interface {
	//区域列表
	ZoneList() []string
	//检查区域存在性
	ContainZone(zoneId string) bool
	//添加区域
	AddZone(zoneId string) error
	//移除区域
	RemoveZone(zoneId string) error
}

type ITeamGroup interface {
	//队伍列表
	TeamList() []string
	//检查队伍存在性
	ContainTeam(roomId string) bool
	//添加房间
	AddTeam(roomId string) error
	//移除房间
	RemoveTeam(roomId string) error
}

type IRoomGroup interface {
	//房间列表
	RoomList() []string
	//检查房间存在性
	ContainRoom(roomId string) bool
	//添加房间
	AddRoom(roomId string) error
	//移除房间
	RemoveRoom(roomId string) error
}

type IUserGroup interface {
	//用户列表
	UserList() []string
	//检查用户
	ContainUser(userId string) bool
	//加入用户,进行唯一性检查
	AcceptUser(userId string) error
	//从组中移除用户
	DropUser(userId string) error
}

//组
type IEntityGroup interface {
	//接纳实体的类型
	EntityType() EntityType
	//最大实例数
	MaxLen() int
	//实体数量
	Len() int
	//实体已满
	IsFull() bool

	//包含实体id
	Entities() []string
	//包含实体id
	CopyEntities() []string
	//检查实体是否属于当前组
	ContainEntity(entityId string) bool

	//加入实体到组,进行唯一性检查
	Accept(entity string) error
	//加入实体到组,进行唯一性检查
	AcceptMulti(entityId []string) (count int, err error)
	//从组中移除实体
	Drop(entityId string) error
	//从组中移除实体
	DropMulti(entityId []string) (count int, err error)
}
