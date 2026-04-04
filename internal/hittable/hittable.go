package hittable

import (
	"raytracer/internal/primitives"
)

// HitRecord stores details of a ray hitting an object.
type HitRecord struct {
	P, Normal primitives.Vector
	T         float64
	FrontFace bool
	Material  Material
}

// SetFaceNormal determines if a hit is on the front or back face and sets the normal accordingly.
func (r HitRecord) SetFaceNormal(ray primitives.Ray, outwardNormal primitives.Vector) HitRecord {
	ret := r

	frontFace := ray.Direction.Dot(outwardNormal) < 0
	if frontFace {
		ret.Normal = outwardNormal
	} else {
		ret.Normal = outwardNormal.Negate()
	}

	return ret
}

// Hittable defines the interface for objects that can be hit by a ray.
type Hittable interface {
	Hit(ray primitives.Ray, rayT primitives.Interval) (*HitRecord, bool)
}
