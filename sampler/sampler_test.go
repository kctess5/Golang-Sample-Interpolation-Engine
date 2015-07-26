package sampler

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSampler(t *testing.T) {
	sampler := Gaussian2DPixelSampler(0, 0, 1, 4)

	data := []float64{100, 100, 100, 100}

	sampler.AddSample(Coord(0, 0), data)

	assert.Equal(t, sampler.getValue(), data, "single sample should be value")
}

func TestSampler_offset(t *testing.T) {
	sampler := Gaussian2DPixelSampler(1, 1, 1, 4)

	data := []float64{100, 100, 100, 100}

	sampler.AddSample(Coord(1, 1), data)

	assert.Equal(t, sampler.getValue(), data, "single sample equal offset should be value")
}

// func TestFrameSampler(t *testing.T) {
// 	// fs :=
// 	GaussianFrameSampler(10, 10, 4, 1, 10)

// 	// data := []float64{100, 100, 100, 100}

// 	// fs.AddSample(Coord(0,0), data)
// }
