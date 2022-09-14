package utils

import (
	"math"
	"math/rand"
)

var colors = [][3]int{{255, 21, 21}, {255, 243, 21}, {21, 87, 255}, {21, 255, 208}, {255, 21, 224}}

func RandomColor() [3]int {
	return colors[rand.Intn(len(colors))]
}

func CalcDistance(x1 float32, x2 float32, y1 float32, y2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2)))
}

func CalcDistance64(x1 float64, x2 float64, y1 float64, y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
}

func Lerp() {

}

func Clip(v float32, min float32, max float32) float32 {
	return float32(math.Min(math.Max(float64(v), float64(min)), float64(max)))
}

func Clip64(v float64, min float64, max float64) float64 {
	return math.Min(math.Max(v, min), max)
}

func SurfaceArea(r float32) float32 {
	return float32(math.Pi * math.Pow(float64(r), 2))
}

func SurfaceArea64(r float64) float64 {
	return math.Pi * math.Pow(r, 2)
}

func IntersectionArea(r1 float64, r2 float64, dist float64) float64 {
	if dist < math.Abs(r1-r2) {
		return SurfaceArea64(math.Min(r1, r2))
	}
	if dist > math.Abs(r1+r2) {
		return 0
	}
	dist1 := (r1*r1 - r2*r2 + dist*dist) / (2 * dist)
	dist2 := dist - dist1
	h1 := math.Sqrt(r1*r1 - dist1*dist1)
	h2 := math.Sqrt(r2*r2 - dist2*dist2)
	return r1*r1*math.Acos(dist1/r1) - dist1*h1 + r2*r2*math.Acos(dist2/r2) - dist2*h2
}
