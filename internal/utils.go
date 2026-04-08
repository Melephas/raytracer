package internal

import (
	"math"
	"math/rand/v2"
)

func LinearToGamma(x float64) float64 {
	gamma := 1.7
	if x > 0 {
		return math.Pow(x, 1/gamma)
	}

	return 0
}

// DegreesToRadians converts degrees to radians.
func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

//// RadiansToDegrees converts radians to degrees.
//func RadiansToDegrees(radians float64) float64 {
//	return radians * (180 / math.Pi)
//}

// RandomFloat returns a random float in [0, 1).
func RandomFloat() float64 {
	return rand.Float64()
}

func RandomFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
