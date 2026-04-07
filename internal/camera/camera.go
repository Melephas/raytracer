package camera

import (
	"fmt"
	"io"
	"log"
	"math"
	"raytracer/internal"
	"raytracer/internal/hittable"
	"raytracer/internal/primitives"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// Camera represents a virtual camera that renders a scene.
type Camera struct {
	Position, FirstPixel      primitives.Vector
	PixelDV, PixelDU          primitives.Vector
	Up, LookAt                primitives.Vector
	U, V, W                   primitives.Vector
	AspectRatio, SamplesScale float64
	ImageWidth, ImageHeight   int
	SamplesPerPixel, MaxDepth int
	Parallel                  bool
	FOV                       float64
	SpaceColor                primitives.Vector
}

// DefaultCamera returns a camera with default values.
func DefaultCamera() *Camera {
	return &Camera{
		Position:        primitives.Vector{},
		Up:              primitives.Vector{J: 1},
		LookAt:          primitives.Vector{K: -1},
		AspectRatio:     16.0 / 9.0,
		ImageWidth:      1920,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		FOV:             90,
	}
}

// Render renders the scene to the provided output writer using the given world.
func (c *Camera) Render(out io.Writer, world hittable.Hittable) error {
	//c.Initialise()

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
		//progressbar.OptionClearOnFinish(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	if c.Parallel {
		var waitGroup sync.WaitGroup
		var mutex sync.Mutex
		output := make([]primitives.Vector, c.ImageHeight*c.ImageWidth)
		calculatePixel := func(i, j int, mutex *sync.Mutex, output []primitives.Vector) {
			// calculatePixel := func() {
			pixelColor := primitives.Vector{}
			for range c.SamplesPerPixel {
				r := c.GetRay(i, j)
				pixelColor = pixelColor.Add(c.RayColor(r, c.MaxDepth, world))
			}
			finalColor := pixelColor.Scale(c.SamplesScale)

			mutex.Lock()
			output[i+(j*c.ImageWidth)] = finalColor
			mutex.Unlock()
		}

		for j := range c.ImageHeight {
			for i := range c.ImageWidth {
				waitGroup.Go(func() {
					calculatePixel(i, j, &mutex, output)
				})
			}
			if err := progressBar.Add(1); err != nil {
				return err
			}
		}

		waitGroup.Wait()

		mutex.Lock()
		for i := range len(output) {
			if err := primitives.WriteColor(out, output[i]); err != nil {
				log.Fatalf("Error writing color to file: %v", err)
			}
		}
		mutex.Unlock()
	} else {
		for j := range c.ImageHeight {
			for i := range c.ImageWidth {
				pixelColor := primitives.Vector{}
				for range c.SamplesPerPixel {
					r := c.GetRay(i, j)
					pixelColor = pixelColor.Add(c.RayColor(r, c.MaxDepth, world))
				}

				if err := primitives.WriteColor(out, pixelColor.Scale(c.SamplesScale)); err != nil {
					log.Fatalf("Error writing color to file: %v", err)
				}
			}
			if err := progressBar.Add(1); err != nil {
				return err
			}
		}
	}

	return nil
}

// Initialise calculates the derived camera parameters.
func (c *Camera) Initialise() {
	imageHeight := int(float64(c.ImageWidth) / c.AspectRatio)
	c.ImageHeight = max(imageHeight, 1)

	c.SamplesScale = 1.0 / float64(c.SamplesPerPixel)

	// Determine viewport dimensions.
	focalLength := (c.Position.Sub(c.LookAt)).Length()
	theta := internal.DegreesToRadians(c.FOV)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h * focalLength
	viewportWidth := float64(c.ImageWidth) / float64(c.ImageHeight) * viewportHeight

	// Calculate the camera's basis vectors.
	c.W = c.Position.Sub(c.LookAt).Normalize()
	c.U = c.Up.Cross(c.W).Normalize()
	c.V = c.W.Cross(c.U) // |W X U| = 1 for normalized vectors

	// Calculate the viewport basis vectors.
	viewportU := c.U.Scale(viewportWidth)
	viewportV := c.V.Scale(viewportHeight).Negate()

	// Calculate the delta vectors for each pixel.
	c.PixelDU = viewportU.Scale(1.0 / float64(c.ImageWidth))
	c.PixelDV = viewportV.Scale(1.0 / float64(c.ImageHeight))

	// Calculate the position of the upper left pixel. WARNING: I have no idea why this works.
	viewportUpperLeft := c.Position.Sub(c.W.Scale(focalLength)).Sub(viewportU.Scale(0.5)).Sub(viewportV.Scale(0.5))
	c.FirstPixel = viewportUpperLeft.Add(c.PixelDU.Add(c.PixelDV).Scale(0.5))
}

// RayColor computes the color for a given ray in the world.
func (c *Camera) RayColor(r primitives.Ray, depth int, world hittable.Hittable) primitives.Vector {
	if depth <= 0 {
		return primitives.Vector{}
	}

	rec, hit := world.Hit(r, primitives.Interval{Min: 0.001, Max: math.Inf(1)})
	if hit {
		scattered, attenuation, ok := rec.Material.Scatter(r, *rec)
		if !ok {
			return primitives.Vector{}
		}
		return attenuation.ColorMultiply(c.RayColor(scattered, depth-1, world))
	}

	// Sky color.
	//unitDirection := r.Direction.Normalize()
	//a := (unitDirection.Y() + 1.0) * 0.5
	//return primitives.Vector{I: 1, J: 1, K: 1}.Scale(1 - a).Add(primitives.Vector{I: 0.5, J: 0.7, K: 1.0}.Scale(a))
	return c.SpaceColor
}

// GetRay returns a new ray for the pixel at (i, j) with random sampling.
func (c *Camera) GetRay(i, j int) primitives.Ray {
	offset := c.SampleSquare()
	pixelSample := c.FirstPixel.Add(c.PixelDU.Scale(float64(i) + offset.X())).Add(c.PixelDV.Scale(float64(j) + offset.Y()))

	rayOrigin := c.Position
	rayDirection := pixelSample.Sub(rayOrigin).Normalize()

	return primitives.Ray{Origin: rayOrigin, Direction: rayDirection}
}

// SampleSquare returns a random 2D vector within a [-0.5, 0.5] unit square.
func (c *Camera) SampleSquare() primitives.Vector {
	return primitives.Vector{I: internal.RandomFloat() - 0.5, J: internal.RandomFloat() - 0.5}
}

// String returns a string representation of the Camera.
func (c *Camera) String() string {
	return fmt.Sprintf("Camera {\n\tPosition: %v\n\tFirstPixel: %v\n\tPixelDV: %v\n\tPixelDU: %v\n\tUp: %v\n\tLookAt: %v\n\tU: %v\n\tV: %v\n\tW: %v\n\tAspectRatio: %v\n\tSamplesScale: %v\n\tImageWidth: %d\n\tImageHeight: %d\n\tSamplesPerPixel: %d\n\tMaxDepth: %d\n\tParallel: %v\n\tFOV: %v\n}",
		c.Position, c.FirstPixel, c.PixelDV, c.PixelDU, c.Up, c.LookAt, c.U, c.V, c.W, c.AspectRatio, c.SamplesScale, c.ImageWidth, c.ImageHeight, c.SamplesPerPixel, c.MaxDepth, c.Parallel, c.FOV)
}
