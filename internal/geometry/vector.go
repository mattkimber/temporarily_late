package geometry

import "math"

type Vector3 struct {
	X, Y, Z float64
}

type Vector2 struct {
	X, Y float64
}

func DegToRad(angle float64) float64 {
	return (angle / 180.0) * math.Pi
}

func (a Vector2) Add(b Vector2) Vector2 {
	return Vector2{a.X + b.X, a.Y + b.Y}
}

func (a Vector2) DistanceSquared(b Vector2) float64 {
	return ((a.X - b.X) * (a.X - b.X)) + ((a.Y - b.Y) * (a.Y - b.Y))
}

func (a Vector2) LengthSquared() float64 {
	return a.Dot(a)
}

func (a Vector2) Dot(b Vector2) float64 {
	return (a.X * b.X) + (a.Y * b.Y)
}

func (a Vector2) DivideByVector(by Vector2) Vector2 {
	return Vector2{a.X / by.X, a.Y / by.Y}
}

func Zero() Vector3 {
	return Vector3{0, 0, 0}
}

func UnitX() Vector3 {
	return Vector3{1, 0, 0}
}

func UnitY() Vector3 {
	return Vector3{0, 1, 0}
}

func UnitZ() Vector3 {
	return Vector3{0, 0, 1}
}

func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector3) Subtract(b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector3) MultiplyByConstant(by float64) Vector3 {
	return Vector3{a.X * by, a.Y * by, a.Z * by}
}

func (a Vector3) MultiplyByVector(by Vector3) Vector3 {
	return Vector3{a.X * by.X, a.Y * by.Y, a.Z * by.Z}
}

func (a Vector3) DivideByConstant(by float64) Vector3 {
	return Vector3{a.X / by, a.Y / by, a.Z / by}
}

func (a Vector3) DivideByVector(by Vector3) Vector3 {
	return Vector3{a.X / by.X, a.Y / by.Y, a.Z / by.Z}
}

func (a Vector3) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a Vector3) Normalise() Vector3 {
	return a.DivideByConstant(a.Length())
}

func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{
		(a.Y * b.Z) - (a.Z * b.Y),
		(a.Z * b.X) - (a.X * b.Z),
		(a.X * b.Y) - (a.Y * b.X),
	}
}

func (a Vector3) Dot(b Vector3) float64 {
	return (a.X * b.X) + (a.Y * b.Y) + (a.Z * b.Z)
}

func (a Vector3) Lerp(b Vector3, amt float64) Vector3 {
	return a.MultiplyByConstant(1 - amt).Add(b.MultiplyByConstant(amt))
}

func (a Vector3) Equals(b Vector3) bool {
	const epsilon = 1e-12
	return a.Subtract(b).Length() < epsilon
}
