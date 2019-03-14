//
//Created by xuzhuoxi
//on 2019-03-13.
//@author xuzhuoxi
//
package proxy

import "github.com/xuzhuoxi/infra-go/netx"

type IServerProxy interface {
	netx.ISockSender

	SetServer(server netx.ISockServer)
}
