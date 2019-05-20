//
//Created by xuzhuoxi
//on 2019-03-15.
//@author xuzhuoxi
//
package mmo

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"github.com/xuzhuoxi/snail/engine/mmo/manager"
)

type IMMOManager interface {
	basis.IManagerBase
	netx.ISockServerSetter
	netx.IAddressProxySetter

	GetEntityManager() manager.IEntityManager
	GetUserManager() manager.IUserManager
	GetBroadcastManager() manager.IBroadcastManager
}

func NewIMMOManager() IMMOManager {
	return NewMMOManager()
}

func NewMMOManager() *MMOManager {
	return &MMOManager{logger: logx.DefaultLogger()}
}

//----------------------------

type MMOManager struct {
	entityMgr manager.IEntityManager
	userMgr   manager.IUserManager
	bcMgr     manager.IBroadcastManager
	varMgr    manager.IVariableManager
	logger    logx.ILogger
}

func (m *MMOManager) InitManager() {
	if nil != m.entityMgr {
		return
	}
	m.entityMgr = manager.NewIEntityManager()
	m.entityMgr.InitManager()
	m.userMgr = manager.NewIUserManager(m.entityMgr)
	m.userMgr.InitManager()
	m.bcMgr = manager.NewIBroadcastManager(m.entityMgr, nil, nil)
	m.bcMgr.InitManager()
	m.varMgr = manager.NewIVariableManager(m.entityMgr, m.bcMgr)
	m.varMgr.InitManager()
	m.SetLogger(m.logger)
}

func (m *MMOManager) DisposeManager() {
	m.varMgr.DisposeManager()
	m.bcMgr.DisposeManager()
	m.userMgr.DisposeManager()
	m.entityMgr.DisposeManager()
}

func (m *MMOManager) SetSockServer(server netx.ISockServer) {
	if nil != m.bcMgr {
		m.bcMgr.SetSockServer(server)
	}
}

func (m *MMOManager) SetAddressProxy(proxy netx.IAddressProxy) {
	if nil != m.bcMgr {
		m.bcMgr.SetAddressProxy(proxy)
	}
}

func (m *MMOManager) SetLogger(logger logx.ILogger) {
	m.logger = logger
	if nil != m.entityMgr {
		m.entityMgr.SetLogger(logger)
	}
	if nil != m.userMgr {
		m.userMgr.SetLogger(logger)
	}
	if nil != m.bcMgr {
		m.bcMgr.SetLogger(logger)
	}
	if nil != m.varMgr {
		m.varMgr.SetLogger(logger)
	}
}

func (m *MMOManager) GetEntityManager() manager.IEntityManager {
	return m.entityMgr
}

func (m *MMOManager) GetUserManager() manager.IUserManager {
	return m.userMgr
}

func (m *MMOManager) GetBroadcastManager() manager.IBroadcastManager {
	return m.bcMgr
}
