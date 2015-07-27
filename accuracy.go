package main

import "fmt"
import "math"

func main() {
	a := float64(0.0)
	for i := 0; i < 100; i++ {
		a += float64(i) * 2 / 3
		a -= float64(i) * 1 / 3
		a -= float64(i) * 1 / 3
	}
	var i int
	for i = 0; math.Abs(a) < 1.0; i++ {
		a += a
		// fmt.Println(a)
	}
	fmt.Printf("After %d iterations, a = %e\n", i, a)
}
