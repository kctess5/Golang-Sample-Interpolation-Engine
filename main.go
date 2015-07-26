package main

import (
	"./sampler"
	"./utils"
	"fmt"
	"math"
	"time"
)

func main() {
	start := time.Now()

	width := 400
	height := 400

	fs := sampler.GaussianFrameSampler(width, height, 4, 1, 12)

	for x := 0; x < width; x += 10 {
		for y := 0; y < height; y += 10 {
			fs.AddSample(
				sampler.Coord(x, y),
				[]float64{
					255 * math.Sin(float64(x)/float64(width)),
					10,
					255 * math.Sin(float64(y)/float64(height)),
					255},
			)
		}
	}

	utils.SaveImage(width, height, fs.Rasterize(), "test.png")

	fmt.Println("Saved Image after:", time.Since(start), "Seconds")
}
