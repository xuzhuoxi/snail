//
//Created by xuzhuoxi
//on 2019-02-19.
//@author xuzhuoxi
//
package basis

type XY struct {
	X float64
	Y float64
}

type XYZ struct {
	X float64
	Y float64
	Z float64
}

func (xyz XYZ) XY() XY {
	return XY{X: xyz.X, Y: xyz.Y}
}

//判断两点是否相近
//用于转发附近消息
func NearXY(pos1 XY, pos2 XY, distance float64) bool {
	x12 := pos1.X - pos2.X
	y12 := pos1.Y - pos2.Y
	return (x12*x12 + y12*y12) <= distance*distance
}

//判断两点是否相近
//用于转发附近消息
func NearXYZ(pos1 XYZ, pos2 XYZ, distance float64) bool {
	x12 := pos1.X - pos2.X
	y12 := pos1.Y - pos2.Y
	if 0 == pos1.Z && 0 == pos2.Z {
		return NearXY(pos1.XY(), pos2.XY(), distance)
	} else {
		z12 := pos1.Z - pos2.Z
		return (x12*x12 + y12*y12 + z12*z12) <= distance*distance
	}
}
