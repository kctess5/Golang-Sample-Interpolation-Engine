package main

import (
	"./sampler"
	"./utils"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// abstract frame size
const width = 640
const height = 480

// scaling an image will strictly make it higher resolution
// (ideally) it should look the same, just bigger, and slower
const scale = 1.0

// sample density
const x_step = 10 * scale
const y_step = 10 * scale

const iterations = 1 // number of passes

/*
	this affects the frame sampler sampling algorithm. Setting it
	higher will allow for cause the algorithm to sample a larger search
	space. If it is too low, then the render will have holes. To high
	and it will be needlessly slow
*/
const max_d = 30 * scale
const positionChaos = 10 * scale // sample jitter

// constant for gaussian sampling
const gauss = 1.2

func transform(c utils.Coord, width, height float64) utils.Coord {
	// add jitter
	return utils.Randomize(c, positionChaos)
}

func colorize(c utils.Coord, width, height float64) utils.Color {
	x, y := c.X(), c.Y()
	w, h, wh := width, height, width*height

	// a simple trigonometric fade

	r := 1.0 * math.Sin(x*y/wh)
	g := 1.0 * math.Sin(y/h)
	b := 1.0 * math.Sin(x/w)
	a := 1.0

	return utils.NewColorFrac(r, g, b, a)
}

func main() {
	fmt.Print("Initializing Rendering Engine. ")

	cpus := 8
	runtime.GOMAXPROCS(cpus) // use all CPUs

	rand.Seed(time.Now().UTC().UnixNano())
	start := time.Now()

	// these are the scaled dims of the output image, useful later
	s_height, s_width := scale*height, scale*width

	// initialize Gaussian frame sampler. This aggregates samples
	// into rasterized images progressively.
	fs := sampler.GaussianFrameSampler(int(s_width), int(s_height), 4,
		gauss, float64(max_d), cpus)

	fmt.Println("Adding samples using:", cpus, "CPUs.")

	var wg sync.WaitGroup

	for i := 0; i < cpus; i++ {
		wg.Add(1) // increment wait counter

		// starts one goroutine per cpu. Each goroutine covers a
		// fraction of the work area

		step := x_step * float64(cpus)
		x_bound, y_bound := s_width+x_step, s_height+y_step

		go func(i, width, height int) {

			/*
				Start offset by i, then increment by the specified step size,
				go one extra step at the end for good measure'

				O(width/x_step * height/y_step * iterations * max_d^2)
			*/
			for x := x_step * float64(i); x < x_bound; x += step {
				for y := y_step; y < y_bound; y += y_step {

					coords := utils.NewCoord(x, y)

					for i := 0.0; i < iterations; i++ {

						/*
							calculate sample position as a funciton of x, y. this
							allows us to transform the distribution of samples
							very easily
						*/
						p := transform(coords, s_width, s_height)

						// calculate color as a function of position in frame
						c := colorize(coords, s_width, s_height)

						fs.AddSample(p, c) // O(max_d^2)
					}
				}
			}
			defer wg.Done() // decrement wait counter
		}(i, width, height)
	}

	// wait for all goroutines to finish
	wg.Wait()

	fmt.Print("Rasterizing... ")
	r_start := time.Now()

	// Rasterize calls are actually cheap operations, since all of the
	// work is done progressively in the for loop
	raster := fs.Rasterize()

	// Save to PNG and give stats
	fmt.Println(time.Since(r_start))
	utils.SaveImage(int(s_width), int(s_height), raster, "test")
	fmt.Println("Finished in:", time.Since(start))
}
