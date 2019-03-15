//
//Created by xuzhuoxi
//on 2019-03-15.
//@author xuzhuoxi
//
package mmo

import (
	"github.com/pkg/errors"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
	"github.com/xuzhuoxi/snail/engine/mmo/manager"
)

type IMMOManager interface {
	InitMMO() error
	netx.ISockServerSetter
	netx.IAddressProxySetter
	logx.ILoggerSetter

	GetEntityManager() manager.IEntityManager
	GetUserManager() manager.IUserManager
	GetBroadcastManager() manager.IBroadcastManager
}

func NewIMMOManager() IMMOManager {
	return NewMMOManager()
}

func NewMMOManager() *MMOManager {
	return &MMOManager{}
}

//----------------------------

type MMOManager struct {
	entityMgr manager.IEntityManager
	userMgr   manager.IUserManager
	bcMgr     manager.IBroadcastManager
	logger    logx.ILogger
}

func (m *MMOManager) InitMMO() error {
	if nil != m.entityMgr {
		return errors.New("Manager is already init. ")
	}
	m.entityMgr = manager.NewIEntityManager()
	m.entityMgr.AddEventListener(basis.EventSetVariable, m.onEntityVar)
	m.entityMgr.AddEventListener(basis.EventSetMultiVariable, m.onEntityVar)
	m.userMgr = manager.NewIUserManager(m.entityMgr)
	m.bcMgr = manager.NewIBroadcastManager(m.entityMgr, nil, nil)
	return nil
}

func (m *MMOManager) SetServer(server netx.ISockServer) {
	if nil != m.bcMgr {
		m.bcMgr.SetServer(server)
	}
}

func (m *MMOManager) SetAddressProxy(proxy netx.IAddressProxy) {
	if nil != m.bcMgr {
		m.bcMgr.SetAddressProxy(proxy)
	}
}

func (m *MMOManager) SetLogger(logger logx.ILogger) {
	m.logger = logger
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

//------------------------------

func (m *MMOManager) onEntityVar(evd *eventx.EventData) {
	varSet := evd.Data.(basis.VarSet)
	target := varSet["Target"]
	m.logger.Traceln("onEntityVar", target)
}
