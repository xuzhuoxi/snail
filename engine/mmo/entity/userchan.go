//
//Created by xuzhuoxi
//on 2019-03-07.
//@author xuzhuoxi
//
package entity

import "github.com/xuzhuoxi/snail/engine/mmo/basis"

func NewIUserSubscriber() basis.IUserSubscriber {
	return &UserSubscriber{UserBlackList: *NewUserBlackList(), UserWhiteList: *NewUserWhiteList()}
}

func NewIUserWhiteList() basis.IUserWhiteList {
	return NewUserWhiteList()
}
func NewIUserBlackList() basis.IUserBlackList {
	return NewUserBlackList()
}

func NewUserSubscriber() *UserSubscriber {
	return &UserSubscriber{UserBlackList: *NewUserBlackList(), UserWhiteList: *NewUserWhiteList()}
}

func NewUserWhiteList() *UserWhiteList {
	return &UserWhiteList{whiteGroup: NewEntityListGroup(basis.EntityUser)}
}
func NewUserBlackList() *UserBlackList {
	return &UserBlackList{blackGroup: NewEntityListGroup(basis.EntityUser)}
}

type UserSubscriber struct {
	UserBlackList
	UserWhiteList
}

func (c *UserSubscriber) OnActive(targetId string) bool {
	return c.OnWhite(targetId) && !c.OnBlack(targetId)
}

type UserWhiteList struct {
	whiteGroup basis.IEntityGroup
}

func (c *UserWhiteList) Whites() []string {
	return c.whiteGroup.Entities()
}

func (c *UserWhiteList) AddWhite(targetId string) error {
	return c.whiteGroup.Accept(targetId)
}

func (c *UserWhiteList) RemoveWhite(targetId string) error {
	return c.whiteGroup.Drop(targetId)
}

func (c *UserWhiteList) OnWhite(targetId string) bool {
	return c.whiteGroup.ContainEntity(targetId)
}

type UserBlackList struct {
	blackGroup basis.IEntityGroup
}

func (c *UserBlackList) Blacks() []string {
	return c.blackGroup.Entities()
}

func (c *UserBlackList) AddBlack(targetId string) error {
	return c.blackGroup.Accept(targetId)
}

func (c *UserBlackList) RemoveBlack(targetId string) error {
	return c.blackGroup.Drop(targetId)
}

func (c *UserBlackList) OnBlack(targetId string) bool {
	return c.blackGroup.ContainEntity(targetId)
}
