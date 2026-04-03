package primitives

import (
	"fmt"
	"io"
	"math"
	"raytracer/internal"
	"strconv"
)

// Vector represents a 3D Vector or color.
type Vector struct {
	I, J, K float64
}

//type Color Vec3
//type Point Vec3

// --- Getter functions ---

// X returns the first component of the Vector.
func (v Vector) X() float64 {
	return v.I
}

// Y returns the second component of the Vector.
func (v Vector) Y() float64 {
	return v.J
}

// Z returns the third component of the Vector.
func (v Vector) Z() float64 {
	return v.K
}

// R returns the first component of the Vector, alias for X.
func (v Vector) R() float64 {
	return v.I
}

// G returns the second component of the Vector, alias for Y.
func (v Vector) G() float64 {
	return v.J
}

// B returns the third component of the Vector, alias for Z.
func (v Vector) B() float64 {
	return v.K
}

// --- Arithmetic functions ---

// Negate returns a flipped version of the Vector.
func (v Vector) Negate() Vector {
	return Vector{
		I: -v.I,
		J: -v.J,
		K: -v.K,
	}
}

// Add returns a new Vector that is the sum of the two vectors.
func (v Vector) Add(other Vector) Vector {
	return Vector{
		I: v.I + other.I,
		J: v.J + other.J,
		K: v.K + other.K,
	}
}

// Sub returns a Vector that is the sum of the two vectors, with the second primitives flipped (v - other).
func (v Vector) Sub(other Vector) Vector {
	return v.Add(other.Negate())
}

// Scale multiplies all Vector components by a scalar.
func (v Vector) Scale(s float64) Vector {
	return Vector{
		I: v.I * s,
		J: v.J * s,
		K: v.K * s,
	}
}

// LengthSquared returns the sum of the squares of the Vector components.
func (v Vector) LengthSquared() float64 {
	return (v.X() * v.X()) + (v.Y() * v.Y()) + (v.Z() * v.Z())
}

// Length returns the Euclidean length of the Vector.
func (v Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// String returns a string representation of the Vector.
func (v Vector) String() string {
	return "(" + strconv.FormatFloat(v.X(), 'f', -1, 64) + ", " + strconv.FormatFloat(v.Y(), 'f', -1, 64) + ", " + strconv.FormatFloat(v.Z(), 'f', -1, 64) + ")"
}

// Normalize scales the Vector to a unit length.
func (v Vector) Normalize() Vector {
	return v.Scale(1.0 / v.Length())
}

// Dot returns the dot product of this Vector and another (v · other).
func (v Vector) Dot(other Vector) float64 {
	return (v.X() * other.X()) + (v.Y() * other.Y()) + (v.Z() * other.Z())
}

// Cross returns the cross-product of this Vector and another (v x other).
func (v Vector) Cross(other Vector) Vector {
	return Vector{
		I: (v.J * other.K) - (v.K * other.J),
		J: (v.K * other.I) - (v.I * other.K),
		K: (v.I * other.J) - (v.J * other.I),
	}
}

// Random returns a Vector with random components in [0, 1).
func Random() Vector {
	return Vector{internal.RandomFloat(), internal.RandomFloat(), internal.RandomFloat()}
}

func RandomInRange(i Interval) Vector {
	return Vector{
		RandomFloatInInterval(i),
		RandomFloatInInterval(i),
		RandomFloatInInterval(i),
	}
}

func RandomUnitVec3() Vector {
	for {
		p := RandomInRange(Interval{-1, 1})
		l := p.LengthSquared()
		if 1e-160 < l && l <= 1 {
			return p.Normalize()
		}
	}
}

func RandomInHemisphere(normal Vector) Vector {
	unitVector := RandomUnitVec3()
	if unitVector.Dot(normal) > 0 {
		return unitVector
	}
	return unitVector.Negate()
}

// WriteColor writes a color to the given writer.
func WriteColor(w io.Writer, color Vector) error {
	intensity := Interval{Min: 0.0, Max: 0.999}
	r := int(256 * intensity.Clamp(internal.LinearToGamma(color.X())))
	g := int(256 * intensity.Clamp(internal.LinearToGamma(color.Y())))
	b := int(256 * intensity.Clamp(internal.LinearToGamma(color.Z())))

	if _, err := fmt.Fprintf(w, "%d %d %d\n", r, g, b); err != nil {
		return err
	}

	return nil
}
