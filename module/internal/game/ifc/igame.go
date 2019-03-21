package ifc

import "github.com/xuzhuoxi/snail/module/imodule"

type ILoginServer interface {
	//登录
	Login()
	//登出
	Logout()
}

type IGameStatus interface {
	//服务器运行时间(秒)
	GetPassTime() int64
	//服务器状态系数
	GetStatePriority() float64
	//详细状态
	DetailState() *imodule.ServiceStateDetail
	//转化为简单的状态
	ToSimpleState() imodule.ServiceState
}
