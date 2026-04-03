package main

// HitRecord stores details of a ray hitting an object.
type HitRecord struct {
	P, Normal Vec3
	T         float64
	FrontFace bool
}

// SetFaceNormal determines if a hit is on the front or back face and sets the normal accordingly.
func (r *HitRecord) SetFaceNormal(ray Ray, outwardNormal Vec3) {
	r.FrontFace = ray.Direction.Dot(outwardNormal) < 0
	if r.FrontFace {
		r.Normal = outwardNormal
	} else {
		r.Normal = outwardNormal.Negate()
	}
}

// Hittable defines the interface for objects that can be hit by a ray.
type Hittable interface {
	Hit(ray Ray, rayT Interval, hitRecord *HitRecord) bool
}
