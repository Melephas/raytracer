package main

import (
	"math"
	"strconv"
)

// Vec3 represents a 3D vector or color.
type Vec3 struct {
	I, J, K float64
}

//type Color Vec3
//type Point Vec3

// --- Getter functions ---

// X returns the first component of the vector.
func (v Vec3) X() float64 {
	return v.I
}

// Y returns the second component of the vector.
func (v Vec3) Y() float64 {
	return v.J
}

// Z returns the third component of the vector.
func (v Vec3) Z() float64 {
	return v.K
}

// R returns the first component of the vector, alias for X.
func (v Vec3) R() float64 {
	return v.I
}

// G returns the second component of the vector, alias for Y.
func (v Vec3) G() float64 {
	return v.J
}

// B returns the third component of the vector, alias for Z.
func (v Vec3) B() float64 {
	return v.K
}

// --- Arithmetic functions ---

// Negate returns a flipped version of the vector.
func (v Vec3) Negate() Vec3 {
	return Vec3{
		I: -v.I,
		J: -v.J,
		K: -v.K,
	}
}

// Add returns a new vector that is the sum of the two vectors.
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		I: v.I + other.I,
		J: v.J + other.J,
		K: v.K + other.K,
	}
}

// Sub returns a vector that is the sum of the two vectors, with the second vector flipped (v - other).
func (v Vec3) Sub(other Vec3) Vec3 {
	return v.Add(other.Negate())
}

// Scale multiplies all vector components by a scalar.
func (v Vec3) Scale(s float64) Vec3 {
	return Vec3{
		I: v.I * s,
		J: v.J * s,
		K: v.K * s,
	}
}

// LengthSquared returns the sum of the squares of the vector components.
func (v Vec3) LengthSquared() float64 {
	return (v.X() * v.X()) + (v.Y() * v.Y()) + (v.Z() * v.Z())
}

// Length returns the Euclidean length of the vector.
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// String returns a string representation of the vector.
func (v Vec3) String() string {
	return "(" + strconv.FormatFloat(v.X(), 'f', -1, 64) + ", " + strconv.FormatFloat(v.Y(), 'f', -1, 64) + ", " + strconv.FormatFloat(v.Z(), 'f', -1, 64) + ")"
}

// Normalize scales the vector to a unit length.
func (v Vec3) Normalize() Vec3 {
	return v.Scale(1.0 / v.Length())
}

// Dot returns the dot product of this vector and another (v · other).
func (v Vec3) Dot(other Vec3) float64 {
	return (v.X() * other.X()) + (v.Y() * other.Y()) + (v.Z() * other.Z())
}

// Cross returns the cross-product of this vector and another (v x other).
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		I: (v.J * other.K) - (v.K * other.J),
		J: (v.K * other.I) - (v.I * other.K),
		K: (v.I * other.J) - (v.J * other.I),
	}
}

// Random returns a vector with random components in [0, 1).
func Random() Vec3 {
	return Vec3{RandomFloat(), RandomFloat(), RandomFloat()}
}

func RandomInRange(i Interval) Vec3 {
	return Vec3{
		RandomFloatInInterval(i),
		RandomFloatInInterval(i),
		RandomFloatInInterval(i),
	}
}

func RandomUnitVec3() Vec3 {
	for {
		p := RandomInRange(Interval{-1, 1})
		l := p.LengthSquared()
		if 1e-160 < l && l <= 1 {
			return p.Normalize()
		}
	}
}

func RandomInHemisphere(normal Vec3) Vec3 {
	unitVector := RandomUnitVec3()
	if unitVector.Dot(normal) > 0 {
		return unitVector
	}
	return unitVector.Negate()
}
