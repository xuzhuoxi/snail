//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"errors"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

func NewIUserIndex() basis.IUserIndex {
	return &UserIndex{userIdMap: make(map[string]basis.IUserEntity)}
}

func NewUserIndex() *UserIndex {
	return &UserIndex{userIdMap: make(map[string]basis.IUserEntity)}
}

type UserIndex struct {
	userIdMap map[string]basis.IUserEntity
	mu        sync.RWMutex
}

func (i *UserIndex) CheckUser(userId string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.userIdMap[userId]
	return ok
}

func (i *UserIndex) GetUser(userId string) basis.IUserEntity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.userIdMap[userId]
}

func (i *UserIndex) AddUser(user basis.IUserEntity) error {
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

func (i *UserIndex) RemoveUser(userId string) (basis.IUserEntity, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.userIdMap[userId]
	if ok {
		delete(i.userIdMap, userId)
		return e, nil
	}
	return nil, errors.New("UserIndex.RemoveUser Error: No User(" + userId + ")")
}

func (i *UserIndex) UpdateUser(user basis.IUserEntity) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if nil == user {
		return errors.New("UserIndex.UpdateUser Error: user is nil")
	}
	i.userIdMap[user.UID()] = user
	return nil
}
