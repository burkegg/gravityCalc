package pkg
//
//import "math"
//
//type Vector struct {
//	X float64
//	Y float64
//}
//
//// Calculate the unit vector
//func Normalize(v Vector) (u Vector) {
//	// if length is zero, set x = 1, y = 0
//	vLen := Length(v)
//	if vLen == 0 {
//		u = Vector{
//			1,
//			0,
//		}
//		return u
//	}
//	x := v.X / vLen
//	y := v.Y / vLen
//	u = Vector{
//		x,
//		y,
//	}
//	u.X = v.X / vLen
//	u.Y = v.Y / vLen
//	return u
//}
//
//func Subtract(v1 Vector, v2 Vector) (res Vector) {
//	x := v2.X - v1.X
//	y := v2.Y - v1.Y
//	res = Vector{
//		X: x,
//		Y: y,
//	}
//	return res
//}
//
//func Length(v Vector) (s float64) {
//	s = math.Sqrt(LengthSq(v))
//	return s
//}
//
//func LengthSq(v Vector) (s2 float64) {
//	s2 = v.X*v.X + v.Y*v.Y
//	return s2
//}
//
//func DistSq(v1 Vector, v2 Vector) (s2 float64) {
//	s2 = math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2)
//	return s2
//}
//
//func VectorDist(v1 Vector, v2 Vector) (s float64) {
//	s = math.Sqrt(DistSq(v1, v2))
//	return s
//}
//
//// Multiply a vector by a scalar and return result
//func ScalarMult(v Vector, s float64) (res Vector) {
//	res.X = v.X * s
//	res.Y = v.Y * s
//	return res
//}
//
//func ScalarDivide(v Vector, s float64) (res Vector) {
//	res.X = v.X / s
//	res.Y = v.Y / s
//	return res
//}
//
//func Add(v1 Vector, v2 Vector) (f Vector) {
//	f.X = v1.X + v2.X
//	f.Y = v1.Y + v2.Y
//	return f
//}
