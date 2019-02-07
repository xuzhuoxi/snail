package impl

import (
	"github.com/xuzhuoxi/snail/module/intfc"
	"github.com/xuzhuoxi/util-go/encodingx"
	"github.com/xuzhuoxi/util-go/netx"
	"time"
)

type ModuleGame struct {
	intfc.ModuleBase
	remoteMap map[string]netx.IRPCClient
	state     intfc.GameServerState
	starting  int64

	codecs *encodingx.GobCodecs
}

func (m *ModuleGame) Init() {
	m.codecs = encodingx.NewCodecs()
	m.remoteMap = make(map[string]netx.IRPCClient)
}

func (m *ModuleGame) Run() {
	m.starting = time.Now().UnixNano()
	go CheckRPC(m)
}

func (m *ModuleGame) Save() {
	panic("implement me")
}

func (m *ModuleGame) OnDestroy() {
}

func (m *ModuleGame) Destroy() {
}
