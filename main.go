package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	// AspectRatio is the target aspect ratio for the rendered image.
	AspectRatio = 16.0 / 9.0
	// ImageWidth is the width of the rendered image in pixels.
	ImageWidth = 1920
	//ImageWidth  = 1280
)

// OutputFile is the name of the file to which the render will be saved.
var OutputFile string

func main() {
	// Open output file and schedule cleanup.
	outputWriter, closeFunc, err := GetOutputWriter()
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer closeFunc()

	// World setup.
	world := NewHittableList()
	world.Add(Sphere{Vec3{0, 0, -1}, 0.5})
	world.Add(Sphere{Vec3{0, -100.5, -1}, 100})

	// Camera setup.
	camera := DefaultCamera()
	camera.AspectRatio = AspectRatio
	camera.ImageWidth = ImageWidth
	camera.SamplesPerPixel = 100

	if err := camera.Render(outputWriter, world); err != nil {
		log.Fatalf("Error rendering image: %v", err)
	}

	// Flush output file.
	if err := outputWriter.Flush(); err != nil {
		log.Fatalf("Error flushing output: %v", err)
	}

	fmt.Printf("Render complete. Output written to %s\n", OutputFile)
}

// GetOutputWriter returns a writer for the output file and a cleanup function to close the file.
func GetOutputWriter() (*bufio.Writer, func(), error) {
	if len(os.Args) >= 2 {
		OutputFile = os.Args[1]
	}

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
