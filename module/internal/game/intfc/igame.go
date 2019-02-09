package intfc

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
}
