package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func FailGracefully(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SaveImage(width, height int, frame [][][]float64, fn string) {
	out, err := os.Create("./" + fn)
	FailGracefully(err)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			r, g, b, a := uint8(frame[x][y][0]),
				uint8(frame[x][y][1]),
				uint8(frame[x][y][2]),
				uint8(frame[x][y][3])

			fill := color.RGBA{r, g, b, a}

			draw.Draw(img, image.Rect(x, y, x+1, y+1), &image.Uniform{fill}, image.ZP, draw.Src)
		}
	}

	// ok, write out the data into the new PNG file

	FailGracefully(png.Encode(out, img))
}
