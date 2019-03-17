//
//Created by xuzhuoxi
//on 2019-03-17.
//@author xuzhuoxi
//
package index

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	a1 := NewIChannelIndex()
	a2 := NewIRoomIndex()
	a3 := NewITeamCorpsIndex()
	a4 := NewITeamIndex()
	a5 := NewIUserIndex()
	a6 := NewIZoneIndex()
	fmt.Println(a1.AddChannel(nil))
	fmt.Println(a2.AddRoom(nil))
	fmt.Println(a3.AddCorps(nil))
	fmt.Println(a4.AddTeam(nil))
	fmt.Println(a5.AddUser(nil))
	fmt.Println(a6.AddZone(nil))
}
