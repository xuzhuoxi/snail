package impl

import (
	"github.com/xuzhuoxi/snail/module/imodule"
	"github.com/xuzhuoxi/util-go/encodingx"
	"github.com/xuzhuoxi/util-go/netx"
)

type ModuleGame struct {
	imodule.ModuleBase
	rpcRemoteMap map[string]netx.IRPCClient
	state        *imodule.ServiceStateDetail

	gobBuffEncoder encodingx.IGobBuffEncoder
	gobBuffDecoder encodingx.IGobBuffDecoder
}

func (m *ModuleGame) Init() {
	m.gobBuffEncoder = encodingx.NewGobBuffEncoder()
	m.gobBuffDecoder = encodingx.NewGobBuffDecoder()
	m.rpcRemoteMap = make(map[string]netx.IRPCClient)
	m.state = imodule.NewServiceState(m.GetId(), imodule.DefaultStatsInterval)
}

func (m *ModuleGame) Run() {
	m.state.Start()
	go CheckRPC(m)
}

func (m *ModuleGame) Save() {
	panic("implement me")
}

func (m *ModuleGame) OnDestroy() {
}

func (m *ModuleGame) Destroy() {
}
