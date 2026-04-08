package scene

import (
	"encoding/json"
	"io"
	"raytracer/internal/camera"
	"raytracer/internal/hittable"
	"raytracer/internal/primitives"
	"raytracer/internal/shape"
)

type Material struct {
	Type            string     `json:"type"`
	Albedo          [3]float64 `json:"albedo"`
	Fuzz            float64    `json:"fuzz"`
	RefractionIndex float64    `json:"refractionIndex"`
}

type Shape struct {
	Type     string     `json:"type"`
	Center   [3]float64 `json:"center"`
	Radius   float64    `json:"radius"`
	Material string     `json:"material"`
}

type Camera struct {
	Position     [3]float64 `json:"position"`
	LookAt       [3]float64 `json:"lookAt"`
	Up           [3]float64 `json:"up"`
	FOV          float64    `json:"fov"`
	SpaceColor   [3]float64 `json:"spaceColor"`
	GroundColor  [3]float64 `json:"groundColor"`
	DefocusAngle float64    `json:"defocusAngle"`
	FocusDist    float64    `json:"focusDist"`
}

type Manager struct {
	Materials map[string]Material `json:"materials"`
	Scene     []Shape             `json:"scene"`
	Camera    Camera              `json:"camera"`
}

func LoadScene(r io.Reader) (*Manager, error) {
	var manager Manager
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &manager); err != nil {
		return nil, err
	}
	return &manager, nil
}

func (m Material) ToMaterial() hittable.Material {
	switch m.Type {
	case "lambertian":
		return hittable.Lambertian{
			Albedo: ConvertToVector(m.Albedo),
		}
	case "metal":
		return hittable.Metal{
			Albedo: ConvertToVector(m.Albedo),
			Fuzz:   m.Fuzz,
		}
	case "dielectric":
		return hittable.Dielectric{
			RefractionIndex: m.RefractionIndex,
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
			Center:   ConvertToVector(s.Center),
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

func (m *Manager) GetCamera() *camera.Camera {
	cam := camera.DefaultCamera()
	cam.Position = ConvertToVector(m.Camera.Position)
	cam.LookAt = ConvertToVector(m.Camera.LookAt)
	cam.Up = ConvertToVector(m.Camera.Up)
	cam.FOV = m.Camera.FOV
	cam.SpaceColor = ConvertToVector(m.Camera.SpaceColor)
	cam.GroundColor = ConvertToVector(m.Camera.GroundColor)
	cam.DefocusAngle = m.Camera.DefocusAngle
	cam.FocusDist = m.Camera.FocusDist
	return cam
}

func ConvertToVector(v [3]float64) primitives.Vector {
	return primitives.Vector{
		I: v[0],
		J: v[1],
		K: v[2],
	}
}
