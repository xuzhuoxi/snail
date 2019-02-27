//
//Created by xuzhuoxi
//on 2019-02-27.
//@author xuzhuoxi
//
package extension

type ISnailInitExtension interface {
	InitExtension() error
}

type ISnailSaveExtension interface {
	SaveExtension() error
}

type ISnailDestroyExtension interface {
	DestroyExtension() error
}
