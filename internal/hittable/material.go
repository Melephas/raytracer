package hittable

import (
	"math"
	"math/rand"
	"raytracer/internal/primitives"
)

type Material interface {
	Scatter(ray primitives.Ray, hitRecord HitRecord) (scattered primitives.Ray, attenuation primitives.Vector, ok bool)
}

type Lambertian struct {
	Albedo primitives.Vector
}

func (m Lambertian) Scatter(_ primitives.Ray, hitRecord HitRecord) (scattered primitives.Ray, attenuation primitives.Vector, ok bool) {
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

func (m Dielectric) Scatter(ray primitives.Ray, hitRecord HitRecord) (primitives.Ray, primitives.Vector, bool) {
	attenuation := primitives.Vector{I: 1.0, J: 1.0, K: 1.0}
	var refractionRatio float64
	if hitRecord.FrontFace {
		refractionRatio = 1.0 / m.RefractionIndex
	} else {
		refractionRatio = m.RefractionIndex
	}

	unitDirection := ray.Direction.Normalize()
	cosTheta := math.Min(unitDirection.Negate().Dot(hitRecord.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1.0
	var direction primitives.Vector

	if cannotRefract || reflectance(cosTheta, refractionRatio) > rand.Float64() {
		direction = unitDirection.Reflect(hitRecord.Normal)
	} else {
		direction = unitDirection.Refract(hitRecord.Normal, refractionRatio)
	}

	scattered := primitives.Ray{Origin: hitRecord.P, Direction: direction}
	return scattered, attenuation, true
}

func reflectance(cosine, refractionRatio float64) float64 {
	// Use Schlick's approximation for reflectance.
	r0 := (1 - refractionRatio) / (1 + refractionRatio)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
