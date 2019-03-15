//
//Created by xuzhuoxi
//on 2019-03-14.
//@author xuzhuoxi
//
package basis

type IEntityContainer interface {
	NumChildren() int
	Full() bool

	Contains(entity IEntity) (isContains bool)
	ContainsById(entityId string) (isContains bool)
	GetChildById(entityId string) (entity IEntity, ok bool)
	ReplaceChildInto(entity IEntity) error
	AddChild(entity IEntity) error
	RemoveChild(entity IEntity) error
	RemoveChildById(entityId string) (entity IEntity, ok bool)

	ForEachChild(each func(child IEntity) (interruptCurrent bool, interruptRecurse bool))
	ForEachChildByType(entityType EntityType, each func(child IEntity), recurse bool)
}
