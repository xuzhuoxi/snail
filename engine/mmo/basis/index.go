//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

type IEntityIndex interface {
	EntityType() EntityType
	//检查存在
	Check(id string) bool
	//获取one
	Get(id string) IEntity
	//添加
	Add(entity IEntity) error
	//从索引中移除
	Remove(id string) (IEntity, error)
	//更新
	Update(entity IEntity) error
}

type IWorldIndex interface {
	IEntityIndex
	//检查World是否存在
	CheckWorld(worldId string) bool
	//获取World
	GetWorld(worldId string) IWorldEntity
	//添加一个新World到索引中
	AddWorld(world IWorldEntity) error
	//从索引中移除一个World
	RemoveWorld(worldId string) (IWorldEntity, error)
	//更新一个新World到索引中
	UpdateWorld(zone IWorldEntity) error
}

type IZoneIndex interface {
	IEntityIndex
	//检查Zone是否存在
	CheckZone(zoneId string) bool
	//获取Zone
	GetZone(zoneId string) IZoneEntity
	//添加一个新Zone到索引中
	AddZone(zone IZoneEntity) error
	//从索引中移除一个Zone
	RemoveZone(zoneId string) (IZoneEntity, error)
	//更新一个新Zone到索引中
	UpdateZone(zone IZoneEntity) error
}

//房间索引
type IRoomIndex interface {
	IEntityIndex
	//检查Room是否存在
	CheckRoom(roomId string) bool
	//获取Room
	GetRoom(roomId string) IRoomEntity
	//添加一个新Room到索引中
	AddRoom(room IRoomEntity) error
	//从索引中移除一个Room
	RemoveRoom(roomId string) (IRoomEntity, error)
	//从索引中更新一个Room
	UpdateRoom(room IRoomEntity) error
}

type ITeamCorpsIndex interface {
	IEntityIndex
	//检查Corps是否存在
	CheckCorps(corpsId string) bool
	//获取Corps
	GetCorps(corpsId string) ITeamCorpsEntity
	//添加一个新Corps到索引中
	AddCorps(corps ITeamCorpsEntity) error
	//从索引中移除一个Corps
	RemoveCorps(corpsId string) (ITeamCorpsEntity, error)
	//更新一个新Corps到索引中
	UpdateCorps(corps ITeamCorpsEntity) error
}

//队伍索引
type ITeamIndex interface {
	IEntityIndex
	//检查Team是否存在
	CheckTeam(teamId string) bool
	//获取Team
	GetTeam(teamId string) ITeamEntity
	//添加一个新Team到索引中
	AddTeam(team ITeamEntity) error
	//从索引中移除一个Team
	RemoveTeam(teamId string) (ITeamEntity, error)
	//从索引中更新一个Team
	UpdateTeam(team ITeamEntity) error
}

//玩家索引
type IUserIndex interface {
	IEntityIndex
	//检查User是否存在
	CheckUser(userId string) bool
	//获取User
	GetUser(userId string) IUserEntity
	//添加一个新User到索引中
	AddUser(user IUserEntity) error
	//从索引中移除一个User
	RemoveUser(userId string) (IUserEntity, error)
	//从索引中更新一个User
	UpdateUser(user IUserEntity) error
}

//频道索引
type IChannelIndex interface {
	IEntityIndex
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
