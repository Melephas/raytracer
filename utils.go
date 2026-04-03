package main

import (
	"fmt"
	"io"
	"math"
	"math/rand/v2"
)

// WriteColor writes a color to the given writer.
func WriteColor(w io.Writer, color Vec3) error {
	intensity := Interval{0.0, 0.999}
	r := int(256 * intensity.Clamp(color.X()))
	g := int(256 * intensity.Clamp(color.Y()))
	b := int(256 * intensity.Clamp(color.Z()))

	if _, err := fmt.Fprintf(w, "%d %d %d\n", r, g, b); err != nil {
		return err
	}

	return nil
}

// DegreesToRadians converts degrees to radians.
func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

// RadiansToDegrees converts radians to degrees.
func RadiansToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

// RandomFloat returns a random float in [0, 1).
func RandomFloat() float64 {
	return rand.Float64()
}

// RandomFloatInInterval returns a random float in the given interval, i.e. [min, max).
func RandomFloatInInterval(i Interval) float64 {
	return i.Min + RandomFloat()*(i.Size())
}
