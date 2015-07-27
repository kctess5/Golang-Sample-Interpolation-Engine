package main

import (
	"./sampler"
	"./utils"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"
)

// const width = 640
// const height = 480

// const width = 1280
// const height = 960

// const width = 1920
// const height = 1080

const scale = 2
const size = 1

const x_step = 50 * scale * size
const y_step = 50 * scale * size
const iterations = 1
const max_d = 200 * scale * size
const positionChaos = 1 * scale * size
const swirlChaos = 1 * scale * size
const positionChaos3 = 3 * scale * size

// const x_step = 14 * scale
// const y_step = 14 * scale
// const iterations = 2
// const max_d = 35 * scale
// const positionChaos = 10 * scale
// const swirlChaos = 4 * scale
// const positionChaos3 = 3 * scale

const gauss = 100 * scale
const swirl = 150 * scale * size
const angle = math.Pi / 1

const src = "./output/test_8.png"

func PSin(i float64) float64 {
	return (math.Sin(i) + 1) / 2
}

func transform(c utils.Coord, width, height int) utils.Coord {
	// return utils.Randomize(c, positionChaos)
	return utils.Randomize(
		utils.Swirl(utils.Randomize(c, swirlChaos),
			float64(width)/(2/scale),
			float64(height)/(2/scale),
			// angle, swirl,
			angle*(utils.GetRand(1)), (.75+utils.GetRand(0.25))*swirl,
		),
		positionChaos)
}

func transform3(c utils.Coord, width, height int) utils.Coord {
	// return utils.Randomize(c, positionChaos)
	return utils.Randomize(
		utils.Swirl(c,
			float64(width)/(2/scale),
			float64(height)/(2/scale),
			angle*utils.GetRand(30), utils.GetRand(30)*swirl,
			// -.5*angle*utils, swirl,
		),
		positionChaos3)
}

func transform2(c utils.Coord) utils.Coord {
	return utils.Randomize(c, swirlChaos)
	// return c
}

func colorize(x, y float64, width, height int) utils.Color {
	// Sqrt := math.Sqrt

	// fmt.Println(100 * math.Pow(PSin(float64(y*x)/float64(width*width)), 2))

	return utils.NewColorFrac(
		math.Pow(PSin(math.Pi+float64(y-x)/float64(height)), 2), // Red
		math.Pow(PSin(float64(y*x)/float64(width*width)), 2),    // Blue
		math.Pow(PSin(float64(x+y)/float64(height+width)), 2),   // Green
		1)
}

func process(c utils.Color, x, y float64, width, height int) utils.Color {
	r, g, b, a := c.FloatRGBA()

	w := float64(width)
	h := float64(height)
	wh := float64(width * height)

	r *= math.Exp(2*math.Pi + 30*x*y/wh)
	g *= math.Exp(2 + 2*y/h)
	b *= math.Tan(2*math.Pi + 3*x/w)

	return utils.NewColor(r, g, b, a)
}

func ConvColor(c color.Color) utils.Color {
	r, g, b, a := c.RGBA()

	return utils.NewColor(float64(r), float64(g), float64(b), float64(a))
}

func main() {
	utils.StartNetProfile("profile")
	rand.Seed(time.Now().UTC().UnixNano())
	start := time.Now()

	width, height, img := utils.LoadImage(src)

	s_width := int(float64(width) * scale)
	s_height := int(float64(height) * scale)

	fs := sampler.GaussianFrameSampler(
		s_width, s_height, 4, gauss, float64(max_d))

	fmt.Println("Adding Samples")

	for i := 0.0; i < iterations; i++ {
		for x := 0.0; x < float64(width)*scale; x += x_step {
			for y := 0.0; y < float64(height)*scale; y += y_step {
				fs.AddSample(
					transform(utils.NewCoordInt(int(x), int(y)), width, height),
					process(ConvColor(img.At(int(x/scale), int(y/scale))), x, y, width, height),
				)
			}
		}
	}

	fmt.Println("Saveing Image")
	utils.SaveImage(s_width, s_height, fs.Rasterize(), "test")
	fmt.Println("Saved Image after:", time.Since(start), "Seconds")
}
