//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package mmo

import (
	"github.com/pkg/errors"
	"sync"
)

//用户实体
type IUserEntity interface {
	IEntity
	IInitEntity
	IChannelSubscriber
	IVariableSupport

	//用户名
	UserName() string

	GetLocation() (zoneId string, roomId string)
	SetZone(zoneId string, roomId string)
	SetRoom(roomId string)
	GetPosition() XYZ
	SetPosition(pos XYZ)
}

//玩家索引
type IUserIndex interface {
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

//-----------------------------------------------

type UserEntity struct {
	Uid  string //用户标识，唯一，内部使用
	Name string //用户名，唯一
	Nick string //用户昵称

	Addr   string
	ZoneId string
	RoomId string
	locMu  sync.RWMutex

	Pos   XYZ
	posMu sync.RWMutex

	ChannelSubscriber
	VariableSupport
}

func (e *UserEntity) UID() string {
	return e.Uid
}

func (e *UserEntity) UserName() string {
	return e.Name
}

func (e *UserEntity) NickName() string {
	return e.Nick
}

func (e *UserEntity) EntityType() EntityType {
	return EntityUser
}

func (e *UserEntity) InitEntity() {
	e.ChannelSubscriber = *NewChannelSubscriber()
	e.VariableSupport = *NewVariableSupport()
}

func (e *UserEntity) GetLocation() (zoneId string, roomId string) {
	e.locMu.RLock()
	defer e.locMu.RUnlock()
	return e.ZoneId, e.RoomId
}

func (e *UserEntity) SetZone(zoneId string, roomId string) {
	e.locMu.Lock()
	defer e.locMu.Unlock()
	if zoneId != e.ZoneId {
		e.ZoneId = zoneId
		e.ChannelSubscriber.AddWhiteChannel(zoneId)
	}
	if roomId != e.RoomId {
		e.RoomId = roomId
		e.ChannelSubscriber.AddWhiteChannel(roomId)
	}
}

func (e *UserEntity) SetRoom(roomId string) {
	e.locMu.Lock()
	defer e.locMu.Unlock()
	if roomId == e.RoomId {
		return
	}
	e.RoomId = roomId
	e.ChannelSubscriber.AddWhiteChannel(roomId)
}

func (e *UserEntity) GetPosition() XYZ {
	e.posMu.RLock()
	defer e.posMu.RUnlock()
	return e.Pos
}

func (e *UserEntity) SetPosition(pos XYZ) {
	e.posMu.Lock()
	defer e.posMu.Unlock()
	e.Pos = pos
}

//-----------------------------------------------

func NewIUserIndex() IUserIndex {
	return &UserIndex{userIdMap: make(map[string]IUserEntity)}
}

func NewUserIndex() *UserIndex {
	return &UserIndex{userIdMap: make(map[string]IUserEntity)}
}

type UserIndex struct {
	userIdMap map[string]IUserEntity
	mu        sync.RWMutex
}

func (i *UserIndex) CheckUser(userId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.userIdMap[userId]
	return ok
}

func (i *UserIndex) GetUser(userId string) IUserEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.userIdMap[userId]
}

func (i *UserIndex) AddUser(user IUserEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == user {
		return errors.New("UserIndex.AddUser Error: user is nil")
	}
	userId := user.UID()
	if i.CheckUser(userId) {
		return errors.New("UserIndex.AddUser Error: UserId(" + userId + ") Duplicate")
	}
	i.userIdMap[userId] = user
	return nil
}

func (i *UserIndex) RemoveUser(userId string) (IUserEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.userIdMap[userId]
	if ok {
		delete(i.userIdMap, userId)
		return e, nil
	}
	return nil, errors.New("UserIndex.RemoveUser Error: No User(" + userId + ")")
}

func (i *UserIndex) UpdateUser(user IUserEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == user {
		return errors.New("UserIndex.UpdateUser Error: user is nil")
	}
	i.userIdMap[user.UID()] = user
	return nil
}
