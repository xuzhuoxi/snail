package impl

import (
	"encoding/binary"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/snail/module/imodule"
)

var GobOrder = binary.BigEndian

type ModuleGame struct {
	imodule.ModuleBase
	rpcRemoteMap map[string]netx.IRPCClient
	state        *imodule.ServiceStateDetail

	gobBuffEncoder encodingx.IGobBuffEncoder
	gobBuffDecoder encodingx.IGobBuffDecoder
}

func (m *ModuleGame) Init() {
	m.gobBuffEncoder = encodingx.NewDefaultGobBuffEncoder()
	m.gobBuffDecoder = encodingx.NewDefaultGobBuffDecoder()
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
