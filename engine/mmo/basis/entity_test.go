//
//Created by xuzhuoxi
//on 2019-03-17.
//@author xuzhuoxi
//
package basis

import (
	"fmt"
	"testing"
)

func TestEntityType(t *testing.T) {
	fmt.Println(EntityNone, EntityWorld, EntityZone, EntityRoom, EntityUser, EntityTeamCorps, EntityTeam, EntityChannel)
	fmt.Println("---")
	fmt.Println(EntityAll.Match(EntityWorld))
	fmt.Println(EntityWorld.Match(EntityAll))
	fmt.Println(EntityWorld.Match(EntityZone))
	fmt.Println("---")
	fmt.Println(EntityAll.Include(EntityWorld))
	fmt.Println(EntityWorld.Include(EntityAll))
	fmt.Println(EntityWorld.Include(EntityZone))
	fmt.Println(EntityWorld.Include(EntityNone))
}
