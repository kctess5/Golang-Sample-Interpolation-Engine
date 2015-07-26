package sampler

import (
	// "image/color"
	// "fmt"
	"math"
)

func dim1(w int) []float64 {
	// allocate composed 1d array
	a := make([]float64, w)
	return a
}

func dim2(w, h int) [][]float64 {
	// allocate composed 2d array
	a := make([][]float64, w)
	for i := range a {
		a[i] = make([]float64, h)
	}
	return a
}

func dim3(w, h, d int) [][][]float64 {
	// allocate composed 3d array
	a := make([][][]float64, w)
	for i := range a {
		a[i] = make([][]float64, h)
		for j := range a[i] {
			a[i][j] = make([]float64, d)
		}
	}
	return a
}

type coord []float64

func Coord(x, y int) coord {
	return []float64{float64(x), float64(y)}
}

func (c coord) x(x ...float64) float64 {
	if len(x) > 0 {
		c[0] = x[0]
	}
	return c[0]
}
func (c coord) y(y ...float64) float64 {
	if len(y) > 0 {
		c[1] = y[0]
	}
	return c[1]
}

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
		output[i] = v / (ws.mass / ws.norm)
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

func (fs *FrameSampler) AddSample(c coord, value []float64) {

	var top, bottom, left, right int

	if fs.max_d > 0 {
		top = int(math.Min(float64(fs.height), c.y()+fs.max_d))
		left = int(math.Max(0, c.x()-fs.max_d))
		right = int(math.Min(float64(fs.width), c.x()+fs.max_d))
		bottom = int(math.Max(0, c.y()-fs.max_d))
	} else {
		top = fs.height
		left = 0
		right = fs.width
		bottom = 0
	}

	for x := left; x < right; x++ {
		for y := bottom; y < top; y++ {

			r_x := (float64(x) - c.x())
			r_y := (float64(y) - c.y())

			r := r_x*r_x + r_y*r_y

			if r <= fs.max_d*fs.max_d {
				fs.pixelSamplers[x][y].AddSample(c, value)
			}
		}
	}
}

func (fs *FrameSampler) Rasterize() [][][]float64 {
	frame := dim3(fs.width, fs.height, fs.depth)

	// fmt.Println(frame)

	for x := 0; x < fs.width; x++ {
		for y := 0; y < fs.height; y++ {
			frame[x][y] = fs.pixelSamplers[x][y].getValue()
		}
	}
	return frame
}

// func FrameSampler()

//
//
//
//
//
//
//
//
//
//
// f\left(x\right) = a \exp{\left(- { \frac{(x-b)^2 }{ 2 c^2} } \right)}
// a*exp()
