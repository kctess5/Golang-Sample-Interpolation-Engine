package sampler

import (
	"../utils"
	"image/color"
	"math"
	"sync"
)

// Distribution. Takes a vector, returns a weight
type Dist func([]float64) float64

// Constructs a gaussian distribution function with the given
// parameters
func GaussianFactory(x_center, y_center, u_0 float64) Dist {

	bottomExp := -2 * u_0 * u_0
	bottomTerm := math.Sqrt(2*math.Pi) * u_0

	return func(coords []float64) float64 {

		x_r := coords[0] - x_center
		y_r := coords[1] - y_center

		r := x_r*x_r + y_r*y_r

		// gaussian function
		return math.Exp(r/bottomExp) / bottomTerm
	}
}

/*
	A single pixel sampler.
*/
type WeightedSampler struct {
	data    []float64
	mass    float64
	norm    float64
	distFxn Dist
	lock    sync.Mutex
}

/*
	Returns current belief of this samplers aggregate value.
	Normalizes it with respect to the gathered gaussian mass
*/
func (ws *WeightedSampler) getValue() []float64 {
	output := make([]float64, len(ws.data))

	for i, v := range ws.data {
		output[i] = ws.norm * v / ws.mass
	}
	return output
}

/*
	Add in a sample with an associated color weight
*/
func (ws *WeightedSampler) AddSample(coords, value []float64) {
	sampleWeight := ws.distFxn(coords)

	if sampleWeight <= 0 || coords[0] < 0.0 || coords[1] < 0.0 {
		// fail in the bad cases
		return
	}

	// acquite mutex lock on this pixel sampler
	ws.lock.Lock()

	// add in sample value
	for i, v := range value {
		ws.data[i] += v * sampleWeight
	}

	// update gaussian mass
	ws.mass += sampleWeight

	ws.lock.Unlock()
}

func Gaussian2DPixelSampler(x, y, gauss float64, w int) *WeightedSampler {
	return &WeightedSampler{
		data:    make([]float64, w),
		mass:    0,
		norm:    1,
		distFxn: GaussianFactory(x, y, gauss),
		lock:    sync.Mutex{},
	}
}

/*
	An abstract frame object, comprised of a 2d array of frame samplers
*/
type FrameSampler struct {
	pixelSamplers [][]*WeightedSampler
	max_d         float64
	width         int
	height        int
	depth         int
	max_CPU       int
}

func GaussianFrameSampler(
	w, h, d int,
	gauss, max_d float64,
	max_CPU int) *FrameSampler {
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
		max_CPU:       max_CPU,
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

/*
	Add in a sample, update nearby pixel samplers within max_d radius
*/
func (fs *FrameSampler) AddSample(c utils.Coord, value color.Color) {

	var top, bottom, left, right int

	if fs.max_d > 0 {
		// makes sure not to index outside of the pixel sampler bounds
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

	// iterate over local box
	for x := left; x < right; x++ {
		for y := bottom; y < top; y++ {

			r_x := (float64(x) - c.X())
			r_y := (float64(y) - c.Y())

			r := r_x*r_x + r_y*r_y

			// make sure that the pixel is within the specified radius bound
			if r <= fs.max_d*fs.max_d {
				// convert color type and add sample
				r, g, b, a := value.RGBA()

				fs.pixelSamplers[x][y].AddSample(c,
					[]float64{float64(r), float64(g), float64(b), float64(a)})
			}
		}
	}
}

/*
	Iterate over all pixel samplers and collect current belief.
	Multithreaded.
*/
func (fs *FrameSampler) Rasterize() *utils.Frame {
	frame := utils.Dim3(fs.width, fs.height, fs.depth)

	// start up one goroutine per thread to collect samples.
	var wg sync.WaitGroup
	for i := 0; i < fs.max_CPU; i++ {
		wg.Add(1)

		go func(i, max int) {
			for x := i; x < fs.width; x += max {
				for y := 0; y < fs.height; y++ {
					// Collect pixel value
					frame[x][y] = fs.pixelSamplers[x][y].getValue()
				}
			}
			defer wg.Done()

		}(i, fs.max_CPU)
	}
	// wait for all samples
	wg.Wait()

	return utils.NewFrame(frame, fs.width, fs.height)
}

func (fs *FrameSampler) Width() int {
	return fs.width
}

func (fs *FrameSampler) Height() int {
	return fs.height
}
