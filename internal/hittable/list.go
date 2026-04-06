package hittable

import (
	"raytracer/internal/primitives"
)

// List maintains a collection of hittable objects.
type List struct {
	Objects []Hittable
}

// NewList creates and returns an empty List.
func NewList() *List {
	return &List{
		Objects: make([]Hittable, 0),
	}
}

// Add appends a hittable object to the list.
func (l *List) Add(hittable Hittable) {
	l.Objects = append(l.Objects, hittable)
}

// Hit determines if a ray hits any object in the list and returns the closest hit.
func (l *List) Hit(ray primitives.Ray, rayT primitives.Interval) (*HitRecord, bool) {
	var closestHit *HitRecord
	hitAnything := false
	closestSoFar := rayT.Max

	for _, object := range l.Objects {
		hitRec, hit := object.Hit(ray, primitives.Interval{Min: rayT.Min, Max: closestSoFar})
		if hit {
			hitAnything = true
			closestSoFar = hitRec.T
			closestHit = hitRec
		}
	}

	if hitAnything {
		return closestHit, true
	}

	return nil, false
}
