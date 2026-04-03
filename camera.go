package main

import (
	"fmt"
	"io"
	"log"
	"math"

	"github.com/schollz/progressbar/v3"
)

// Camera represents a virtual camera that renders a scene.
type Camera struct {
	Position, FirstPixel      Vec3
	PixelDV, PixelDU          Vec3
	AspectRatio, SamplesScale float64
	ImageWidth, ImageHeight   int
	SamplesPerPixel           int
}

// DefaultCamera returns a camera with default values.
func DefaultCamera() *Camera {
	return &Camera{
		Position:        Vec3{0, 0, 0},
		AspectRatio:     AspectRatio,
		ImageWidth:      ImageWidth,
		SamplesPerPixel: 100,
	}
}

// Render renders the scene to the provided output writer using the given world.
func (c *Camera) Render(out io.Writer, world Hittable) error {
	c.Initialise()

	// Write output file header.
	if _, err := fmt.Fprintf(out, "P3\n%d %d\n255\n", c.ImageWidth, c.ImageHeight); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	// Render image and write pixel colors.
	progressBar := progressbar.NewOptions(c.ImageHeight,
		progressbar.OptionSetDescription("Rendering"),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for j := range c.ImageHeight {
		for i := range c.ImageWidth {
			pixelColor := Vec3{0, 0, 0}
			for range c.SamplesPerPixel {
				r := c.GetRay(i, j)
				pixelColor = pixelColor.Add(c.RayColor(r, world))
			}

			if err := WriteColor(out, pixelColor.Scale(c.SamplesScale)); err != nil {
				log.Fatalf("Error writing color to file: %v", err)
			}
		}
		if err := progressBar.Add(1); err != nil {
			return err
		}
	}

	return nil
}

// Initialise calculates the derived camera parameters.
func (c *Camera) Initialise() {
	imageHeight := int(float64(c.ImageWidth) / c.AspectRatio)
	if imageHeight < 1 {
		c.ImageHeight = 1
	} else {
		c.ImageHeight = imageHeight
	}

	c.SamplesScale = 1.0 / float64(c.SamplesPerPixel)

	c.Position = Vec3{0, 0, 0}

	// Determine viewport dimensions.
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := float64(c.ImageWidth) / float64(c.ImageHeight) * viewportHeight

	// Calculate the viewport basis vectors.
	viewportU := Vec3{viewportWidth, 0.0, 0.0}
	viewportV := Vec3{0.0, -viewportHeight, 0.0}

	// Calculate the delta vectors for each pixel.
	c.PixelDU = viewportU.Scale(1.0 / float64(c.ImageWidth))
	c.PixelDV = viewportV.Scale(1.0 / float64(c.ImageHeight))

	// Calculate the position of the upper left pixel. WARNING: I have no idea why this works.
	viewportUpperLeft := c.Position.Sub(Vec3{0, 0, focalLength}).Sub(viewportU.Scale(0.5)).Sub(viewportV.Scale(0.5))
	c.FirstPixel = viewportUpperLeft.Add(c.PixelDU.Add(c.PixelDV).Scale(0.5))
}

// RayColor computes the color for a given ray in the world.
func (c *Camera) RayColor(r Ray, world Hittable) Vec3 {
	var rec HitRecord
	if world.Hit(r, Interval{0, math.Inf(1)}, &rec) {
		return rec.Normal.Add(Vec3{1, 1, 1}).Scale(0.5)
	}

	unitDirection := r.Direction.Normalize()
	a := (unitDirection.Y() + 1.0) * 0.5
	return Vec3{1, 1, 1}.Scale(1 - a).Add(Vec3{0.5, 0.7, 1.0}.Scale(a))
}

// GetRay returns a new ray for the pixel at (i, j) with random sampling.
func (c *Camera) GetRay(i, j int) Ray {
	offset := c.SampleSquare()
	pixelSample := c.FirstPixel.Add(c.PixelDU.Scale(float64(i) + offset.X())).Add(c.PixelDV.Scale(float64(j) + offset.Y()))

	rayOrigin := c.Position
	rayDirection := pixelSample.Sub(rayOrigin).Normalize()

	return Ray{rayOrigin, rayDirection}
}

// SampleSquare returns a random 2D vector within a [-0.5, 0.5] unit square.
func (c *Camera) SampleSquare() Vec3 {
	return Vec3{RandomFloat() - 0.5, RandomFloat() - 0.5, 0}
}
