//
//Created by xuzhuoxi
//on 2019-03-15.
//@author xuzhuoxi
//
package manager

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"math"
	"sync"
)

type IBroadcastManager interface {
	basis.IManagerBase
	netx.ISockServerSetter
	netx.IAddressProxySetter

	//以下为基础方法------

	//广播整个实体
	//target为环境实体
	//source可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastEntity(source basis.IUserEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error
	//广播部分用户
	//targets为用户实体IUserEntity的UID集合
	//source可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastUsers(source basis.IUserEntity, targets []string, handler func(entity basis.IUserEntity)) error
	//设置附近值
	SetNearDistance(distance float64)
	//广播当前用户所在区域
	//source不能为nil
	BroadcastCurrent(source basis.IUserEntity, excludeBlack bool, handler func(entity basis.IUserEntity)) error
	//以下为业务型方法------

	//环境实体变量更新
	NotifyEnvVars(varTarget basis.IEntity, vars basis.VarSet)
	//用户实体变量更新
	NotifyUserVars(source basis.IUserEntity, vars basis.VarSet)
	//用户实体变量更新
	NotifyUserVarsCurrent(source basis.IUserEntity, vars basis.VarSet)
}

func NewIBroadcastManager(entityMgr IEntityManager, sockServer netx.ISockServer, addressProxy netx.IAddressProxy) IBroadcastManager {
	return NewBroadcastManager(entityMgr, sockServer, addressProxy)
}

func NewBroadcastManager(entityMgr IEntityManager, sockServer netx.ISockServer, addressProxy netx.IAddressProxy) *BroadcastManager {
	return &BroadcastManager{entityMgr: entityMgr, sockServer: sockServer, addressProxy: addressProxy, logger: logx.DefaultLogger(), distance: math.MaxFloat64}
}

//----------------------------------

type BroadcastManager struct {
	entityMgr    IEntityManager
	sockServer   netx.ISockServer
	addressProxy netx.IAddressProxy
	logger       logx.ILogger
	broadcastMu  sync.RWMutex
	distance     float64
}

func (m *BroadcastManager) InitManager() {
	return
}

func (m *BroadcastManager) DisposeManager() {
	return
}

func (m *BroadcastManager) SetLogger(logger logx.ILogger) {
	m.logger = logger
}

func (m *BroadcastManager) SetServer(server netx.ISockServer) {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	m.sockServer = server
}

func (m *BroadcastManager) SetAddressProxy(addressProxy netx.IAddressProxy) {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	m.addressProxy = addressProxy
}

func (m *BroadcastManager) BroadcastEntity(source basis.IUserEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.RLock()
	defer m.broadcastMu.RUnlock()
	if nil == target {
		return errors.New(fmt.Sprintf("Target is nil. "))
	}
	if userTarget, ok := target.(basis.IUserEntity); ok {
		if nil != source && (checkSame(source, userTarget) || checkBlack(source, userTarget)) { //本身 或 黑名单
			return nil
		}
		handler(userTarget)
		return nil
	}
	if entityContainer, ok := target.(basis.IEntityContainer); ok { //容器判断
		if nil == source {
			entityContainer.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
				handler(child.(basis.IUserEntity))
			}, true)
		} else {
			entityContainer.ForEachChild(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
				if basis.EntityUser != child.EntityType() || checkSame(source, child) { //不是用户实体 或 是自己本身
					return
				}
				if checkBlack(source, child) { //黑名单
					return false, true
				}
				handler(child.(basis.IUserEntity))
				return
			})
		}
	}
	return nil
}

func (m *BroadcastManager) BroadcastUsers(source basis.IUserEntity, targets []string, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.RLock()
	defer m.broadcastMu.RUnlock()
	if len(targets) == 0 {
		return errors.New("Targets's len is 0")
	}
	userIndex := m.entityMgr.UserIndex()
	for _, targetId := range targets {
		if targetUser := userIndex.GetUser(targetId); nil != targetUser { //目标用户存在
			if nil != source {
				if checkSame(source, targetUser) { //本身
					continue
				}
				if checkBlack(source, targetUser) { //黑名单
					continue
				}
			}
			handler(targetUser)
		}
	}
	return nil
}

func (m *BroadcastManager) SetNearDistance(distance float64) {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	m.distance = distance
}

func (m *BroadcastManager) BroadcastCurrent(source basis.IUserEntity, excludeBlack bool, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.RLock()
	defer m.broadcastMu.RUnlock()
	if nil == source {
		return errors.New(fmt.Sprintf("Source is nil. "))
	}
	idType, id := source.GetLocation()
	if parentEntity, ok := m.entityMgr.GetEntity(idType, id); ok {
		if ec, ok2 := parentEntity.(basis.IEntityContainer); ok2 {
			ec.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
				if checkSame(source, child) { //本身
					return
				}
				if userChild, ok := child.(basis.IUserEntity); ok {
					if excludeBlack && checkBlack(source, userChild) { //黑名单
						return
					}
					if !basis.NearXYZ(source.GetPosition(), userChild.GetPosition(), m.distance) { //位置不相近
						return
					}
					handler(userChild)
				}
			}, false)
		}
	}
	return nil
}

//-----------------------------

func (m *BroadcastManager) NotifyEnvVars(varTarget basis.IEntity, vars basis.VarSet) {
}

func (m *BroadcastManager) NotifyUserVars(source basis.IUserEntity, vars basis.VarSet) {
}

func (m *BroadcastManager) NotifyUserVarsCurrent(source basis.IUserEntity, vars basis.VarSet) {
}

//-----------------------------

func checkSame(source basis.IUserEntity, target basis.IEntity) bool {
	return source.UID() == target.UID()
}

func checkBlack(source basis.IUserEntity, target basis.IEntity) bool {
	if source.OnBlack(target.UID()) { //source黑名单
		return true
	}
	if userTarget, ok := target.(basis.IUserEntity); ok { //target黑名单
		return userTarget.OnBlack(source.UID())
	}
	return false
}
