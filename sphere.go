package main

import "math"

// Sphere represents a 3D sphere.
type Sphere struct {
	Center Vec3
	Radius float64
}

// Hit determines if a ray hits the sphere within the given interval.
func (s Sphere) Hit(ray Ray, rayT Interval, hitRecord *HitRecord) bool {
	oc := s.Center.Sub(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.Dot(oc) - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}

	sqrtDiscriminant := math.Sqrt(discriminant)

	// Find the nearest root that lies in the given interval.
	root := (h - sqrtDiscriminant) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtDiscriminant) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	hitRecord.T = root
	hitRecord.P = ray.At(root)
	hitRecord.Normal = hitRecord.P.Sub(s.Center).Scale(1 / s.Radius)
	outwardNormal := hitRecord.P.Sub(s.Center).Scale(1 / s.Radius)
	hitRecord.SetFaceNormal(ray, outwardNormal)
	return true
}
