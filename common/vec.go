package common

import (
	"errors"
	"math"
)

type Number interface {
	int32 | float32 | int8
}

type Vec[t Number] struct {
	X t
	Y t
}

func NewVec[t Number](x, y t) Vec[t] {
	return Vec[t]{
		X: x,
		Y: y,
	}
}

func (v Vec[t]) Get() (t, t) {
	return v.X, v.Y
}

func (v *Vec[t]) Clone() *Vec[t] {
	return &Vec[t]{
		v.X,
		v.Y,
	}

}

func (v *Vec[t]) SetXY(a Vec[t]) {
	v.X = a.X
	v.Y = a.Y
}

func (v *Vec[t]) Abs() *Vec[t] {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v
}

func (v *Vec[t]) CapAt(a Vec[t]) *Vec[t] {
	if v.X > 0 && v.X > a.X {
		v.X = a.X
	} else if v.X < 0 && v.X < -a.X {
		v.X = -a.X
	}
	if v.Y > 0 && v.Y > a.Y {
		v.Y = a.Y
	} else if v.Y < 0 && v.Y < -a.Y {
		v.Y = -a.Y
	}
	return v
}

func (v Vec[t]) Magnitude() t {
	return t(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vec[t]) rotate(origin Vec[t], angle float32) {
	angle = (angle / 180) * math.Pi

	cos := t(math.Cos(float64(angle)))
	sin := t(math.Sin(float64(angle)))

	pos := v.Clone().Sub(origin)

	// Rotate
	x := pos.X*cos - pos.Y*sin
	y := pos.X*sin + pos.Y*cos

	// Translate back
	v.X = x + origin.X
	v.Y = y + origin.Y
}

func (v Vec[t]) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X))) * (180 / math.Pi)
}

func (v Vec[t]) Normalize() *Vec[t] {
	mag := v.Magnitude()
	if mag == 0 {
		return &Vec[t]{}
	}
	return &Vec[t]{
		X: v.X / mag,
		Y: v.Y / mag,
	}
}

func (v *Vec[t]) Negate() *Vec[t] {
	return &Vec[t]{X: -v.X, Y: -v.Y}
}

func (v *Vec[t]) MultScalar(a t) *Vec[t] {
	v.X *= a
	v.Y *= a
	return v
}

func DotProduct[t Number](a, b Vec[t]) t {
	return a.X*b.X + a.Y*b.Y
}

func (v *Vec[t]) Add(a Vec[t]) *Vec[t] {
	v.X += a.X
	v.Y += a.Y
	return v
}

func (v *Vec[t]) Sub(a Vec[t]) *Vec[t] {
	v.X -= a.X
	v.Y -= a.Y
	return v
}

func (v *Vec[t]) AddScalars(x, y t) *Vec[t] {
	v.X += x
	v.Y += y
	return v
}

func (v *Vec[t]) SubScalars(x, y t) *Vec[t] {
	v.X -= x
	v.Y -= y
	return v
}

func FromAtoBVec[t Number](a, b Vec[t]) Vec[t] {
	return Vec[t]{X: b.X - a.X, Y: b.Y - a.Y}
}

func DistanceAtoBVecByEuclidean[t Number](a, b Vec[t]) t {
	return a.Clone().Sub(b).Magnitude()
}

func DistanceAtoBVecShev[t Number](a, b Vec[t]) t {
	diff := a.Clone().Sub(b)
	diff = diff.Abs()
	if diff.X < diff.Y {
		return diff.Y
	}
	return diff.X
}

func (v Vec[t]) Normal() Vec[t] {
	return Vec[t]{X: -v.Y, Y: v.X}
}

func (v *Vec[t]) IsNull() bool {
	return v.X == 0 && v.Y == 0
}

func (v *Vec[t]) IsEqual(a Vec[t]) bool {
	return v.X == a.X && v.Y == a.Y
}

func SubVecs[t Number](a, b Vec[t]) Vec[t] {
	return Vec[t]{
		a.X - b.X,
		a.Y - b.Y,
	}
}

func AddVecs[t Number](v ...Vec[t]) (res Vec[t]) {
	for i := range v {
		res.Add(v[i])
	}
	return res
}

func CastVec[from, to Number](a Vec[from]) Vec[to] {
	return Vec[to]{
		X: to(a.X),
		Y: to(a.Y),
	}
}

func CmpVecSort[t Number](a, b Vec[t]) int {
	return int(DistanceAtoBVecByEuclidean(a, b))
}

func GetNeighborhoodPos(pos Vec[int32], size Vec[int32]) (res []Vec[int32], err error) {

	if size.X <= 0 || size.Y <= 0 {
		return nil, errors.New("invalid neighborhood size")
	}
	if size.X%2 == 0 || size.Y%2 == 0 {
		return nil, errors.New("neighborhood size must be odd")
	}
	halfX := size.X / 2
	halfY := size.Y / 2

	res = make([]Vec[int32], size.X*size.Y)
	var worldX, worldY int32
	for ix := range size.X {
		for iy := range size.Y {
			worldX = pos.X + (ix - halfX)
			worldY = pos.Y + (iy - halfY)

			res = append(res, NewVec(worldX, worldY))
		}
	}
	return res, nil
}
