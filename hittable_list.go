package main

// HittableList maintains a collection of hittable objects.
type HittableList struct {
	Objects []Hittable
}

// NewHittableList creates and returns an empty HittableList.
func NewHittableList() *HittableList {
	return &HittableList{
		Objects: make([]Hittable, 0),
	}
}

// Add appends a hittable object to the list.
func (l *HittableList) Add(hittable Hittable) {
	l.Objects = append(l.Objects, hittable)
}

// Hit determines if a ray hits any object in the list and returns the closest hit.
func (l *HittableList) Hit(ray Ray, rayT Interval, hitRecord *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := rayT.Max

	for _, object := range l.Objects {
		if object.Hit(ray, Interval{rayT.Min, closestSoFar}, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*hitRecord = tempRec
		}
	}

	return hitAnything
}
