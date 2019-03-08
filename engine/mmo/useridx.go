//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package mmo

import (
	"errors"
	"sync"
)

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
