package ifc

import "github.com/xuzhuoxi/snail/module/imodule"

type ILoginServer interface {
	//登录
	Login()
	//登出
	Logout()
}

type IGameSockStatus interface {
	//服务器运行时间(秒)
	GetPassSecond() int64
	//服务器详细状态
	GetStateDetail() imodule.IServiceStateDetail
	//服务器简单状态
	GetStateSimple() imodule.ServiceState
}
