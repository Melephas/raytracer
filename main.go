package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"raytracer/internal/camera"
	"raytracer/internal/hittable"
	"raytracer/internal/primitives"
	"raytracer/internal/shape"
)

const (
	// AspectRatio is the target aspect ratio for the rendered image.
	AspectRatio = 16.0 / 9.0
	// ImageWidth is the width of the rendered image in pixels.
	//ImageWidth = 1920
	ImageWidth = 1280
)

// OutputFile is the name of the file to which the render will be saved.
var (
	OutputFile      string
	SamplesPerPixel int
)

func init() {
	flag.IntVar(&SamplesPerPixel, "s", 10, "Number of samples per pixel")
	flag.StringVar(&OutputFile, "o", "output.ppm", "Name of the output file")
}

func main() {
	flag.Parse()

	// Open output file and schedule cleanup.
	outputWriter, closeFunc, err := GetOutputWriter()
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer closeFunc()

	// Material setup.
	silver := hittable.Metal{Albedo: primitives.Vector{I: 0.8, J: 0.8, K: 0.8}, Fuzz: 0.1}
	gold := hittable.Metal{Albedo: primitives.Vector{I: 0.8, J: 0.65, K: 0.1}, Fuzz: 0.9}
	gray := hittable.Lambertian{Albedo: primitives.Vector{I: 0.5, J: 0.5, K: 0.5}}
	green := hittable.Lambertian{Albedo: primitives.Vector{I: 81.0 / 256.0, J: 214.0 / 256.0, K: 84.0 / 256.0}}

	// World setup.
	world := hittable.NewList()
	world.Add(shape.Sphere{Center: primitives.Vector{I: 0, J: 0.5, K: -1.5}, Radius: 0.5, Material: gray})
	world.Add(shape.Sphere{Center: primitives.Vector{I: -1, J: 0, K: -1}, Radius: 0.5, Material: silver})
	world.Add(shape.Sphere{Center: primitives.Vector{I: 1, J: 0, K: -1}, Radius: 0.5, Material: gold})
	world.Add(shape.Sphere{Center: primitives.Vector{I: 0, J: -100.5, K: -1}, Radius: 100, Material: green})

	// Camera setup.
	cam := camera.DefaultCamera()
	cam.AspectRatio = AspectRatio
	cam.ImageWidth = ImageWidth
	cam.SamplesPerPixel = SamplesPerPixel

	if err := cam.Render(outputWriter, world); err != nil {
		log.Fatalf("\nError rendering image: %v", err)
	}
	fmt.Println()

	// Flush output file.
	if err := outputWriter.Flush(); err != nil {
		log.Fatalf("Error flushing output: %v", err)
	}

	fmt.Printf("Render complete. Output written to %s\n", OutputFile)
}

// GetOutputWriter returns a writer for the output file and a cleanup function to close the file.
func GetOutputWriter() (*bufio.Writer, func(), error) {
	if OutputFile == "" {
		OutputFile = "output.ppm"
	}

	file, err := os.OpenFile(OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}

	closeFunc := func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}

	return bufio.NewWriter(file), closeFunc, nil
}
