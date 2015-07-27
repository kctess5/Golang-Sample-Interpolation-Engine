package sampler

import (
	// "image/color"
	// "fmt"
	"../utils"
	"image/color"
	"math"
)

type Dist func([]float64) float64

func GaussianFactory(x_center, y_center, u_0 float64) Dist {

	bottomExp := -2 * u_0 * u_0
	bottomTerm := math.Sqrt(2*math.Pi) * u_0

	return func(coords []float64) float64 {

		x_r := coords[0] - x_center
		y_r := coords[1] - y_center

		r := x_r*x_r + y_r*y_r

		return math.Exp(r/bottomExp) / bottomTerm

	}
}

type WeightedSampler struct {
	data    []float64
	mass    float64
	norm    float64
	distFxn Dist
}

func (ws *WeightedSampler) getValue() []float64 {
	output := make([]float64, len(ws.data))

	for i, v := range ws.data {
		output[i] = ws.norm * v / ws.mass
	}
	return output
}

func (ws *WeightedSampler) AddSample(coords, value []float64) {
	sampleWeight := ws.distFxn(coords)

	for i, v := range value {
		ws.data[i] += v * sampleWeight
	}

	ws.mass += sampleWeight
}

// args:
func Gaussian2DPixelSampler(x, y, gauss float64, w int) *WeightedSampler {
	return &WeightedSampler{
		data:    make([]float64, w),
		mass:    0,
		norm:    1,
		distFxn: GaussianFactory(x, y, gauss),
	}
}

type FrameSampler struct {
	pixelSamplers [][]*WeightedSampler
	max_d         float64
	width         int
	height        int
	depth         int
}

func GaussianFrameSampler(w, h, d int, gauss, max_d float64) *FrameSampler {
	dim2PixelSampler := make([][]*WeightedSampler, w)

	for i := range dim2PixelSampler {
		dim2PixelSampler[i] = make([]*WeightedSampler, h)
	}

	fs := &FrameSampler{
		pixelSamplers: dim2PixelSampler,
		max_d:         max_d,
		width:         w,
		height:        h,
		depth:         d,
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			fs.pixelSamplers[x][y] = Gaussian2DPixelSampler(
				float64(x), float64(y), gauss, d,
			)
		}
	}
	return fs
}

func (fs *FrameSampler) AddSample(c utils.Coord, value color.Color) {

	var top, bottom, left, right int

	if fs.max_d > 0 {
		top = int(math.Min(float64(fs.height), c.Y()+fs.max_d))
		left = int(math.Max(0, c.X()-fs.max_d))
		right = int(math.Min(float64(fs.width), c.X()+fs.max_d))
		bottom = int(math.Max(0, c.Y()-fs.max_d))
	} else {
		top = fs.height
		left = 0
		right = fs.width
		bottom = 0
	}

	for x := left; x < right; x++ {
		for y := bottom; y < top; y++ {

			r_x := (float64(x) - c.X())
			r_y := (float64(y) - c.Y())

			r := r_x*r_x + r_y*r_y

			if r <= fs.max_d*fs.max_d {
				r, g, b, a := value.RGBA()
				fs.pixelSamplers[x][y].AddSample(c,
					[]float64{float64(r), float64(g), float64(b), float64(a)})
			}
		}
	}
}

func (fs *FrameSampler) Rasterize() *utils.Frame {
	frame := utils.Dim3(fs.width, fs.height, fs.depth)

	// fmt.Println(frame)

	for x := 0; x < fs.width; x++ {
		for y := 0; y < fs.height; y++ {
			frame[x][y] = fs.pixelSamplers[x][y].getValue()
		}
	}
	return utils.NewFrame(frame, fs.width, fs.height)
}
