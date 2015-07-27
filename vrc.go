package main

import (
	"./raycasting"
	"./sampler"
	"./utils"
	"fmt"
	"time"
)

const width = 200.0
const height = 200.0

const max_d = 2

func baseline(sub int, fs *sampler.FrameSampler, scene *raycasting.Scene) {
	for x := 0; x < fs.Width(); x += 1 / sub {
		for y := 0; y < fs.Height(); y += 1 / sub {

		}
	}
}

func main() {
	fmt.Print("Initializing Ray Casting Engine.")

	cpus := 8
	// runtime.GOMAXPROCS(cpus)
	// utils.StartNetProfile("profile")
	// rand.Seed(time.Now().UTC().UnixNano())
	start := time.Now()

	scene := raycasting.NewScene()
	scene.Add(raycasting.Sphere())

	fs := sampler.GaussianFrameSampler(int(width), int(height), 4, 1,
		float64(max_d), cpus)

	baseline(1, fs, scene)

	raster := fs.Rasterize()
	utils.SaveImage(int(width), int(height), raster, "test_vrc")
	fmt.Println("Finished in:", time.Since(start))
}
