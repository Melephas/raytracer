package main

import (
	"math"
)

// Interval represents a range between a minimum and maximum value.
type Interval struct {
	Min, Max float64
}

// DefaultInterval returns an interval spanning from negative to positive infinity.
func DefaultInterval() Interval {
	return Interval{math.Inf(-1), math.Inf(1)}
}

// NewInterval creates and returns a new interval with the given bounds.
func NewInterval(min, max float64) Interval {
	return Interval{min, max}
}

// Size returns the length of the interval.
func (i Interval) Size() float64 {
	return i.Max - i.Min
}

// Contains returns true if the interval contains the value x (inclusive).
func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

// Surrounds returns true if the interval strictly surrounds the value x (exclusive).
func (i Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

// Clamp returns x if it is within the interval, or the nearest bound if it is outside.
func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	} else if x > i.Max {
		return i.Max
	}

	return x
}
