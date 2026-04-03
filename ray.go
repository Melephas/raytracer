package main

// Ray represents a directed line with an origin and direction.
type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// At calculates the position along the ray at a given distance t.
func (r Ray) At(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
