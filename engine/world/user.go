//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

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
}

//-----------------------------------------------

type UserEntity struct {
	Uid  string //用户标识，唯一，内部使用
	Name string //用户名，唯一
	Nick string //用户昵称

	Addr   string
	ZoneId string
	RoomId string

	Pos XYZ

	ChannelSubscriber ChannelSubscriber
	VariableSupport   VariableSupport
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

func (e *UserEntity) InitEntity() {
	e.ChannelSubscriber = NewChannelSubscriber()
	e.VariableSupport = NewVariableSupport()
}

func (e *UserEntity) TouchingChannels() []string {
	return e.ChannelSubscriber.TouchingChannels()
}

func (e *UserEntity) CopyTouchingChannels() []string {
	return e.ChannelSubscriber.CopyTouchingChannels()
}

func (e *UserEntity) TouchChannel(chanId string) error {
	return e.ChannelSubscriber.TouchChannel(chanId)
}

func (e *UserEntity) UnTouchChannel(chanId string) error {
	return e.ChannelSubscriber.UnTouchChannel(chanId)
}

func (e *UserEntity) InChannel(chanId string) bool {
	return e.ChannelSubscriber.InChannel(chanId)
}

func (e *UserEntity) SetVar(key string, value interface{}) {
	e.VariableSupport.SetVar(key, value)
}

func (e *UserEntity) GetVar(key string) interface{} {
	return e.VariableSupport.GetVar(key)
}

func (e *UserEntity) CheckVar(key string) bool {
	return e.VariableSupport.CheckVar(key)
}

func (e *UserEntity) RemoveVar(key string) {
	e.VariableSupport.RemoveVar(key)
}

//-----------------------------------------------

func NewIUserIndex() IUserIndex {
	return &UserIndex{userIdMap: make(map[string]IUserEntity)}
}

func NewUserIndex() UserIndex {
	return UserIndex{userIdMap: make(map[string]IUserEntity)}
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
		return errors.New("AddUser nil!")
	}
	userId := user.UID()
	if i.CheckUser(userId) {
		return errors.New("UserEntity Repeat At :" + userId)
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
	return nil, errors.New("RemoveUser Error: No UserEntity[" + userId + "]")
}
