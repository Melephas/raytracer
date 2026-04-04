package hittable

import "raytracer/internal/primitives"

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
