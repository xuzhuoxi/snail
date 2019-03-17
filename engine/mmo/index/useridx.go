//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package index

import (
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

func NewIUserIndex() basis.IUserIndex {
	return NewUserIndex()
}

func NewUserIndex() *UserIndex {
	return &UserIndex{EntityIndex: *NewEntityIndex("UserIndex", basis.EntityUser)}
}

type UserIndex struct {
	EntityIndex
}

func (i *UserIndex) CheckUser(userId string) bool {
	return i.EntityIndex.Check(userId)
}

func (i *UserIndex) GetUser(userId string) basis.IUserEntity {
	entity := i.EntityIndex.Get(userId)
	if nil != entity {
		return entity.(basis.IUserEntity)
	}
	return nil
}

func (i *UserIndex) AddUser(user basis.IUserEntity) error {
	return i.EntityIndex.Add(user)
}

func (i *UserIndex) RemoveUser(userId string) (basis.IUserEntity, error) {
	c, err := i.EntityIndex.Remove(userId)
	if nil != c {
		return c.(basis.IUserEntity), err
	}
	return nil, err
}

func (i *UserIndex) UpdateUser(user basis.IUserEntity) error {
	return i.EntityIndex.Update(user)
}
