## Sample cloud interpolation engine in Go

This is a simple progressively enhancing rendering engine. 

![ex3](https://github.mit.edu/chwalsh/rendering/raw/pretty/example/test_5.png)

The current configuration gives a plesent fragmented trigonometric fade. Run with:
``` Bash
go run main.go
```

A nice little development enviromnent with:
``` Bash
npm install
npm run-script watch --silent # re-runs on saves
```

### Examples: [more...](./example)

![ex2](https://github.mit.edu/chwalsh/rendering/raw/pretty/example/26.07.2015.04.49.34.548.png)
![ex1](https://github.mit.edu/chwalsh/rendering/raw/pretty/example/26.07.2015.01.59.17.129.png)
![ex4](https://github.mit.edu/chwalsh/rendering/raw/pretty/example/test_6.png)

This is the general idea of how to use to. See the [well](./main.go) [documented](./utils/utils.go) [source](./sampler/sampler.go) [code](./sampler/sampler_test.go) for specifics.
``` Go
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

fs := sampler.GaussianFrameSampler(width, height, 4,
		gauss, float64(max_d), 1)

for x := 0; x < width; x += x_step {
	for y := 0; y < height; y += y_step {
		coords := utils.NewCoord(x, y)
		p := utils.Randomize(coords, jitter)
		c := colorize(coords, s_width, s_height)

		fs.AddSample(p, c) // O(max_d^2)
	}
}

fs.Rasterize()

```
Note: This code is a simplified illustration. See [main.go](./main.go) for the full (well commented!) implementation.

## This is just the beginning! 

Much more is possible by modifying the colorize and transform functions to modifying various parameters. The interpolation engine generally handles making the images look good, so you can focus on abusing the parameters.

### Ray Tracing?

I made this out of interest in ray tracing. I made this to make the implementation of a ray tracer easier. The ray tracer just has to supply the image grid with samples, and it will handle interpolating between them and creating the raster.

## Monte Carlo importance first sampling

This could be easily modified to also return sample importance on sample inserts. That could be combined with a priority queue to act as an easy importance first sampling algorithm for a ray tracer, or similar.

### What about GPUs?

Yes, this would be faster in a GPU. Once I refine my algorithms and have more of an idea where I want to take this project, I am considering either incorporating (through a C wrapper) or moving to Cuda C.

I have some ideas for how to do this with a preprocessor, to remove all of the nasty boilerplate. Golang has pretty good support for that kind of thing, and I've [recently gotten pretty familiar](https://github.com/kctess5/Go-lexer-parser) with parsing and abstract syntax trees...

## Profiling and more concurrency!

The code is decently fast out of the box, I generally see ~1-3s render times for 640x480 samples, and then when I ramp up the resolution it's around 30s-2m for high res images. I've optimized this a tad, but I need sit down with pprof for a little while. More on this later... 

I suspect that there is some sub-optimal memory usage going on, and that some careful refactoring with more channels could help things.