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

func MapRPCFunction(handler *RPCHandler, key string, f func(args *RPCArgs, reply *RPCReply) error) {
	handler.handlerMap[key] = f
}

func NewRPCHandler(logger logx.ILogger) *RPCHandler {
	handler := new(RPCHandler)
	handler.Logger = logger
	handler.handlerMap = make(map[string]func(args *RPCArgs, reply *RPCReply) error)
	return handler
}

//RPCHandler要求
//全部方法必须是func(args *RPCArgs, reply *RPCReply) error
//不然会报警告
type RPCHandler struct {
	Logger     logx.ILogger
	handlerMap map[string]func(args *RPCArgs, reply *RPCReply) error
}

func (g *RPCHandler) OnRPCCall(args *RPCArgs, reply *RPCReply) error {
	//g.Log.Infoln("\tOnRPCCall:", args, reply)
	handler, ok := g.handlerMap[args.Cmd]
	if ok {
		return handler(args, reply)
	} else {
		if nil != g.Logger {
			g.Logger.Fatalln("\tRPCHandler Map Error: ", args.Cmd)
		}
		return nil
	}
	return nil
}
