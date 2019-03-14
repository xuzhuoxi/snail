//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

type IEntityChild interface {
	GetParent() string
	NoneParent() bool

	SetParent(ownerId string)
	ClearParent()
}
