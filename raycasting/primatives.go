package raycasting

import (
	"math"
)

type Vec1 struct {
	d [1]float64
}

func NewVec1(x float64) Vec1 {
	return Vec1{[1]float64{x}}
}

func (v *Vec1) Len() float64 {
	return v.d[0]
}

func (v *Vec1) Len2() float64 {
	return v.d[0] * v.d[0]
}

type Vec2 struct {
	d [2]float64
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{[2]float64{x, y}}
}

func (v *Vec2) Len() float64 {
	r_2 := v.d[0]*v.d[0] + v.d[1]*v.d[1]
	return math.Sqrt(r_2)
}

func (v *Vec2) Len2() float64 {
	return v.d[0]*v.d[0] + v.d[1]*v.d[1]
}

type Vec3 struct {
	d [3]float64
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{[3]float64{x, y, z}}
}

func (v *Vec3) Len() float64 {
	r_2 := v.d[0]*v.d[0] + v.d[1]*v.d[1] + v.d[2]*v.d[2]
	return math.Sqrt(r_2)
}

func (v *Vec3) Len2() float64 {
	return v.d[0]*v.d[0] + v.d[1]*v.d[1] + v.d[2]*v.d[2]
}

type Vec4 struct {
	d [4]float64
}

func NewVec4(x, y, z, a float64) Vec4 {
	return Vec4{[4]float64{x, y, z, a}}
}

func (v *Vec4) Len() float64 {
	r_2 := v.d[0]*v.d[0] + v.d[1]*v.d[1]
	r_2 += v.d[2]*v.d[2] + v.d[3]*v.d[3]
	return math.Sqrt(r_2)
}

func (v *Vec4) Len2() float64 {
	return v.d[0]*v.d[0] + v.d[1]*v.d[1] + v.d[2]*v.d[2] + v.d[3]*v.d[3]
}
