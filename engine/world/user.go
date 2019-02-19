//
//Created by xuzhuoxi
//on 2019-02-18.
//@author xuzhuoxi
//
package world

import (
	"github.com/pkg/errors"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

type IUserChannel interface {
	//订阅的通信频道
	TouchingChannels() []string

	//订阅频道
	TouchChannel(chanId string) error
	//取消频道订阅
	UnTouchChannel(chanId string) error
	//处于频道
	InChannel(chanId string) bool
}

type IUser interface {
	IEntity
	IUserChannel
	UserName() string
	NickName() string

	SetUserVariables(vars Variables)
}

type User struct {
	Uid  string
	Name string
	Nick string

	Addr     string
	ZoneId   string
	RoomId   string
	Channels []string

	Pos XYZ

	chanMu sync.RWMutex
}

func (u *User) UID() string {
	return u.Uid
}

func (u *User) UserName() string {
	return u.Name
}

func (u *User) NickName() string {
	return u.Nick
}

func (u *User) TouchingChannels() []string {
	return u.Channels
}

func (u *User) TouchChannel(chanId string) error {
	if u.InChannel(chanId) {
		return errors.New("TouchChannel Error :" + chanId)
	}
	u.chanMu.Lock()
	defer u.chanMu.Unlock()
	u.Channels = append(u.Channels, chanId)
	return nil
}

func (u *User) UnTouchChannel(chanId string) error {
	index, ok := slicex.IndexString(u.Channels, chanId)
	if !ok {
		return errors.New("UnTouchChannel Error :" + chanId)
	}
	u.chanMu.Lock()
	u.Channels = append(u.Channels[:index], u.Channels[index+1:]...)
	return nil
}

func (u *User) InChannel(chanId string) bool {
	_, ok := slicex.IndexString(u.Channels, chanId)
	return ok
}
