//
//Created by xuzhuoxi
//on 2019-03-16.
//@author xuzhuoxi
//
package manager

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/snail/engine/mmo/basis"
)

type IVariableManager interface {
	basis.IManagerBase
}

func NewIVariableManager(entityManager IEntityManager, broadcastManager IBroadcastManager) IVariableManager {
	return NewVariableManager(entityManager, broadcastManager)
}

func NewVariableManager(entityManager IEntityManager, broadcastManager IBroadcastManager) *VariableManager {
	return &VariableManager{entityMgr: entityManager, bcMgr: broadcastManager, logger: logx.DefaultLogger()}
}

//--------------------------------

type VariableManager struct {
	entityMgr IEntityManager
	bcMgr     IBroadcastManager
	logger    logx.ILogger
}

func (m *VariableManager) InitManager() {
	m.entityMgr.AddEventListener(basis.EventVariableChanged, m.onEntityVar)
}

func (m *VariableManager) DisposeManager() {
	m.entityMgr.RemoveEventListener(basis.EventVariableChanged, m.onEntityVar)
}

func (m *VariableManager) SetLogger(logger logx.ILogger) {
	m.logger = logger
}

func (m *VariableManager) onEntityVar(evd *eventx.EventData) {
	data := evd.Data.([]interface{})
	currentTarget := data[0].(basis.IEntity)
	varSet := data[1].(basis.VarSet)
	if nil != m.logger {
		m.logger.Traceln("onEntityVar", currentTarget.UID(), varSet)
	}
	if currentTarget.EntityType() == basis.EntityUser {
		m.bcMgr.NotifyUserVars(currentTarget.(basis.IUserEntity), varSet)
	} else {
		m.bcMgr.NotifyEnvVars(currentTarget, varSet)
	}
}
