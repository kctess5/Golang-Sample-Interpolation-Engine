package utils

import (
	"fmt"
	"image"
	"image/color"
	// "image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	// "runtime/pprof"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// func StartProfile(destpath string) {
// 	f, err := os.Create(destpath)
// 	FailGracefully(err)
// 	pprof.StartCPUProfile(f)
// }

func StartNetProfile(destpath string) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

const ColorMax = 0xFFFF

func GetRand(a float64) float64 {
	return a * (rand.Float64()*2 - 1)
}

/*
	Transform the input by a swirl.

	Args:
	- c: input coordinates
	- x_o, y_o: swirl center x, swirl center y
	- rot: amount of rotation in radians
	- effect: intensity
	Returns:
	- Coord: translated coordinates
*/
func Swirl(c Coord, x_o, y_o, rot, effect float64) Coord {
	x, y := c.X(), c.Y()

	x = x - x_o
	y = y - y_o

	angle := rot * math.Exp(-(x*x+y*y)/(effect*effect))

	u := math.Cos(angle)*x + math.Sin(angle)*y + x_o
	v := -math.Sin(angle)*x + math.Cos(angle)*y + y_o

	return NewCoord(u, v)
}

/*
	Add [-amount:amount] jitter to input coordinate

	Args:
	- c: input coordinates
	- amount: range of jitter
	Returns:
	- Coord: translated coordinates
*/
func Randomize(c Coord, amount float64) Coord {
	c.Y(c.Y() + GetRand(amount))
	c.X(c.X() + GetRand(amount))
	return c
}

/*
	A basic 64 bit coordinate warpper.
*/
type Coord []float64

func NewCoord(x, y float64) Coord {
	return Coord{x, y}
}

func NewCoordInt(x, y int) Coord {
	return NewCoord(float64(x), float64(y))
}

func (c Coord) X(x ...float64) float64 {
	if len(x) > 0 {
		c[0] = x[0]
	}
	return c[0]
}
func (c Coord) Y(y ...float64) float64 {
	if len(y) > 0 {
		c[1] = y[0]
	}
	return c[1]
}

/*
	Exit if a problem occurs
*/
func FailGracefully(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func makeTimestamp() string {
	t := time.Now()
	return t.Format("02.01.2006.15.04.05.000")
}

/*
	A slightly more convenient color implementation
*/
type Color struct {
	d []float64
}

func NewColor(r, g, b, a float64) Color {
	return Color{[]float64{r, g, b, a}}
}

func NewColorFrac(r, g, b, a float64) Color {
	return Color{
		[]float64{r * ColorMax, g * ColorMax, b * ColorMax, a * ColorMax}}
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.d[0]), uint32(c.d[1]), uint32(c.d[2]), uint32(c.d[3])
}

func (c Color) FloatRGBA() (r, g, b, a float64) {
	return c.d[0], c.d[1], c.d[2], c.d[3]
}

func (c Color) Src() []float64 {
	return c.d
}

/*
	Wraps a rasterized frame object. Fulifills the Image interface
*/
type Frame struct {
	image  [][][]float64
	width  int
	height int
}

func NewFrame(img [][][]float64, w, h int) *Frame {
	return &Frame{
		image:  img,
		width:  w,
		height: h,
	}
}

func (f *Frame) setAt(x, y int, c Color) {
	f.image[x][y] = c.d
}

func (f Frame) At(x, y int) color.Color {
	return Color{f.image[x][y]}
}

func (f Frame) ColorModel() color.Model {
	c := image.NRGBA64{}
	return c.ColorModel()
}

func (f Frame) Bounds() image.Rectangle {
	return image.Rect(0, 0, f.width, f.height)
}

func SaveImage(width, height int, frame *Frame, fn string) {
	out, err := os.Create("./output/" + makeTimestamp() + ".png")
	FailGracefully(err)
	out2, err := os.Create("./output/" + fn + ".png")
	FailGracefully(err)

	fmt.Println("Saving As: ", makeTimestamp())

	go FailGracefully(png.Encode(out, frame))
	go FailGracefully(png.Encode(out2, frame))
}

func LoadImage(fn string) (int, int, *Frame) {
	infile, err := os.Open(fn)
	FailGracefully(err)
	src, _, err := image.Decode(infile)
	FailGracefully(err)

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	image := NewFrame(Dim3(w, h, 4), w, h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := src.At(x, y).RGBA()

			// if int(a) == 0 {
			// 	// fmt.Print("test")
			// 	image.setAt(x, y,
			// 		NewColorFrac(.5, .5, .5, .5))
			// } else {
			image.setAt(x, y,
				NewColor(float64(r), float64(g), float64(b), float64(a)))
			// }

		}
	}

	return w, h, image
}

/*
	Convienient multi-dimensional array wrappers
*/
func Dim1(w int) []float64 {
	// allocate composed 1d array
	a := make([]float64, w)
	return a
}

func Dim2(w, h int) [][]float64 {
	// allocate composed 2d array
	a := make([][]float64, w)
	for i := range a {
		a[i] = make([]float64, h)
	}
	return a
}

func Dim3(w, h, d int) [][][]float64 {
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
