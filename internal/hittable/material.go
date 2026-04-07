package hittable

import (
	"math"
	"raytracer/internal/primitives"
)

type Material interface {
	Scatter(ray primitives.Ray, hitRecord HitRecord) (scattered primitives.Ray, attenuation primitives.Vector, ok bool)
}

type Lambertian struct {
	Albedo primitives.Vector
}

func (m Lambertian) Scatter(ray primitives.Ray, hitRecord HitRecord) (scattered primitives.Ray, attenuation primitives.Vector, ok bool) {
	scatterDirection := hitRecord.Normal.Add(primitives.RandomUnitVector())

	// Catch degenerate scatter direction.
	if scatterDirection.NearZero() {
		scatterDirection = hitRecord.Normal
	}

	scattered = primitives.Ray{Origin: hitRecord.P, Direction: scatterDirection}
	attenuation = m.Albedo
	ok = true
	return
}

type Metal struct {
	Albedo primitives.Vector
	Fuzz   float64
}

func (m Metal) Scatter(ray primitives.Ray, hitRecord HitRecord) (scattered primitives.Ray, attenuation primitives.Vector, ok bool) {
	reflected := ray.Direction.Reflect(hitRecord.Normal)
	reflected = reflected.Normalize().Add(primitives.RandomUnitVector().Scale(m.Fuzz))

	scattered = primitives.Ray{Origin: hitRecord.P, Direction: reflected}
	attenuation = m.Albedo
	ok = scattered.Direction.Dot(hitRecord.Normal) > 0
	return
}

type Dielectric struct {
	RefractionIndex float64
}

func (m Dielectric) Scatter(ray primitives.Ray, hitRecord HitRecord) (scattered primitives.Ray, attenuation primitives.Vector, ok bool) {
	attenuation = primitives.Vector{I: 1, J: 1, K: 1}
	var ri float64
	if hitRecord.FrontFace {
		ri = 1.0 / m.RefractionIndex
	} else {
		ri = m.RefractionIndex
	}

	unitDirection := ray.Direction.Normalize()
	cosTheta := math.Min(unitDirection.Negate().Dot(unitDirection), 1)
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1
	var direction primitives.Vector
	if cannotRefract {
		direction = unitDirection.Reflect(hitRecord.Normal)
	} else {
		direction = unitDirection.Refract(hitRecord.Normal, ri)
	}

	scattered = primitives.Ray{Origin: hitRecord.P, Direction: direction}
	ok = true
	return
}
