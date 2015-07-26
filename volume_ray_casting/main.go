// package main

// import (
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"image/draw"
// 	"image/png"
// 	"math/rand"
// 	"os"
// 	"time"
// )

// func failGracefully(err error) {
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

// func main() {
// 	rand.Seed(time.Now().UTC().UnixNano())

// 	out, err := os.Create("./output.png")
// 	failGracefully(err)

// 	// generate some QR code look a like image

// 	const x_max = 780
// 	const y_max = 100

// 	img := image.NewRGBA(image.Rect(0, 0, x_max, y_max))

// 	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)

// 	for y := 0; y < y_max; y += 1 {
// 		for x := 0; x < x_max; x += 1 {

// 			r, g, b := uint8(float64(100)), uint8(float64(100)), uint8(float64(100))

// 			fill := color.RGBA{r, g, b, 255}

// 			draw.Draw(img, image.Rect(x, y, x+1, y+1), &image.Uniform{fill}, image.ZP, draw.Src)

// 			// 		fill := &image.NewRGBA(255, 0, 0, 1)
// 			// 		// if rand.Intn(10)%2 == 0 {
// 			// 		// 	fill = &image.Uniform{color.White}
// 			// 		// }
// 			// 		draw.Draw(img, image.Rect(x, y, x+10, y+10), fill, image.ZP, draw.Src)

// 		}
// 	}

// 	// ok, write out the data into the new PNG file

// 	failGracefully(png.Encode(out, img))

// 	fmt.Println("Generated image to output.png \n")
// }
