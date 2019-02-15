package imodule

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

const (
	ServiceMethod_OnRPCCall = "RPCHandler.OnRPCCall"

	CmdRoute_OnConnected    = "ModRoute.OnConnected"
	CmdRoute_OnDisconnected = "ModRoute.OnDisconnected"
	CmdRoute_UpdateState    = "ModRoute.UpdateState"
)

func MapRPCHandler(h *RPCHandler, key string, f func(args *RPCArgs, reply *RPCReply) error) {
	if nil == h.handler {
		h.handler = make(map[string]func(args *RPCArgs, reply *RPCReply) error)
	}
	h.handler[key] = f
}

type RPCHandler struct {
	Log     logx.ILogger
	handler map[string]func(args *RPCArgs, reply *RPCReply) error
}

type RPCArgs struct {
	From string
	Cmd  string
	Data []byte
}

type RPCReply struct {
	To   string
	Cmd  string
	Data []byte
}

func (g *RPCHandler) OnRPCCall(args *RPCArgs, reply *RPCReply) error {
	//g.Log.Infoln("\tOnRPCCall:", args, reply)
	handler, ok := g.handler[args.Cmd]
	if ok {
		return handler(args, reply)
	} else {
		g.Log.Fatalln("\tRPCHandler Map Error: ", args.Cmd)
		return nil
	}
	return nil
}
