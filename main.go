package main

import (
	"./sampler"
	"./utils"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func PSin(i float64) float64 {
	return (math.Sin(i) + 1) / 2
}
func ConvColor(c color.Color) utils.Color {
	r, g, b, a := c.RGBA()

	return utils.NewColor(float64(r), float64(g), float64(b), float64(a))
}

const width = 640
const height = 480

// const width = 1280
// const height = 960

// const width = 1920
// const height = 1080

const scale = 1.0

const x_step = 30 * scale
const y_step = 30 * scale

const iterations = 1
const max_d = 70 * scale
const positionChaos = 20 * scale

const gauss = 1.5

func f_1(c utils.Coord, x_c, y_c float64) float64 {
	x, y := c.X(), c.Y()
	return math.Cos(x_c*x) - math.Cos(y_c*y)
}

func f_2(c utils.Coord, x_c, y_c float64) float64 {
	x, y := c.X(), c.Y()
	return math.Sin(y*x/x_c) + y*math.Cos(y/y_c)
}

func transform(c utils.Coord, width, height float64) utils.Coord {
	return utils.Randomize(c, positionChaos)
	// return c
}

func colorize(c utils.Coord, width, height float64) utils.Color {

	c_2 := utils.Swirl(c, width/2, height/2, math.Pi, height)

	x, y := c_2.X(), c_2.Y()
	_, h, wh := width, height, width*height

	// x_2 := x + 22
	// y_2 := y + 22

	r := 0.0 * math.Sin(x*y/wh)
	g := 0.0 * math.Sin(y/h)
	b := f_2(c, 2*width, 2*height)
	a := 1.0

	return utils.NewColorFrac(r, g, b, a)
}

func main() {
	fmt.Print("Initializing Rendering Engine. ")

	cpus := 8
	// cpus := 1
	runtime.GOMAXPROCS(cpus)
	// utils.StartNetProfile("profile")

	rand.Seed(time.Now().UTC().UnixNano())
	start := time.Now()

	s_gauss := gauss * math.Pow(scale, 1)
	s_height, s_width := scale*height, scale*width

	fs := sampler.GaussianFrameSampler(int(s_width), int(s_height), 4, s_gauss,
		float64(max_d), cpus)

	fmt.Println("Adding samples using:", cpus, "CPUs.")

	var wg sync.WaitGroup

	for i := 0; i < cpus; i++ {
		wg.Add(1)

		go func(i, max, width, height int) {
			for x := x_step * float64(i); x < float64(width)*scale+x_step; x += x_step * float64(max) {
				for y := y_step; y < float64(height)*scale+y_step; y += y_step {

					coords := utils.NewCoord(x, y)

					for i := 0.0; i < iterations; i++ {
						p := transform(coords, s_width, s_height)
						c := colorize(coords, s_width, s_height)

						fs.AddSample(p, c)
					}
				}
			}
			defer wg.Done()
		}(i, cpus, width, height)
	}

	// for i := 0.0; i < iterations; i++ {
	// 	for x := 0.0; x < float64(width)*scale+x_step; x += x_step {
	// 		for y := 0.0; y < float64(height)*scale+y_step; y += y_step {

	// 			wg.Add(1)

	// 			go func(coords utils.Coord, width, height int, shaders ...utils.Color) {

	// 				p := transform(coords, width, height)
	// 				c := colorize(coords, width, height)

	// 				fs.AddSample(p, c)
	// 				// Decrement the counter when the goroutine completes.
	// 				defer wg.Done()
	// 			}(
	// 				utils.NewCoord(x, y),
	// 				width,
	// 				height,
	// 			)
	// 		}
	// 	}
	// }

	wg.Wait()

	fmt.Print("Rasterizing... ")
	r_start := time.Now()
	raster := fs.Rasterize()
	fmt.Println(time.Since(r_start))

	utils.SaveImage(int(s_width), int(s_height), raster, "test")

	fmt.Println("Finished in:", time.Since(start))
}
