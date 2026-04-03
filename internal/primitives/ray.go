package primitives

// Ray represents a directed line with an origin and direction.
type Ray struct {
	Origin    Vector
	Direction Vector
}

// At calculates the position along the ray at a given distance t.
func (r Ray) At(t float64) Vector {
	return r.Origin.Add(r.Direction.Scale(t))
}
