# Ray Tracer

A ray tracer written in Go. This project implements a path tracing algorithm to generate photorealistic images from scene descriptions provided in JSON format.

This is a personal project for learning and experimenting with Go.

## Example

![output](https://github.com/user-attachments/assets/10a8df2f-2fa2-4af1-bbe4-5718515bb8d3)

Generated with the example scene at 1024 samples-per-pixel

## Features

- **Concurrent Rendering**: Fully leverages multicore processors for significantly faster rendering times using Go's goroutines (~5x speedup).
- **JSON-Driven Scenes**: Define materials and shapes in a clean, human-readable JSON format.
- **Materials Support**:
  - **Lambertian (Diffuse)**: For non-reflective, matte surfaces.
  - **Metal**: For polished, reflective surfaces with adjustable atmospheric "fuzziness."
- **Configurable Parameters**: Control image resolution, samples per pixel (for antialiasing), and ray bounce depth via CLI flags.
- **PPM Output**: Generates images in the simple and widely supported Portable Pixmap (PPM) format.

## Future Features

- **GPU Acceleration**: Leverage GPU acceleration for even faster rendering times, especially for complex scenes.
- **Dielectric Materials**: Support for transparent materials like glass and water.
- **Material Library**: Allow users to define custom materials and reuse them across scenes.
- **Additional Shapes**: Support for shapes that aren't just spheres.
- **Other Image Formats**: Support for more image formats.

## Getting Started

### Prerequisites

- **Go 1.25** or later is required to build and run the project.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/raytracer.git
   cd raytracer
   ```

2. Build the executable:
   ```bash
   go build -o raytracer main.go
   ```

## Usage

Run the ray tracer by specifying an input scene file and various rendering options.

### Commands

```bash
./raytracer -i scene.json -o output.ppm -w 1280 -s 100 -b 50 -p
```

### Options

| Flag | Default      | Description                                                                              |
|:-----|:-------------|:-----------------------------------------------------------------------------------------|
| `-i` | `scene.json` | Path to the input JSON scene file.                                                       |
| `-o` | `output.ppm` | Name of the output image file.                                                           |
| `-w` | `1280`       | Width of the rendered image in pixels (height is calculated based on 16:9 aspect ratio). |
| `-s` | `10`         | Number of samples per pixel (higher values reduce noise but increase render time).       |
| `-b` | `10`         | Maximum number of ray bounces (depth).                                                   |
| `-p` | `false`      | Enable parallel rendering for multi-core performance.                                    |

## Scene Configuration

Scenes are defined in JSON, allowing you to easily manage materials and objects.

### Example `scene.json`

```json
{
    "materials": {
        "silver": {
            "type": "metal",
            "albedo": [0.8, 0.8, 0.8],
            "fuzz": 0.05
        },
        "ground": {
            "type": "lambertian",
            "albedo": [0.5, 0.5, 0.5]
        }
    },
    "scene": [
        {
            "type": "sphere",
            "center": [0, -100.5, -1],
            "radius": 100,
            "material": "ground"
        },
        {
            "type": "sphere",
            "center": [0, 0, -1],
            "radius": 0.5,
            "material": "silver"
        }
    ]
}
```

## How It Works

This program uses a **Path Tracing** approach:

1. **Ray Generation**: For every pixel, rays are cast from the camera into the scene. To achieve anti-aliasing, multiple rays are cast per pixel with slight random offsets (`samples per pixel`).
2. **Intersection Testing**: Each ray is checked for intersections with all objects in the scene. The program currently supports spheres as the primary primitive.
3. **Material Interaction**: When a ray hits an object, the material determines how it scatters:
   - **Lambertian** materials scatter rays in random directions, weighted by the surface normal.
   - **Metal** materials reflect rays based on the angle of incidence, with optional "fuzz" to simulate rougher surfaces.
4. **Recursion**: Scanned rays continue to bounce until they either hit a light source (the background sky color in this case) or reach the maximum `bounce depth`.
5. **Color Accumulation**: The colors from all samples are averaged, gamma-corrected (using a factor of 1.7), and then written to the output file.

## Acknowledgments

This project is inspired by the [_Ray Tracing in One Weekend_](https://raytracing.github.io/books/RayTracingInOneWeekend.html) series.

## License

MIT License. See `LICENSE` for more information.
