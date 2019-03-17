//
//Created by xuzhuoxi
//on 2019-03-15.
//@author xuzhuoxi
//
package manager

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"sync"
)

type IUserManager interface {
	basis.IManagerBase
	//加入世界
	EnterWorld(user basis.IUserEntity, roomId string) error
	//离开世界
	ExitWorld(userId string) error
	//在世界转移
	Transfer(userId string, toRoomId string) error
}

func NewIUserManager(entityMgr IEntityManager) IUserManager {
	return NewUserManager(entityMgr)
}

func NewUserManager(entityMgr IEntityManager) *UserManager {
	return &UserManager{entityMgr: entityMgr, logger: logx.DefaultLogger()}
}

//----------------------------

type UserManager struct {
	entityMgr  IEntityManager
	logger     logx.ILogger
	transferMu sync.RWMutex
}

func (w *UserManager) InitManager() {
	return
}

func (w *UserManager) DisposeManager() {
	return
}

func (m *UserManager) SetLogger(logger logx.ILogger) {
	m.logger = logger
}

func (w *UserManager) EnterWorld(user basis.IUserEntity, roomId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	if nil == user {
		return errors.New("WorldManager.EnterWorld Error: user is nil")
	}
	roomIndex := w.entityMgr.RoomIndex()
	userIndex := w.entityMgr.UserIndex()
	if !roomIndex.CheckRoom(roomId) {
		return errors.New("WorldManager.EnterWorld Error: Room(" + roomId + ") does not exist")
	}
	userId := user.UID()
	if userIndex.CheckUser(userId) {
		oldUser := userIndex.GetUser(userId)
		w.exitCurrentRoom(oldUser)
	}
	userIndex.UpdateUser(user)
	room := roomIndex.GetRoom(roomId)
	room.AddChild(user)
	user.SetLocation(basis.EntityRoom, roomId)
	return nil
}

func (w *UserManager) ExitWorld(userId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	userIndex := w.entityMgr.UserIndex()
	if "" == userId || userIndex.CheckUser(userId) {
		return errors.New("WorldManager.ExitWorld Error: User() does not exist")
	}
	roomIndex := w.entityMgr.RoomIndex()
	user := userIndex.GetUser(userId)
	_, roomId := user.GetLocation()
	if room := roomIndex.GetRoom(roomId); nil != room {
		user.DestroyEntity()
		room.RemoveChild(user)
	}
	return nil
}

func (w *UserManager) Transfer(userId string, toRoomId string) error {
	w.transferMu.Lock()
	defer w.transferMu.Unlock()
	userIndex := w.entityMgr.UserIndex()
	if "" == userId || !userIndex.CheckUser(userId) {
		return errors.New(fmt.Sprintf("EnterWorld Error: user(%s) does not exist", userId))
	}
	roomIndex := w.entityMgr.RoomIndex()
	if "" == toRoomId || !roomIndex.CheckRoom(toRoomId) {
		return errors.New(fmt.Sprintf("EnterWorld Error: Target room(%s) does not exist", toRoomId))
	}
	user := userIndex.GetUser(userId)
	_, roomId := user.GetLocation()
	if roomId == toRoomId {
		return errors.New(fmt.Sprintf("EnterWorld Error: user(%s) already in the room(%s)", userId, toRoomId))
	}
	w.exitCurrentRoom(user)
	toRoom := roomIndex.GetRoom(toRoomId)
	toRoom.AddChild(user)
	user.SetLocation(basis.EntityRoom, roomId)
	return nil
}

func (w *UserManager) exitCurrentRoom(user basis.IUserEntity) error {
	_, roomId := user.GetLocation()
	roomIndex := w.entityMgr.RoomIndex()
	if "" == roomId || !roomIndex.CheckRoom(roomId) {
		return errors.New("WorldManager.exitCurrentRoom Error: room is nil")
	}
	room := roomIndex.GetRoom(roomId)
	err := room.RemoveChild(user)
	if nil != err {
		return err
	}
	user.SetLocation(basis.EntityNone, "")
	return nil
}
