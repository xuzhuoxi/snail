//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package world

type XYZ struct {
	X float64
	Y float64
	Z float64
}

//判断两点是否相近
//用于转发附近消息
func Near(pos1 XYZ, pos2 XYZ, distance float64) bool {
	x12 := pos1.X - pos2.X
	y12 := pos1.Y - pos2.Y
	if 0 == pos1.Z && 0 == pos2.Z {
		return (x12*x12 + y12*y12) <= distance*distance
	} else {
		z12 := pos1.Z - pos2.Z
		return (x12*x12 + y12*y12 + z12*z12) <= distance*distance
	}
}
