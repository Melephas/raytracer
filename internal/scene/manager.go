package scene

import (
	"encoding/json/v2"
	"io"
	"raytracer/internal/hittable"
	"raytracer/internal/primitives"
	"raytracer/internal/shape"
)

type Material struct {
	Type   string     `json:"type"`
	Albedo [3]float64 `json:"albedo"`
	Fuzz   float64    `json:"fuzz"`
}

type Shape struct {
	Type     string     `json:"type"`
	Center   [3]float64 `json:"center"`
	Radius   float64    `json:"radius"`
	Material string     `json:"material"`
}

type Manager struct {
	Materials map[string]Material `json:"materials"`
	Scene     []Shape             `json:"scene"`
}

func LoadScene(r io.Reader) (*Manager, error) {
	var manager Manager
	if err := json.UnmarshalRead(r, &manager); err != nil {
		return nil, err
	}
	return &manager, nil
}

func (m Material) ToMaterial() hittable.Material {
	switch m.Type {
	case "lambertian":
		return hittable.Lambertian{
			Albedo: primitives.Vector{
				I: m.Albedo[0],
				J: m.Albedo[1],
				K: m.Albedo[2],
			},
		}
	case "metal":
		return hittable.Metal{
			Albedo: primitives.Vector{
				I: m.Albedo[0],
				J: m.Albedo[1],
				K: m.Albedo[2],
			},
			Fuzz: m.Fuzz,
		}
	default:
		return nil
	}
}

func (s Shape) ToHittable(manager *Manager) hittable.Hittable {
	material := manager.Materials[s.Material]
	if material.Type == "" { // If this material is not defined, error out.
		return nil
	}

	switch s.Type {
	case "sphere":
		return shape.Sphere{
			Center: primitives.Vector{
				I: s.Center[0],
				J: s.Center[1],
				K: s.Center[2],
			},
			Radius:   s.Radius,
			Material: material.ToMaterial(),
		}
	default:
		return nil
	}
}

func (m *Manager) ToHittables() hittable.Hittable {
	list := hittable.NewList()
	for _, s := range m.Scene {
		newHittable := s.ToHittable(m)
		list.Add(newHittable)
	}
	return list
}
