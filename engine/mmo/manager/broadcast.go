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
	"sync"
)

type IBroadcastManager interface {
	basis.IManagerBase
	netx.ISockServerSetter
	netx.IAddressProxySetter

	//环境实体变量更新
	NotifyEnvVars(varTarget basis.IEntity, vars basis.VarSet)
	//用户实体变量更新
	NotifyUserVars(source basis.IUserEntity, vars basis.VarSet)
	//用户实体变量更新
	NotifyUserVarsCurrent(source basis.IUserEntity, vars basis.VarSet)

	//广播整个实体
	Broadcast(source basis.IEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error
	//广播整个实体，过滤掉黑名单部分
	BroadcastWithoutBlack(source basis.IEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error
	//广播部分用户
	BroadcastUsers(source basis.IEntity, targets []string, handler func(entity basis.IUserEntity)) error
	//广播部分用户
	BroadcastUsersWithoutBlack(source basis.IEntity, targets []string, handler func(entity basis.IUserEntity)) error
	//广播当前用户
	BroadcastCurrent(source basis.IEntity, handler func(entity basis.IUserEntity)) error
	//广播当前用户
	BroadcastCurrentWithoutBlack(source basis.IEntity, handler func(entity basis.IUserEntity)) error
}

func NewIBroadcastManager(entityMgr IEntityManager, sockServer netx.ISockServer, addressProxy netx.IAddressProxy) IBroadcastManager {
	return NewBroadcastManager(entityMgr, sockServer, addressProxy)
}

func NewBroadcastManager(entityMgr IEntityManager, sockServer netx.ISockServer, addressProxy netx.IAddressProxy) *BroadcastManager {
	return &BroadcastManager{entityMgr: entityMgr, sockServer: sockServer, addressProxy: addressProxy}
}

//----------------------------------

type BroadcastManager struct {
	entityMgr    IEntityManager
	sockServer   netx.ISockServer
	addressProxy netx.IAddressProxy
	logger       logx.ILogger
	broadcastMu  sync.RWMutex
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

func (m *BroadcastManager) NotifyEnvVars(varTarget basis.IEntity, vars basis.VarSet) {
}

func (m *BroadcastManager) NotifyUserVars(source basis.IUserEntity, vars basis.VarSet) {
}

func (m *BroadcastManager) NotifyUserVarsCurrent(source basis.IUserEntity, vars basis.VarSet) {
}

func (m *BroadcastManager) Broadcast(source basis.IEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	Broadcast(source, target, handler)
	return nil
}

func (m *BroadcastManager) BroadcastWithoutBlack(source basis.IEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	BroadcastWithoutBlack(source, target, handler)
	return nil
}

func (m *BroadcastManager) BroadcastUsers(source basis.IEntity, targets []string, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	if err := m.checkUserBroadcast(source, targets); nil != err {
		return err
	}
	userIndex := m.entityMgr.UserIndex()
	for _, targetId := range targets {
		if targetEntity := userIndex.GetUser(targetId); targetEntity != nil {
			handler(targetEntity)
		}
	}
	return nil
}

func (m *BroadcastManager) BroadcastUsersWithoutBlack(source basis.IEntity, targets []string, handler func(entity basis.IUserEntity)) error {
	m.broadcastMu.Lock()
	defer m.broadcastMu.Unlock()
	if err := m.checkUserBroadcast(source, targets); nil != err {
		return err
	}
	userIndex := m.entityMgr.UserIndex()
	userSource, userSourceOk := source.(basis.IUserEntity)
	for _, targetId := range targets {
		if targetEntity := userIndex.GetUser(targetId); targetEntity != nil {
			if userSourceOk && userSource.OnBlack(targetEntity.UID()) { //source黑名单
				continue
			}
			if targetEntity.OnBlack(source.UID()) { //target黑名单
				continue
			}
			handler(targetEntity)
		}
	}
	return nil
}

func (m *BroadcastManager) checkUserBroadcast(source basis.IEntity, targets []string) error {
	if nil == source {
		return errors.New(fmt.Sprintf("Source is nil. "))
	}
	if len(targets) == 0 {
		return errors.New(fmt.Sprintf("Target is empty. "))
	}
	return nil
}

func (m *BroadcastManager) BroadcastCurrent(source basis.IEntity, handler func(entity basis.IUserEntity)) error {
	if nil == source {
		return errors.New(fmt.Sprintf("Source is nil. "))
	}
	panic("implement me")
}

func (m *BroadcastManager) BroadcastCurrentWithoutBlack(source basis.IEntity, handler func(entity basis.IUserEntity)) error {
	panic("implement me")
}

//-----------------------------

func Broadcast(source basis.IEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error {
	if err := checkSourceTarget(source, target); nil != err {
		return err
	}
	sourceId := source.UID()
	if userTarget, ok := target.(basis.IUserEntity); ok {
		handler(userTarget)
		return nil
	} else if entityContainer, ok := target.(basis.IEntityContainer); ok {
		var userEntity basis.IUserEntity
		entityContainer.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
			if sourceId != child.UID() { //过滤自己
				userEntity = child.(basis.IUserEntity)
				handler(userEntity)
			}
		}, true)
	}
	return nil
}

func BroadcastWithoutBlack(source basis.IEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error {
	if err := checkSourceTarget(source, target); nil != err {
		return err
	}
	sourceId := source.UID()
	userSource, sourceOk := source.(basis.IUserEntity)
	userTarget, targetOk := target.(basis.IUserEntity)
	if targetOk {
		if userTarget.OnBlack(sourceId) {
			return nil
		}
		if sourceOk && userSource.OnBlack(target.UID()) {
			return nil
		}
		handler(userTarget)
	} else if entityContainer, ok := target.(basis.IEntityContainer); ok {
		var userEntity basis.IUserEntity
		entityContainer.ForEachChild(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
			if sourceId == child.UID() { //过滤自己
				return
			}
			if sourceOk && userSource.OnBlack(child.UID()) { //source黑名单
				return false, true
			}
			if basis.EntityUser == child.EntityType() && sourceId != child.UID() { //是用户实体
				userEntity = child.(basis.IUserEntity)
				if userEntity.OnBlack(sourceId) { //target黑名单
					return
				}
				handler(userEntity)
			}
			return
		})
	}
	return nil
}

func checkSourceTarget(source basis.IEntity, target basis.IEntity) error {
	if nil == source || nil == target {
		return errors.New("Source or target is nil. ")
	}
	if source == target || source.UID() == target.UID() {
		return errors.New("Source is the same as the target. ")
	}
	return nil
}
