package raycasting

func Sphere(p Vec3, s float64) float64 {
	return p.Len() - s
}
